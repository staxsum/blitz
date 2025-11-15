package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const banner = `
██████╗ ██╗     ██╗████████╗███████╗
██╔══██╗██║     ██║╚══██╔══╝╚══███╔╝
██████╔╝██║     ██║   ██║     ███╔╝ 
██╔══██╗██║     ██║   ██║    ███╔╝  
██████╔╝███████╗██║   ██║   ███████╗
╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝
                                     
  Fast Web Form Cracking
  BY: stax 
  https://github.com/staxsum/blitz
`

var (
	targetURL    string
	usernameFile string
	passwordFile string
	threads      int
	timeout      int
	rateLimit    int
	verbose      bool
	skipWAF      bool
	formIndex    int
)

func init() {
	flag.StringVar(&targetURL, "url", "", "Target URL to test (required)")
	flag.StringVar(&usernameFile, "usernames", "usernames.txt", "Path to username wordlist")
	flag.StringVar(&passwordFile, "passwords", "passwords.txt", "Path to password wordlist")
	flag.IntVar(&threads, "threads", 5, "Number of concurrent threads")
	flag.IntVar(&timeout, "timeout", 10, "Request timeout in seconds")
	flag.IntVar(&rateLimit, "rate", 100, "Requests per second limit")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&skipWAF, "skip-waf", false, "Skip WAF detection")
	flag.IntVar(&formIndex, "form", 0, "Form index to target (default: auto-detect)")
}

func main() {
	printBanner()

	flag.Parse()

	// Validate required flags
	if targetURL == "" {
		color.Red("[-] Error: Target URL is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  ./blitz -url http://example.com/login")
		os.Exit(1)
	}

	// Display legal disclaimer (no user confirmation required)
	printLegalDisclaimer()

	// Load wordlists
	usernames, err := loadWordlist(usernameFile)
	if err != nil {
		color.Red("[-] Failed to load usernames: %v", err)
		os.Exit(1)
	}
	color.Green("[+] Loaded %d usernames", len(usernames))

	passwords, err := loadWordlist(passwordFile)
	if err != nil {
		color.Red("[-] Failed to load passwords: %v", err)
		os.Exit(1)
	}
	color.Green("[+] Loaded %d passwords", len(passwords))

	// Initialize the scanner
	scanner, err := NewScanner(targetURL, timeout, verbose)
	if err != nil {
		color.Red("[-] Failed to initialize scanner: %v", err)
		os.Exit(1)
	}

	// Run security checks
	color.Cyan("\n[*] Running security checks...")
	scanner.RunSecurityChecks(skipWAF)

	// Find and analyze forms
	forms, err := scanner.FindForms()
	if err != nil {
		color.Red("[-] Failed to find forms: %v", err)
		os.Exit(1)
	}

	if len(forms) == 0 {
		color.Red("[-] No forms found on target page")
		os.Exit(1)
	}

	color.Green("[+] Found %d form(s)", len(forms))

	// Select target form
	var targetForm *Form
	if formIndex >= 0 && formIndex < len(forms) {
		targetForm = forms[formIndex]
	} else {
		targetForm = selectForm(forms)
	}

	if targetForm == nil {
		color.Red("[-] No suitable form found for testing")
		os.Exit(1)
	}

	// Initialize brute forcer
	bruteForcer := NewBruteForcer(scanner, targetForm, threads, rateLimit)

	// Start brute force attack
	color.Cyan("\n[*] Starting credential testing...")
	color.Yellow("[!] Testing %d username(s) with %d password(s) each", len(usernames), len(passwords))

	result := bruteForcer.Start(usernames, passwords)

	if result != nil {
		color.Green("\n[+] ✓ Valid credentials found!")
		color.Green("    Username: %s", result.Username)
		color.Green("    Password: %s", result.Password)

		// Save results
		saveResults(result)
	} else {
		color.Red("\n[-] No valid credentials found")
	}
}

func printBanner() {
	color.Cyan(banner)
}

func printLegalDisclaimer() {
	color.Yellow("\nLEGAL DISCLAIMER")
	color.White("This tool is designed for authorized security testing only.")
	color.White("Unauthorized access to computer systems is illegal.")
	color.White("You must have explicit permission to test the target system.")
	color.White("\nBy using this tool, you agree that you have proper authorization.")
	color.White("The author assumes NO LIABILITY for misuse or damage.\n")
}

func loadWordlist(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func selectForm(forms []*Form) *Form {
	color.Cyan("\n[*] Available forms:")
	for i, form := range forms {
		color.White("  [%d] %s", i, form.Description())
	}

	// Auto-select first login-like form
	for _, form := range forms {
		if form.UsernameField != "" && form.PasswordField != "" {
			color.Green("[+] Auto-selected form with username and password fields")
			return form
		}
	}

	return nil
}

func saveResults(result *Credential) {
	filename := "blitz_results.txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Yellow("[!] Could not save results: %v", err)
		return
	}
	defer file.Close()

	output := fmt.Sprintf("\n=== Blitz Results ===\nTarget: %s\nUsername: %s\nPassword: %s\nTimestamp: %v\n",
		targetURL, result.Username, result.Password, result.Timestamp)

	file.WriteString(output)
	color.Green("[+] Results saved to %s", filename)
}
