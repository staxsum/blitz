# Blitz

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Version](https://img.shields.io/badge/version-1.0.0-green.svg)](https://github.com/moscovium-mc/blitz/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey.svg)]()
[![Tool Type](https://img.shields.io/badge/tool-bruteforce-red.svg)]()
[![Built for](https://img.shields.io/badge/built%20for-pentesting-red.svg)]()

[![GitHub Stars](https://img.shields.io/github/stars/moscovium-mc/blitz?style=social)](https://github.com/moscovium-mc/blitz/stargazers)
[![Forks](https://img.shields.io/github/forks/moscovium-mc/blitz?style=social)](https://github.com/moscovium-mc/blitz/network/members)

A modern, lightning-fast login page bruteforcer written in Go. Blitz is built for speed and efficiency, leveraging Go's concurrent processing to deliver results 10x faster than traditional Python-based tools.

## What it does

- **Blazing fast concurrent processing** - 10x faster than Python equivalents
- **Smart form and field detection** - Automatically identifies login forms and input fields
- **CSRF and Clickjacking scanner** - Built-in security analysis
- **Cloudflare and WAF detector** - Identifies protection mechanisms
- **SQL injection bypass detection** - Tests for login bypass vulnerabilities
- **Multi-threaded worker pool** - Configurable 5-20 threads
- **Intelligent success detection** - Smart analysis of login responses
- **Rate limiting protection** - Prevents getting blocked during testing
- **Cross-platform** - Windows, Linux, macOS support

## Roadmap

- [ ] Browser automation for JavaScript-heavy sites
- [ ] Proxy support for distributed testing
- [ ] Custom success detection patterns
- [ ] Session management
- [ ] Report generation

## Requirements

- Go 1.21 or higher

## Installation

Clone the repository:
```bash
git clone https://github.com/moscovium-mc/blitz
cd blitz
```

### Windows

**Option 1: Using build.bat**
```powershell
# Simply double-click build.bat
```

**Option 2: Manual build**
```powershell
go mod tidy
go build -o blitz.exe .
```

### Linux/macOS

**Option 1: Using setup.sh**
```bash
chmod +x setup.sh
./setup.sh
```

**Option 2: Manual build**
```bash
go mod tidy
go build -o blitz .
```

## How to use it

Basic scan:
```bash
./blitz -url http://example.com/login
```

With all the options:
```bash
./blitz -url http://example.com/login -threads 10 -usernames users.txt -passwords pass.txt -verbose
```

**Options:**

| Option | Description |
|--------|-------------|
| `-url` | Target login page URL (required) |
| `-threads` | Number of concurrent threads (default: 5, max: 20) |
| `-rate` | Rate limit in seconds between requests (default: 1) |
| `-usernames` | Custom username wordlist file |
| `-passwords` | Custom password wordlist file |
| `-verbose` | Show detailed output during scan |

## Examples

**Basic scan with default settings:**
```bash
./blitz -url http://example.com/login
```

**Fast mode (10 threads):**
```bash
./blitz -url http://target.com/login -threads 10
```

**Stealth mode (slow and careful):**
```bash
./blitz -url http://target.com/login -threads 2 -rate 5
```

**Custom wordlists:**
```bash
./blitz -url http://target.com/login -usernames users.txt -passwords pass.txt
```

**Verbose output for detailed analysis:**
```bash
./blitz -url http://target.com/login -verbose
```

**Maximum speed (use with caution):**
```bash
./blitz -url http://target.com/login -threads 20 -rate 0
```

## Features Breakdown

### Smart Form Detection
Blitz automatically analyzes the target page to identify:
- Login forms and input fields
- Hidden fields (CSRF tokens, session IDs)
- Form submission methods (POST/GET)
- Required vs optional fields

### Security Analysis
Built-in security scanner checks for:
- CSRF token implementation
- Clickjacking protection headers
- Cloudflare and WAF presence
- Common SQL injection vulnerabilities

### Performance
- Written in Go for maximum performance
- Concurrent request processing
- Efficient memory usage
- Smart rate limiting to avoid detection

## Version History

### v1.0.0 (Current)

**Core Features:**
- Multi-threaded login bruteforcing
- Smart form and field detection
- CSRF and Clickjacking scanner
- Cloudflare and WAF detection
- SQL injection bypass testing
- Rate limiting protection
- Cross-platform support

**Technical:**
- Go 1.21+ support
- Configurable worker pool (5-20 threads)
- Intelligent success detection
- Built-in wordlists

## Contributing

Contributions are welcome! Here's how you can help:

### Reporting Bugs

If you find a bug, please open an issue with:
- Your operating system and version
- Go version (`go version`)
- Steps to reproduce the bug
- Expected vs actual behavior
- Any error messages

### Pull Requests

Want to contribute code?

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a pull request

**Guidelines:**
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation as needed
- Keep commits focused and descriptive

## Support

If you find this project useful, consider supporting my work:

<a href="https://buymeacoffee.com/webmoney" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" height="40"></a>

**Crypto donations:**
- <a href="bitcoin:bc1quavqz6cxqzfy4qtvq4zxc4fjgap3s7cmxja0k4"><img src="https://img.shields.io/badge/Bitcoin-000000?style=plastic&logo=bitcoin&logoColor=white" alt="Bitcoin"></a> `bc1quavqz6cxqzfy4qtvq4zxc4fjgap3s7cmxja0k4`
- <a href="ethereum:0x5287af72afbc152b09b3bf20af3693157db9e425"><img src="https://img.shields.io/badge/Ethereum-627EEA?style=plastic&logo=ethereum&logoColor=white" alt="Ethereum"></a> `0x5287af72afbc152b09b3bf20af3693157db9e425`
- <a href="solana:HYZjfEx8NbEMJX1vL1GmGj39zA6TgMsHm5KCHWSZxF4j"><img src="https://img.shields.io/badge/Solana-9945FF?style=plastic&logo=solana&logoColor=white" alt="Solana"></a> `HYZjfEx8NbEMJX1vL1GmGj39zA6TgMsHm5KCHWSZxF4j`
- <a href="monero:86zv6vTDuG35sdBzBpwVAsD71hbt2gjH14qiesyrSsMkUAWHQkPZyY9TreeQ5dXRuP57yitP4Yn13SQEcMK4MhtwFzPoRR1"><img src="https://img.shields.io/badge/Monero-FF6600?style=plastic&logo=monero&logoColor=white" alt="Monero"></a> `86zv6vTDuG35sdBzBpwVAsD71hbt2gjH14qiesyrSsMkUAWHQkPZyY9TreeQ5dXRuP57yitP4Yn13SQEcMK4MhtwFzPoRR1`

## Important Legal Notice

> [!WARNING]
> **FOR AUTHORIZED SECURITY TESTING ONLY**

**Only use Blitz on systems you have explicit permission to test.** This tool is designed for ethical security research, authorized penetration testing, and educational purposes only.

**Unauthorized access to computer systems is illegal** and may result in criminal prosecution under various laws including:
- Computer Fraud and Abuse Act (CFAA) in the United States
- Computer Misuse Act in the United Kingdom
- Similar legislation in other jurisdictions

**You are solely responsible for how you use this tool.** The author assumes NO LIABILITY for any misuse, damage, or illegal activity conducted with Blitz.

**Ethical Use Required:**
- Obtain written authorization before testing
- Respect rate limits and system resources
- Follow responsible disclosure practices
- Comply with all applicable laws and regulations

## License

MIT License - See [LICENSE](LICENSE) for details.

---

<div align="center">

**[Star this repo](https://github.com/moscovium-mc/blitz)** if you find it useful

</div>
