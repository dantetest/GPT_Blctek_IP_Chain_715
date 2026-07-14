package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/data-agent/internal/command"
)

func main() {
	if err := command.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
