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

### Environment Variables

Set the following environment variables to configure the application:

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name
OPENAI_API_KEY=your_openai_key
PORT=8080
```

### Build and Run

#### Locally

```bash
# Download dependencies
go mod download

# Run the application
go run ./cmd/server/main.go
```

#### Using Docker

```bash
# Build the Docker image
docker build -t birdseyeapi_v2 .

# Run the Docker container
docker run -p 8080:8080 --env-file .env birdseyeapi_v2
```

## API Endpoints

### News

- `GET /api/news/`: Get all news articles
- `GET /api/news/:id`: Get a news article by ID
- `POST /api/news/`: Create a new news article
- `POST /api/news/scrape`: Trigger news scraping
- `POST /api/news/summarize`: Summarize a news article

## Data Sources

The API scrapes news from the following sources:

- CloudWatch by Impress
- Hatena
- Zenn
- ZDNet Japan
