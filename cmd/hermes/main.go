package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rugwirobaker/hermes/build"
)

func main() {
	log.SetFlags(0)

	if err := run(context.Background(), os.Args[1:]); err == flag.ErrHelp {
		os.Exit(2)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) (err error) {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	switch cmd {
	case "serve":
		return runServe(ctx, args)
	case "migrate":
		return runMigrate(ctx, args)
	case "version":
		fmt.Printf("version: %s, date: %s\n", build.Info().Version, build.Info().Date)
		return
	default:
		if cmd == "" || cmd == "help" || strings.HasPrefix(cmd, "-") {
			printUsage()
			return flag.ErrHelp
		}
		return fmt.Errorf("litefs %s: unknown command", cmd)
	}
}

func printUsage() {
	fmt.Println(`
hermes is a sms gateway for other applications in the network.

Usage:
	hermes [command]

Available Commands:
	server	  start the hermes server
	migrate	  run database migrations
	version	  print the version of hermes
`[1:])
}
