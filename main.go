package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	var pkt trundl.Packet
	if err := json.Unmarshal(bs, &pkt); err != nil {
		return err
	}

	pkt.DataFrame()
	pkt.FindAlerts()

	jsn, err := json.MarshalIndent(pkt.Alerts, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsn))
	return nil
}
