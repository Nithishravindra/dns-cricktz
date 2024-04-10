package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var lo = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func initConfig() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Bool("version", false, "show build version")
}

func main() {
	lo.Println("I'm trying to exploit dns protocol")
	initConfig()
}
