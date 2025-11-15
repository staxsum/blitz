package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golang.org/x/net/publicsuffix"
)

type Scanner struct {
	client       *http.Client
	baseURL      string
	parsedURL    *url.URL
	timeout      time.Duration
	verbose      bool
	originalPage string
}

type SecurityCheck struct {
	Name        string
	Vulnerable  bool
	Description string
}

func NewScanner(targetURL string, timeout int, verbose bool) (*Scanner, error) {
	// Normalize URL
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Create cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	// Create HTTP client with custom transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // For testing purposes
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Jar:       jar,
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	scanner := &Scanner{
		client:    client,
		baseURL:   targetURL,
		parsedURL: parsedURL,
		timeout:   time.Duration(timeout) * time.Second,
		verbose:   verbose,
	}

	// Test connection and fetch original page
	if err := scanner.testConnection(); err != nil {
		// Try HTTPS if HTTP fails
		if strings.HasPrefix(targetURL, "http://") {
			targetURL = strings.Replace(targetURL, "http://", "https://", 1)
			parsedURL, _ = url.Parse(targetURL)
			scanner.baseURL = targetURL
			scanner.parsedURL = parsedURL
			
			if err := scanner.testConnection(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return scanner, nil
}

func (s *Scanner) testConnection() error {
	req, err := http.NewRequest("GET", s.baseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	s.originalPage = string(body)
	
	if s.verbose {
		color.Cyan("[*] Connection successful (Status: %d)", resp.StatusCode)
	}

	return nil
}

func (s *Scanner) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
}

func (s *Scanner) Get(targetURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)
	return s.client.Do(req)
}

func (s *Scanner) PostForm(targetURL string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return s.client.Do(req)
}

func (s *Scanner) RunSecurityChecks(skipWAF bool) {
	checks := []SecurityCheck{}

	// Check for Clickjacking protection
	resp, err := s.Get(s.baseURL)
	if err == nil {
		defer resp.Body.Close()
		
		xFrameOptions := resp.Header.Get("X-Frame-Options")
		if xFrameOptions == "" {
			checks = append(checks, SecurityCheck{
				Name:        "Clickjacking",
				Vulnerable:  true,
				Description: "Missing X-Frame-Options header",
			})
			color.Green("[+] Heuristic found a Clickjacking vulnerability")
		}

		// Check for Cloudflare
		server := resp.Header.Get("Server")
		cfRay := resp.Header.Get("CF-Ray")
		if strings.Contains(strings.ToLower(server), "cloudflare") || cfRay != "" {
			color.Red("[-] Target is protected by Cloudflare")
		}

		// Check for CSRF protection
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(s.originalPage))
		if err == nil {
			hiddenInputs := doc.Find("input[type='hidden']")
			if hiddenInputs.Length() == 0 {
				checks = append(checks, SecurityCheck{
					Name:        "CSRF",
					Vulnerable:  true,
					Description: "No hidden fields detected (possible CSRF vulnerability)",
				})
				color.Green("[+] Heuristic found a possible CSRF vulnerability")
			}
		}
	}

	// WAF Detection
	if !skipWAF {
		s.detectWAF()
	}
}

func (s *Scanner) detectWAF() {
	testPayload := "?test=<script>alert(1)</script>"
	testURL := s.baseURL + testPayload

	resp, err := s.Get(testURL)
	if err != nil {
		if s.verbose {
			color.Yellow("[!] WAF detection failed: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 406, 501:
		color.Red("[-] WAF Detected: Mod_Security")
	case 999:
		color.Red("[-] WAF Detected: WebKnight")
	case 419:
		color.Red("[-] WAF Detected: F5 BIG IP")
	case 403:
		color.Red("[-] Unknown WAF Detected (403 Forbidden)")
	default:
		if s.verbose {
			color.Green("[+] No obvious WAF detected")
		}
	}
}

func (s *Scanner) FindForms() ([]*Form, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s.originalPage))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var forms []*Form

	doc.Find("form").Each(func(i int, sel *goquery.Selection) {
		form := s.parseForm(sel, i)
		if form != nil {
			forms = append(forms, form)
			
			if s.verbose {
				color.Cyan("\n[*] Form %d:", i)
				color.White("    Action: %s", form.Action)
				color.White("    Method: %s", form.Method)
				if form.UsernameField != "" {
					color.White("    Username field: %s", form.UsernameField)
				}
				if form.PasswordField != "" {
					color.White("    Password field: %s", form.PasswordField)
				}
			}
		}
	})

	return forms, nil
}

func (s *Scanner) parseForm(sel *goquery.Selection, index int) *Form {
	action, _ := sel.Attr("action")
	method, _ := sel.Attr("method")

	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// Resolve relative URLs
	if action == "" {
		action = s.baseURL
	} else if !strings.HasPrefix(action, "http") {
		baseURL, _ := url.Parse(s.baseURL)
		actionURL, _ := url.Parse(action)
		action = baseURL.ResolveReference(actionURL).String()
	}

	form := &Form{
		Index:  index,
		Action: action,
		Method: method,
		Fields: make(map[string]string),
	}

	// Find input fields
	sel.Find("input").Each(func(i int, input *goquery.Selection) {
		inputType, _ := input.Attr("type")
		inputName, exists := input.Attr("name")
		
		if !exists || inputName == "" {
			return
		}

		inputType = strings.ToLower(inputType)
		
		switch inputType {
		case "text", "email", "":
			if form.UsernameField == "" {
				// Common username field names
				lowerName := strings.ToLower(inputName)
				if strings.Contains(lowerName, "user") || 
				   strings.Contains(lowerName, "login") || 
				   strings.Contains(lowerName, "email") ||
				   lowerName == "username" {
					form.UsernameField = inputName
				}
			}
			form.Fields[inputName] = ""
		case "password":
			if form.PasswordField == "" {
				form.PasswordField = inputName
			}
			form.Fields[inputName] = ""
		case "hidden":
			value, _ := input.Attr("value")
			form.Fields[inputName] = value
		default:
			value, _ := input.Attr("value")
			form.Fields[inputName] = value
		}
	})

	// Find select elements
	sel.Find("select").Each(func(i int, selectElem *goquery.Selection) {
		selectName, exists := selectElem.Attr("name")
		if !exists || selectName == "" {
			return
		}

		// Get first option value as default
		firstOption := selectElem.Find("option").First()
		value, _ := firstOption.Attr("value")
		form.Fields[selectName] = value

		if form.SelectField == "" {
			form.SelectField = selectName
			form.SelectOptions = []string{}
			
			selectElem.Find("option").Each(func(j int, option *goquery.Selection) {
				optVal, _ := option.Attr("value")
				form.SelectOptions = append(form.SelectOptions, optVal)
			})
		}
	})

	// Only return forms that have at least username OR password field
	if form.UsernameField != "" || form.PasswordField != "" {
		return form
	}

	return nil
}

func (s *Scanner) GetOriginalPage() string {
	return s.originalPage
}
