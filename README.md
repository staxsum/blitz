# Blitz
Blitz is a modern, lightning-fast login page bruteforcer written in Go.

### Features
- [x] Blazing fast concurrent processing (10x faster than Python)
- [x] Smart form and field detection
- [x] CSRF and Clickjacking scanner
- [x] Cloudflare and WAF detector
- [x] Checks for login bypass via SQL injection
- [x] Multi-threaded worker pool (5-20 threads)
- [x] Intelligent success detection
- [x] Rate limiting protection
- [x] Cross-platform (Windows, Linux, macOS)
- [ ] Browser automation mode
- [ ] Proxy support

#### Requirements
- Go 1.21 or higher

### Installation

Open your terminal and enter:
```bash
git clone https://github.com/staxsum/blitz
```

Now enter the following command:
```bash
cd blitz
```

#### Windows Users
Build Blitz by running:
```powershell
go mod tidy
go build -o blitz.exe .
```

Or simply double-click `build.bat`

#### Linux/macOS Users
Build Blitz by running:
```bash
go mod tidy
go build -o blitz .
```

Or simply run:
```bash
./setup.sh
```

### Usage

Now run Blitz by entering:

**Windows:**
```powershell
.\blitz.exe -url http://example.com/login
```

**Linux/macOS:**
```bash
./blitz -url http://example.com/login
```

Now enter your target URL and Blitz will do its thing:

### Advanced Usage

**Fast mode (10 threads):**
```bash
./blitz -url http://target.com/login -threads 10
```

**Stealth mode (slow):**
```bash
./blitz -url http://target.com/login -threads 2 -rate 5
```

**Custom wordlists:**
```bash
./blitz -url http://target.com/login -usernames users.txt -passwords pass.txt
```

**Verbose output:**
```bash
./blitz -url http://target.com/login -verbose
```

### Contributing

Contributions are welcome! Here's how you can help:

#### Reporting Bugs
If you find a bug, please open an issue with:
- Your operating system
- Go version
- Steps to reproduce the bug
- Expected vs actual behavior

#### Pull Requests
Want to contribute code?

1. Fork the repo
2. Create a branch
3. Commit your changes
4. Push to the branch
5. Open a pull request

### Legal Warning

**FOR AUTHORIZED SECURITY TESTING ONLY**

This tool is for authorized penetration testing only. Unauthorized access to computer systems is illegal. You must have explicit permission before testing any system.

**The author assumes NO LIABILITY for misuse.**

### License
MIT License

### Acknowledgments
- Inspired by Blazy (Python-based bruteforcer)
- Built with Go for performance and concurrency
- Thanks to all contributors

### Contact
- GitHub: [github.com/staxsum/blitz](https://github.com/staxsum/blitz)
- Issues: [github.com/staxsum/blitz/issues](https://github.com/staxsum/blitz/issues)