# Route53 DNS Records Reader

A Go application that retrieves all DNS records from AWS Route53.

## Prerequisites

- Go 1.21 or later
- AWS credentials configured (via environment variables, AWS credentials file, or IAM role)
- IAM permissions for Route53:
  - `route53:ListHostedZones`
  - `route53:ListResourceRecordSets`

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Ensure AWS credentials are configured (the SDK will automatically use the default credential chain):
   - AWS credentials file (~/.aws/credentials)
   - Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
   - IAM role (if running on EC2/ECS/Lambda)
   - SSO credentials

## Usage

### List all DNS records
```bash
go run main.go
```

### Create a subdomain
Edit `main.go` and uncomment the `createSubdomain` example, then update:
- `hostedZoneID`: Your Route53 hosted zone ID (e.g., "Z1234567890ABC")
- `subdomain`: Full subdomain name (e.g., "test123.example.com")
- `ipAddress`: IP address for the A record (e.g., "192.0.2.1")
- `ttl`: Time to live in seconds (e.g., 300)

### Delete a subdomain
Edit `main.go` and uncomment the `deleteSubdomain` example with the same parameters used during creation.

## Features

The application provides:
1. **List all DNS records** - Display all hosted zones and their DNS records
2. **Create subdomain** - Add a new A record subdomain to any hosted zone
3. **Delete subdomain** - Remove an existing A record subdomain

## Functions

- `getAllZonesAndRecords()` - Retrieves and displays all zones and records
- `listHostedZones()` - Gets all hosted zones with pagination
- `listResourceRecordSets()` - Gets all DNS records for a zone with pagination
- `createSubdomain()` - Creates a new A record subdomain
- `deleteSubdomain()` - Deletes an existing A record subdomain
- `displayZoneRecords()` - Displays records for a specific zone
- `displayRecord()` - Formats and prints a single DNS record

## Example Output

```
Found 2 hosted zones

Hosted Zone: example.com. (ID: /hostedzone/Z1234567890ABC)
----------------------------------------
  Name: example.com.
  Type: A
  TTL: 300
    Value: 192.0.2.1

  Name: www.example.com.
  Type: CNAME
  TTL: 300
    Value: example.com.
```

## Disclaimer

This is a sample application for educational and testing purposes only. It should not be used in production environments without proper error handling, validation, logging, and security measures. Always follow AWS best practices and your organization's security policies when working with DNS records in production.
