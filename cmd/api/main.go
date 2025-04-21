package main

import (
	"context"
	"go-microservice/boot"
)

func main() {
	ctx := context.Background()
	boot.Start(ctx)
}
