# book-stock-manager

## WORK IN PROGRESS

Backend REST API for a book stock manager.

Project for creating a **QR Code-Based Stock and Sales Management System** organized during the *BeSmart* event by the Intermedia Student Activity Unit at Amikom University, Purwokerto.

Documentation for the API can be accessed [here](https://crazydw4rf.github.io/book-stock-manager).

## Usage Guide

### Configuration

1. Copy the example configuration file to the configuration file that will be used
   ```bash
   cp .env.example .env
   ```

   Note: The application can automatically load environment variables without an `.env` file if those variables are already set in the system.

2. Adjust the values in the `.env` file according to your needs
   ```bash
   # Application configuration
   APP_HOST=127.0.0.1
   APP_PORT=8080

   # JWT configuration
   JWT_ACCESS_TOKEN_SECRET=your_secret_key_here
   JWT_REFRESH_TOKEN_SECRET=your_secret_key_here

   # Database configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=username
   DB_PASSWORD=password
   DB_NAME=book_stock
   ```

### Database Migration

Run database migrations to create necessary tables:

```bash
# In the project root directory
go run db/migrate.go db/migrations up
```

### How to Build (Linux)

```bash
# Store module name in a variable for easier use
export MODULE="github.com/crazydw4rf/book-stock-manager"

# Build the application (debug mode)
go build -o book-stock-manager ./cmd/app/main.go

# Build with a custom version
go build -ldflags="-X ${MODULE}/internal/config.APP_VERSION=1.0.0" -o book-stock-manager ./cmd/app/main.go

# Build for production (release mode)
go build -ldflags="-s -w -X ${MODULE}/internal/config.APP_VERSION=1.0.0 -X ${MODULE}/internal/config.APP_ENV=production" -o book-stock-manager ./cmd/app/main.go
```

### How to Run (Linux)

```bash
# Run from the build result
./book-stock-manager

# Or run directly
go run ./cmd/app/main.go
```

After the application is running, the Swagger UI for API documentation can be accessed at: `http://localhost:8080/docs/`

## TODO
- [ ] Fix Dockerfile and compose.yml
- [ ] Add unit and integration tests
- [ ] Implement user authentication
