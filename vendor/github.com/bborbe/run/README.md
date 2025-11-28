# Run

[![Go Reference](https://pkg.go.dev/badge/github.com/bborbe/run.svg)](https://pkg.go.dev/github.com/bborbe/run)
[![Go Report Card](https://goreportcard.com/badge/github.com/bborbe/run)](https://goreportcard.com/report/github.com/bborbe/run)

A Go library for parallel function execution with different error handling strategies and utilities.

## Features

- **Parallel Execution**: Run multiple functions concurrently with various error handling strategies
- **Error Handling**: Choose from different error handling approaches (fail-fast, collect all, etc.)
- **Retry Logic**: Built-in retry functionality with configurable backoff
- **Context Support**: Full context.Context support for cancellation and timeouts  
- **Utilities**: Additional helpers for delayed execution, parallel skipping, and more

## Installation

```bash
go get github.com/bborbe/run
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/bborbe/run"
)

func main() {
    ctx := context.Background()
    
    // Run functions in parallel, stop on first error
    err := run.CancelOnFirstError(ctx,
        func(ctx context.Context) error {
            fmt.Println("Task 1")
            return nil
        },
        func(ctx context.Context) error {
            fmt.Println("Task 2") 
            return nil
        },
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Execution Strategies

### Parallel Execution with Error Handling

```go
// Cancel all remaining functions when first one finishes
err := run.CancelOnFirstFinish(ctx, funcs...)

// Cancel all remaining functions when first error occurs
err := run.CancelOnFirstError(ctx, funcs...)

// Run all functions, collect all errors
err := run.All(ctx, funcs...)

// Run functions sequentially
err := run.Sequential(ctx, funcs...)
```

### Retry with Backoff

```go
retryableFunc := run.Retry(run.Backoff{
    Delay:   time.Second,
    Factor:  2.0,
    Retries: 3,
    IsRetryAble: func(err error) bool {
        // Custom logic to determine if error is retryable
        return true
    },
}, func(ctx context.Context) error {
    // Your function that may fail
    return someOperation(ctx)
})

err := retryableFunc(ctx)
```

### Delayed Execution

```go
// Delay execution by 5 seconds
delayedFunc := run.Delayed(func(ctx context.Context) error {
    fmt.Println("This runs after 5 seconds")
    return nil
}, 5*time.Second)

err := delayedFunc(ctx)
```

### Prevent Parallel Execution

```go
skipper := run.NewParallelSkipper()

// Wrap function to prevent parallel execution
protectedFunc := skipper.SkipParallel(func(ctx context.Context) error {
    // This function will be skipped if already running
    return expensiveOperation(ctx)
})
```

## Core Types

The library is built around two main interfaces:

```go
// Func is a function that can be executed with context
type Func func(context.Context) error

// Runnable interface for objects that can be run
type Runnable interface {
    Run(ctx context.Context) error
}
```

## Advanced Usage

### Custom Error Handling

```go
// Get a channel of errors for custom processing
errorChan := run.Run(ctx, funcs...)
for err := range errorChan {
    if err != nil {
        log.Printf("Function failed: %v", err)
    }
}
```

### Background Runners

```go
// For long-running background tasks
runner := run.NewBackgroundRunner()
// ... (see source code for full API)
```

## Examples

### Web Server with Graceful Shutdown

```go
func startServer(ctx context.Context) error {
    server := &http.Server{Addr: ":8080"}
    
    return run.CancelOnFirstFinish(ctx,
        func(ctx context.Context) error {
            // Start server
            return server.ListenAndServe()
        },
        func(ctx context.Context) error {
            // Wait for context cancellation
            <-ctx.Done()
            // Graceful shutdown
            return server.Shutdown(context.Background())
        },
    )
}
```

### Parallel Data Processing

```go
func processData(ctx context.Context, items []string) error {
    var funcs []run.Func
    
    for _, item := range items {
        item := item // capture loop variable
        funcs = append(funcs, func(ctx context.Context) error {
            return processItem(ctx, item)
        })
    }
    
    // Process all items, collect all errors
    return run.All(ctx, funcs...)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the BSD-style license. See the LICENSE file for details.
