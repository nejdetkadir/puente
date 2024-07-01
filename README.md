![Build and test](https://github.com/nejdetkadir/puente/actions/workflows/main.yml/badge.svg?branch=main)
![Go Version](https://img.shields.io/badge/go_version-_1.22.2-007d9c.svg)

# Puente

Puente is a lightweight routing library for AWS Lambda using API Gateway events in Go. It allows you to define routes and handle HTTP methods easily.

## Features

- Define routes for different HTTP methods (GET, POST, PUT, PATCH, DELETE).
- Group routes for better organization.
- Custom error handling.

## Incoming Features 🚀

- Middleware support. (e.g. logging, authentication etc.)
- Lambda function URL request handling natively.

## Installation

To install Puente, run:

```sh
go get github.com/nejdetkadir/puente
```

## Usage

Here's a simple example to get you started:

```go
package main

import (
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/nejdetkadir/puente"
)

func main() {
    lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    router := puente.New()

    router.Get("/hello", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
        return events.APIGatewayProxyResponse{
            StatusCode: 200,
            Body:       "Hello, World!",
        }
    })

    return router.ListenAPIGateway(request), nil
}
```

### Grouping Routes

You can group routes to keep your code organized:

```go
group := router.Group("/api")

group.Get("/users", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       "List of users",
    }
})
```

### Custom Error Handling

You can define custom error handling for your routes:

```go
router.OnError(func(err error) events.APIGatewayProxyResponse {
    return events.APIGatewayProxyResponse{
        StatusCode: 500,
        Body:       "Internal Server Error",
    }
})
```

## Testing

Unit tests are provided to ensure the functionality of the package. To run the tests, use:

```sh
go test ./...
```

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/nejdetkadir/puente. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [code of conduct](https://github.com/nejdetkadir/puente/blob/main/CODE_OF_CONDUCT.md).

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
