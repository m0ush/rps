package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/m0ush/rps/app"
	"github.com/m0ush/rps/trundl"
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

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var rs trundl.Records
	if err := json.Unmarshal(bs, &rs); err != nil {
		return err
	}

	lims, err := app.RegisterLimits()
	if err != nil {
		return err
	}

	// TODO: print or return alerts to be picked up by notification process
	as := app.Sieve(rs, lims...)
	app.AddAlerts(as)

	return nil
}
