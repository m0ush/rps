package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/m0ush/rps/rpsdb"
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

	x := trundl.Pipeline(rs)
	fmt.Println(x)

	return nil
}

func run1() error {
	db, err := rpsdb.Open("sqlite", "./rpsdb/rps.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// secs, err := db.AllSecurities()
	// if err != nil {
	// 	return err
	// }
	// for _, p := range secs {
	// 	fmt.Println(p)
	// }

	lims, err := db.AllLimits()
	if err != nil {
		return err
	}
	for _, l := range lims {
		fmt.Println(l)
	}
	return nil
}
