# BirdseyeAPI v2

This is a Go+Gin implementation of BirdseyeAPI for scraping and managing news sources.

## Features

- REST API for news management
- Automatic news scraping from multiple sources
- News summarization using OpenAI
- Reaction tracking for news articles

## Setup

### Prerequisites

- Go 1.20 or later
- MySQL database
- OpenAI API key (for summarization features)
- Selenium (for reaction scraping)

### Environment Variables

A sample environment file `.env.example` is provided in the repository. Copy this file to `.env` and update the values accordingly.

```bash
# Create your .env file from the example
cp .env.example .env
# Edit the file with your values
nano .env
```

The following environment variables are required:

```
# Database Configuration
MYSQL_ROOT_PASSWORD=your_secure_password
# DB connection settings are hardcoded in db_info.go:
# Username: root
# Host: mysql
# Port: 3306
# DBName: birds_eye

# OpenAI Configuration
OPENAI_MODEL=gpt-4-turbo
OPENAI_CHAT_ENDPOINT=https://api.openai.com/v1/chat/completions

# AWS Configuration (for CDN invalidation)
AWS_REGION=your_aws_region
AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID=your_distribution_id

# Application Settings
GO_API_PORT=8080
BIRDSEYEAPI_EXECUTION_MODE=production
SCRAPING_ARTICLES=10

# Selenium Settings (for reaction scraping with Hatena)
SELENIUM_URL=http://selenium:4444/wd/hub

# Deployment Settings (used by deploy.sh)
# VENUS_SSH_HOST=your_host
# VENUS_HOME=/path/to/deployment
```

### Build and Run

#### Locally

```bash
# Download dependencies
go mod download

# Run the application
go run ./go/src/main.go
```

#### Using Docker Compose (Recommended)

```bash
# Start all services (MySQL, Go, Nginx, Selenium)
docker compose up -d

# Run the application inside the container
docker compose exec go go run ./go/src/main.go

# Or use the convenience script
./run.sh
```

#### Building for Production

```bash
# Build the binary
./build.sh

# Deploy to server
./deploy.sh
```

## API Endpoints

### News

- `GET /news/today-news`: Get all news articles for today
- `GET /news/news-reactions/:news-id`: Get reactions for a specific news article
- `POST /news/scrape`: Trigger news scraping
- `GET /news/trends`: Get trending topics from Google Trends

## Data Sources

The API scrapes news from the following sources:

- CloudWatch by Impress
- Hatena
- Zenn
- ZDNet Japan

News articles are automatically summarized using the OpenAI API. The summary is limited to 200 characters in Japanese with proper line breaks for readability.

## Reaction Scraping

The system also scrapes reactions (comments) for news articles from:

- Hatena Bookmark comments

## Technical Architecture

- **Frontend Proxy**: Nginx serves as a reverse proxy on port 1111
- **API Server**: Go+Gin running on port 8080
- **Database**: MySQL 8.4 for data persistence
- **Selenium**: Firefox headless browser for scraping reactions
- **CDN Invalidation**: AWS CloudFront integration for cache management

## Health Check

A health check endpoint is available at `/HealthCheck` that returns HTTP 200 with 'ok' message.
