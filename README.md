# Go HTTP Client

This project is a simple HTTP client written in Go that provides methods for making GET, POST, and DELETE requests to an HTTP API.

## Project Structure

```
go-http-client
├── cmd
│   └── main.go          # Entry point of the application
├── pkg
│   ├── client
│   │   └── client.go    # HTTP client implementation
│   └── models
│       └── models.go    # Data models for the API
├── go.mod               # Module definition
└── README.md            # Project documentation
```

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-http-client
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

To run the application, use the following command:
```
go run cmd/main.go
```

You can extend the functionality by modifying the `pkg/client/client.go` file to add more HTTP methods or customize the existing ones.

## Contributing

Feel free to submit issues or pull requests for any improvements or features you would like to see in this project.