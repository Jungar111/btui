/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"btui/cmd"
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	err := cmd.Execute(ctx, cmd.New())
	if err != nil {
		os.Exit(1)
	}
}
