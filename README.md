# BirdseyeAPI v2

This is a Go+Gin implementation of BirdseyeAPI for scraping and managing news sources.

## Features

- REST API for news management
- Automatic news scraping from multiple sources
- News summarization using OpenAI or Anthropic Claude
- Reaction tracking for news articles

## Setup

### Prerequisites

- Go 1.22 or later
- MySQL database
- OpenAI API key or Anthropic Claude API key (for summarization)
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

# AI Summarization (set at least one)
BIRDSEYEAPI_V2_OPENAI_API_KEY=your_openai_api_key
BIRDSEYEAPI_V2_CLAUDE_API_KEY=your_claude_api_key

# AWS Configuration (for CDN invalidation)
AWS_REGION=your_aws_region
AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID=your_distribution_id
BIRDSEYEAPI_V2_AWS_ACCESS_KEY_ID=your_access_key_id
BIRDSEYEAPI_V2_AWS_SECRET_ACCESS_KEY=your_secret_access_key

# Application Settings
GO_API_PORT=8080
# PRODUCTION runs the pre-built binary; any other value drops into an interactive shell
BIRDSEYEAPI_EXECUTION_MODE=PRODUCTION

# Grafana (observability dashboard)
BIRDSEYEAPI_V2_GRAFANA_ADMIN_PASSWORD=your_grafana_password

# Deployment Settings (used by scrape.sh for remote scrape triggering)
# VENUS_SSH_HOST=your_host
```

> **Note:** `OPENAI_MODEL`, `OPENAI_CHAT_ENDPOINT`, and `SELENIUM_URL` are not configurable via environment variables. The OpenAI model (`gpt-4.1-mini`) and endpoint, the Claude model (`claude-3-5-sonnet-20241022`) and endpoint, and the Selenium URL (`http://selenium:4444/wd/hub`) are all hardcoded in the source. `SELENIUM_URL` is set automatically by Docker Compose and does not need to be set manually.

### Build and Run

#### Locally

```bash
# Download dependencies
go mod download

# Run the application
./run.sh
# or directly:
go run ./go/src/main.go
```

#### Using Docker Compose (Recommended)

The Go container is started automatically by `go-entrypoint.sh`. When `BIRDSEYEAPI_EXECUTION_MODE=PRODUCTION`, it runs the pre-built binary at `go/dist/birdseyeapi_v2`; otherwise it drops into an interactive shell for development.

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f go
```

#### Building for Production

```bash
# Build the binary (output: go/dist/birdseyeapi_v2)
./build.sh
```

After building, restart the Go container to pick up the new binary:

```bash
docker compose restart go
```

#### Triggering a Remote Scrape

`scrape.sh` SSHes into the configured host and calls the scrape endpoint:

```bash
./scrape.sh
```

This requires `VENUS_SSH_HOST` to be set in `.env`.

## API Endpoints

### News

- `GET /news/today-news` — Get news articles for today (falls back up to 10 days if today has none)
- `GET /news/news-reactions/:news-id` — Get Hatena Bookmark reactions for a specific article
- `POST /news/scrape` — Trigger news and reaction scraping (runs in background; returns 409 if already running)
- `GET /news/trends` — Get trending topics from Google Trends

## Data Sources

The API scrapes news from the following sources (up to 15 articles each):

- CloudWatch by Impress
- Hatena Bookmark hot entries (IT)
- Zenn (daily tech articles)
- ZDNet Japan

News articles are automatically summarized using the OpenAI API (`gpt-4.1-mini`) by default. The summary is limited to 200 characters in Japanese with appropriate line breaks. A Claude (`claude-3-5-sonnet-20241022`) summarizer is also available in the codebase.

## Reaction Scraping

The system scrapes reactions for news articles from:

- Hatena Bookmark comments (via Selenium/Firefox)

Scraping is serialized: concurrent `POST /news/scrape` requests are rejected with HTTP 409 to prevent multiple Selenium sessions from running simultaneously.

## Technical Architecture

| Component | Image | Port (host:container) | Role |
|---|---|---|---|
| Nginx | nginx:1.31.2 | 1111:1111 | Reverse proxy, rate limiting |
| Go API | golang:1.24 | 8080:8080 | Application server |
| MySQL | mysql:9.3 | 3307:3306 | Data persistence |
| Selenium | selenium/standalone-firefox:133.0 | 4444:4444 | Headless Firefox for reaction scraping |
| Loki | grafana/loki:3.3.2 | 3100:3100 | Log aggregation |
| Promtail | grafana/promtail:3.3.2 | — | Log shipping to Loki |
| Grafana | grafana/grafana:11.4.0 | 3000:3000 | Log and metrics dashboard |

The MySQL database (`birds_eye`) is created automatically on first boot via `mysql/create_db.sql`.

## Health Check

The `/HealthCheck` endpoint is handled directly by Nginx (not the Go application) and returns HTTP 200 with `ok`.
