# Project Overview

This project is a Go-based REST API for querying UK company data. It provides a fast and efficient way to search for companies within a specific geographic bounding box. The API is built using the Gin framework and uses a SQLite database to store the company and postcode data.

The data is imported from two main sources:

*   **Companies House:** Provides basic company data for UK companies.
*   **Ordnance Survey CodePoint Open:** Provides postcode data for the UK.

The API exposes several endpoints for searching and retrieving company data, as well as for health checks and monitoring.

## Building and Running

The project can be built and run using either Go commands or Docker.

### Go

To build and run the project using Go, you will need to have Go installed on your system.

1.  **Build the application:**

    ```sh
    go build -tags=jsoniter -o company-data .
    ```

2.  **Run the API server:**

    ```sh
    ./company-data api-server --db ./data/companies_data.db --port 8080
    ```

### Docker

To build and run the project using Docker, you will need to have Docker installed on your system.

1.  **Build the Docker image:**

    ```sh
    docker build -t company-data-api .
    ```

2.  **Run the Docker container:**

    ```sh
    docker run -p 8080:8080 -v $PWD/data:/app/data company-data-api api-server
    ```

## Development Conventions

*   **API:** The API is built using the Gin framework.
*   **Database:** The project uses a SQLite database to store the company and postcode data.
*   **Data Import:** The data is imported from CSV files using custom Go scripts.
*   **Dependencies:** The project uses Go modules to manage dependencies.
*   **Linting:** The project uses `golangci-lint` for linting.
*   **Testing:** The project uses `gotestsum` for testing and generates JUnit reports.
*   **CI/CD:** The project uses GitHub Actions for CI/CD. The workflow builds and tests the application, and then publishes a Docker image to GitHub Container Registry.
