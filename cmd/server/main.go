package main

import (
	"os"

	"openapi-mcp/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args))
}
