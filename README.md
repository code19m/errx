# `errx`: Advanced Error Handling for Go

`errx` is a flexible error handling package designed to provide structured, extensible, and developer-friendly error management in Go projects. It extends Go's built-in `error` interface with additional methods, supports gRPC integration, and includes utilities for detailed error tracing and contextual debugging.

## Features

- **Extended Error Interface**: Add structured details, validation fields, and trace information to errors.
- **gRPC Compatibility**: Convert errors to/from gRPC status codes seamlessly.
- **Error Tracing**: Automatically track the origin and flow of errors through the system.
- **Customizable Options**: Use functional options to customize errors on creation or wrapping.
- **Rich Metadata**: Attach contextual information and validation fields for debugging and logging.
- **Integration Utilities**: Utilities for extracting or converting errors with functions like `AsErrorX`, `GetCode`, and `GetType`.

## Installation

```bash
go get github.com/code19m/errx
```

## Usage

### 0. Helper function for use in examples
```go
func logError(err error) {
	if err == nil {
		return
	}
	e := errx.AsErrorX(err)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Error("Error occurred",
		slog.String("err_code", e.Code()),
		slog.String("err_type", e.Type().String()),
		slog.String("err_message", e.Error()),
		slog.String("err_trace", e.Trace()),
		slog.Any("err_fields", e.Fields()),
		slog.Any("err_details", e.Details()),
	)
}
```

### 1. Creating Custom Errors

```go
package main

import (
	"log/slog"
	"os"

	"github.com/code19m/errx"
)

func main() {
	err := errx.New("Resource not found",
		errx.WithCode("NOT_FOUND"),                     // Set error code
		errx.WithType(errx.T_NotFound),                 // Set error type
		errx.WithTrace(),                               // Add trace
		errx.WithDetails(errx.M{"resource_id": "123"}), // Add additional details
		errx.WithFields(errx.M{"username": "invalid"}), // Add validation fields
	)

	logError(err)
}
```

**Output**

```json
{
  "time": "2024-12-08T14:36:34.715364+05:00",
  "level": "ERROR",
  "msg": "Error occurred",
  "err_code": "NOT_FOUND",
  "err_type": "T_NotFound",
  "err_message": "[T_NotFound: NOT_FOUND] - Resource not found",
  "err_trace": "[main.go:28] main.main",
  "err_fields": {
    "username": "invalid"
  },
  "err_details": {
    "resource_id": "123"
  }
}
```

---


### 2. Wrapping Existing Errors

```go
package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/code19m/errx"
)

func main() {
	baseErr := errors.New("database connection failed")
	err := errx.Wrap(
		baseErr,
		errx.WithDetails(errx.M{
			"host": "localhost",
			"port": "5432",
		}),
	)

	logError(err)
}
```

**Output**

```json
{
  "time": "2024-12-08T14:42:44.488966+05:00",
  "level": "ERROR",
  "msg": "Error occurred",
  "err_code": "INTERNAL",
  "err_type": "T_Internal",
  "err_message": "[T_Internal: INTERNAL] - database connection failed",
  "err_trace": "[main.go:30] main.main",
  "err_fields": null,
  "err_details": {
    "host": "localhost",
    "port": "5432"
  }
}
```

---

### 3. Error trace information

```go
package main

import (
	"log/slog"
	"os"

	"github.com/code19m/errx"
)

func main() {
	err := errx.Wrap(firstFunc())

	logError(err)
}

func firstFunc() error {
	return errx.Wrap(secondFunc())
}

func secondFunc() error {
	return errx.New("some error occurred", errx.WithTrace())
}
```

**Output**

```json
{
  "time": "2024-12-08T14:49:21.362634+05:00",
  "level": "ERROR",
  "msg": "Error occurred",
  "err_code": "INTERNAL",
  "err_type": "T_Internal",
  "err_message": "[T_Internal: INTERNAL] - some error occurred",
  "err_trace": "[main.go:28] main.main ➡️ [main.go:34] main.firstFunc ➡️ [main.go:38] main.secondFunc",
  "err_fields": null,
  "err_details": null
}
```

---

## Error Types

The package defines several error types for categorizing errors:

| Type                | Description                          |
|---------------------|--------------------------------------|
| `T_Internal`        | Internal server errors               |
| `T_Validation`      | Input validation errors              |
| `T_NotFound`        | Resource not found errors            |
| `T_Conflict`        | Conflicting resource errors          |
| `T_Authentication`  | Authentication-related errors        |
| `T_Forbidden`       | Permission-related errors            |

## Functional Options

| Option              | Description                          |
|---------------------|--------------------------------------|
| `WithCode`          | Sets a machine-readable error code   |
| `WithType`          | Sets the error type                  |
| `WithTrace`         | Adds trace information               |
| `WithPrefix`        | Adds a prefix to trace and details   |
| `WithDetails`       | Adds debugging details               |
| `WithFields`        | Sets validation-related fields       |


## Testing

Unit tests cover the package functionality. To run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
