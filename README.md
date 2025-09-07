# Company Data API

A fast REST API for querying UK company data by geographic bounding box, built with Go, SQLite, and Gin. It imports official datasets from Companies House and Ordnance Survey CodePoint Open, providing spatial search capabilities for company records.

## Example Usage

#### Search for companies within a bounding box:

```http
GET /v1/company-data/search?bbox=425000,450000,435000,460000
```

Example JSON response (see [models/company_data.go](./models/company_data.go) for the exact schema):

```json
{
    "results": [
        {
            "company_number": "12345678",
            "company_name": "ACME WIDGETS LIMITED",
            "reg_address_post_code": "AB12 3CD",
            "easting": 426000,
            "northing": 451000
        },
        {
            "company_number": "87654321",
            "company_name": "BETA SERVICES LLP",
            "reg_address_post_code": "AB12 3CD",
            "easting": 427500,
            "northing": 452500
        }
    ],
    "attribution": [
        "Basic Company Data (UK Gov, Companies House), https://download.companieshouse.gov.uk/en_output",
        "CodePoint Open (UK Gov, OS Data Hub), https://osdatahub.os.uk/downloads/open/CodePointOpen"
    ],
    "last_updated": "2025-06-30T00:00:00Z"
}
```

#### Group companies by postcode within a bounding box:

```http
GET /v1/company-data/search/by-postcode?bbox=425000,450000,435000,460000
```

The JSON response is similar to previously, but results are grouped by postcode.

#### Health check:

```http
GET /healthz
```

#### Prometheus metrics:

```http
GET /metrics
```

#### Swagger/OpenAPI documentation:

```http
GET /swagger/index.html
```

OpenAPI documentation is auto-generated from annotated Go handler functions using [swaggo/swag](https://github.com/swaggo/swag). See `internal/search.go` for endpoint annotations.

## Architecture Overview

```mermaid
flowchart TD
    subgraph CLI
        X[company-data]
        X1[api-server]
        X2[import]
        X --> X1
        X --> X2
    end

    subgraph Data Import
        A[Zip Files: Companies House, CodePoint Open]
        B[Import Scripts]
        C[SQLite DB]
        X2 --> B
        A --> B --> C
    end

    subgraph API Server
        D[main.go]
        E[internal/search.go]
        F[repositories/search.go]
        G[models/company_data.go]
        X1 --> D
        C --> F
        F --> E
        E --> D
    end

    subgraph Client
        H[HTTP Client]
        H --> D
    end
```

-   **Data Import:**
    -   `internal/company_data.go` and `internal/code_point.go` handle importing zipped CSV data into SQLite.
-   **Database:**
    -   `internal/migration.sql` defines the schema for company and postcode data.
-   **API Server:**
    -   `main.go` sets up the Gin HTTP server, routes, and middleware.
    -   `internal/search.go` implements search endpoints.
    -   `repositories/search.go` performs SQL queries joining company and postcode tables.
-   **Models:**
    -   `models/company_data.go` defines the data structures returned by the API.

## Getting Started Locally

### CLI Commands

The application uses [cobra](https://github.com/spf13/cobra) for its command-line interface. The main commands are:

-   `api-server` — Starts the HTTP API server.

    -   Options:
        -   `--db <path>`: Path to Companies data SQLite database (default: `./data/companies_data.db`)
        -   `--port <port>`: Port to run HTTP server on (default: `8080`)
        -   `--debug`: Enable debugging (pprof). **Warning:** Do not enable in production.

-   `import` — Imports Companies House ZIP file into the database.
    -   Options:
        -   `--zip-file <path>`: Path to Companies House .zip file (default: `./data/BasicCompanyDataAsOneFile-2025-07-01.zip`)

Example usage:

```sh
./company-data api-server --db ./data/companies_data.db --port 8080
./company-data import --zip-file ./data/BasicCompanyDataAsOneFile-2025-07-01.zip
```

### 1. Download Data

-   Download company data from [Companies House](https://download.companieshouse.gov.uk/en_output.html)
-   Download CodePoint Open from [OS Data Hub](https://osdatahub.os.uk/downloads/open/CodePointOpen)
-   Place both `.zip` files in the `./data` directory (do **not** decompress).

### 2. Regenerate Swagger definitions

Swagger/OpenAPI docs are generated from code comments. To update the docs after changing endpoints or annotations:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 3. Build and Run

```sh
go build -tags=jsoniter -o company-data .
./company-data api-server --db ./data/companies_data.db --port 8080
```

## Using Docker

Build the image:

```sh
docker build -t company-data-api .
```

Run the container (mount your data directory):

```sh
docker run -p 8080:8080 -v $PWD/data:/app/data company-data-api http
```

-   The binary is built with the `jsoniter` tag for fast JSON serialization.
-   The container runs as a non-root user for security.
-   Health checks are enabled on `/healthz`.
-   Timezone and CA certificates are included for compatibility.

## API Endpoints

| Endpoint                                       | Description                                   |
| ---------------------------------------------- | --------------------------------------------- |
| `/v1/company-data/search?bbox=...`             | Search companies within a bounding box        |
| `/v1/company-data/search/by-postcode?bbox=...` | Group companies by postcode in a bounding box |
| `/healthz`                                     | Health check                                  |
| `/metrics`                                     | Prometheus metrics                            |
| `/swagger/index.html`                          | Swagger UI (OpenAPI documentation)            |
| `/swagger/doc.json`                            | OpenAPI definition (JSON)                     |

## Attribution

-   Basic Company Data (UK Gov, Companies House): https://download.companieshouse.gov.uk/en_output.html
-   CodePoint Open (UK Gov, OS Data Hub): https://osdatahub.os.uk/downloads/open/CodePointOpen

## TODO & Future Enhancements

-   [ ] Add authentication and rate limiting
-   [ ] Support for additional spatial queries (e.g., radius search)
-   [ ] Pagination and filtering options
-   [ ] Docker Compose for easier setup
-   [ ] Automated data refresh/import
-   [X] OpenAPI/Swagger documentation (auto-generated from code)
-   [ ] More robust error handling and logging
-   [ ] Unit and integration tests for import and API layers

## License

See `LICENSE.md` for further details.
