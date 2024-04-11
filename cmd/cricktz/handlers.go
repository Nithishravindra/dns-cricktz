package main

import (
	"fmt"

	"github.com/miekg/dns"
)

type handlers struct {
	// services map[string]Service
	domain string
	help   []dns.RR
}

func (h *handlers) handleHelp(w dns.ResponseWriter, r *dns.Msg) {
	m := &dns.Msg{}
	m.SetReply(r)
	m.Compress = false
	m.Answer = h.help
	w.WriteMsg(m)
}

func (h *handlers) handleDefault(w dns.ResponseWriter, m *dns.Msg) {
	respErr(fmt.Errorf(`unknown query. try: dig help @%s`, h.domain), w, m)
	w.WriteMsg(m)
}

// respErr writes an error message to a DNS response.
func respErr(err error, w dns.ResponseWriter, m *dns.Msg) {
	r, err := dns.NewRR(fmt.Sprintf(". 1 IN TXT \"error: %s\"", err.Error()))
	if err != nil {
		lo.Println(err)
		return
	}

	m.Rcode = dns.RcodeServerFailure
	m.Extra = []dns.RR{r}

	w.WriteMsg(m)
}
