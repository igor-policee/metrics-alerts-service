# Project Structure

```markdown
metrics-alerts-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── counter.go
│   │   └── gauge.go
│   ├── metrics/
│   │   ├── counter.go
│   │   └── gauge.go
│   └── utils/
│       └── helpers.go
├── pkg/
│   └── // (optional for reusable packages)
├── go.mod
└── go.sum
```

## Explanation:
- `cmd/server/main.go`: Entry point for the server application.
- `internal/handlers/`: Contains HTTP handler functions for different metric types.
- `internal/metrics/`: Manages the business logic and data storage for metrics.
- `internal/utils/`: Holds utility functions shared across the application, such as `splitPath`.
- `pkg/`: (Optional) For reusable packages that can be imported by external applications.
- `go.mod` & `go.sum`: Manage project dependencies.