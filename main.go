package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/m0ush/rps/app"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	fnamePtr := flag.String("fin", "trundl/trundl.json", "The file name to parse")
	flag.Parse()

	f, err := os.Open(*fnamePtr)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := app.Runner(f); err != nil {
		return err
	}

	return nil
}
