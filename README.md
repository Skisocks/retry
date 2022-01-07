# open-source-archie/retry [![GitLab pipeline](https://img.shields.io/gitlab/pipeline/open-source-archie/retry)](https://gitlab.com/open-source-archie/retry/builds) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/open-source-archie/retry)](https://goreportcard.com/report/gitlab.com/open-source-archie/retry) ![MIT licence](https://img.shields.io/badge/license-MIT-green) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/open-source-archie/retry.svg)](https://pkg.go.dev/gitlab.com/open-source-archie/retry)
A simple Go retry package.

Use the Retry function with a custom backoff policy to re-execute any retryable function with custom parameters.
Alternatively you can use the default policy to re-execute with generic parameters for ease of use.
## Usage
To install:
```
go get gitlab.com/open-source-archie/retry
```
```go
import (
	"gitlab.com/open-source-archie/retry"
)
```
Wrap the retryable function in an anonymous function and then give that as the first argument.

Use a NewBackoffPolicy to retry with default parameters.
```go
retryableFunction := func() error { return SomeReallyCoolFunction() }
if err := retry.Retry(retryableFunction, retry.NewBackoffPolicy); err != nil {
    return err
}
```
Or use a NewCustomBackoffPolicy for a policy with custom parameters.
```go
myBackoffPolicy := retry.NewCustomBackoffPolicy(5, 1000, 2, 1000, 500)
retryableFunction := func() error { return SomeReallyCoolFunction() }
if err := retry.Retry(retryableFunction, myBackoffPolicy); err != nil {
return err
}
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
