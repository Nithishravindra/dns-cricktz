package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"

	"github.com/miekg/dns"
	flag "github.com/spf13/pflag"
)

var (
	lo = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	ko = koanf.New(".")
)

func initConfig() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)

	if err := ko.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		lo.Printf("error reading config: %v", err)
	}
	ko.Load(posflag.Provider(f, ".", ko), nil)
}

func main() {
	lo.Println("initConfig")
	initConfig()

	var (
		h = &handlers{
			// services: make(map[string]Service),
			domain: ko.MustString("server.domain"),
		}
		mux = dns.NewServeMux()

		help = [][]string{}
	)

	help = append(help, []string{"this is the test help message", "dig x @%s"})
	for _, l := range help {
		r, err := dns.NewRR(fmt.Sprintf("help. 1 TXT \"%s\" \"%s\"", l[0], fmt.Sprintf(l[1], h.domain)))
		if err != nil {
			lo.Fatalf("error preparing: %v", err)
		}

		h.help = append(h.help, r)
	}
	mux.HandleFunc("help.", h.handleHelp)

	fmt.Printf("help: %v\n", help)

	// Start the server.
	server := &dns.Server{
		Addr:    ko.MustString("server.address"),
		Net:     "udp",
		Handler: mux,
	}
	lo.Println("listening on ", ko.String("server.address"))
	if err := server.ListenAndServe(); err != nil {
		lo.Fatalf("error starting server: %v", err)
	}
	defer server.Shutdown()
}
