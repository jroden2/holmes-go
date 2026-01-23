# Holmes-Go 
## Local Diff Checker

![Go Version](https://img.shields.io/badge/Go-1.25.1-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Active-success)

**Holmes** is a lightweight, security-first diff checker written in **Go**, designed for comparing content **locally** without relying on third-party or online services.

Holmes focuses on privacy and safety by ensuring all comparisons run entirely on your machine.

---

## Why Holmes?

Many quick diff tools require pasting files or data into online services, which can be risky when working with sensitive or proprietary information.

Holmes is built with a security-first mindset:

* Runs 100% locally
* No external services or network calls
* No uploads or data sharing
* Suitable for confidential configs, logs, API responses, and internal data

---

## Features

### Text Comparison

Compare plain text files and clearly identify differences.

### JSON Comparison

Compare JSON files in a structure-aware way:

* Ignores formatting and key ordering differences
* Highlights actual data changes

### XML Comparison

Exactly like JSON, you can compare XML files in a structure-aware way:
* Ignores formatting and key ordering differences
* Highlights actual data changes

### UI-Based

Holmes includes a user interface, making it easy to visualise differences without relying on command-line workflows.

---

## Installation

### Build from source

```bash
git clone https://github.com/jroden2/holmes-go.git
cd holmes-go
go build ./cmd/*
```

---

## Security-First by Design

Holmes is intentionally designed to keep all comparisons on your machine.

* No telemetry
* No background services
* No outbound network traffic required

If you are diffing data that should not be uploaded to the internet, Holmes is built for that use case.

---

## Example Use Cases

* Comparing configuration files across environments (e.g. development vs production)
* Reviewing JSON API response changes
* Debugging differences in structured output
* Verifying log output changes during development

---

## Contributing

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Submit a pull request

Bug reports and feature requests can be submitted via GitHub issues.

---

## License

This project is licensed under the MIT License.
