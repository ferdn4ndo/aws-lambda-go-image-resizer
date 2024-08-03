package main

import (
	"flag"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ferdn4ndo/aws-lambda-go-image-resizer/handlers"
)

var (
	selfTest   = flag.Bool("self-test", false, "Execute a self-test (convert a fixture image)")
	runFromCli = flag.Bool("run-from-cli", false, "Execute an image conversion from the command line (providing input and output file)")
)

func main() {
	lambda.Start(new(handlers.GatewayHandler).ServeHTTP)
}
