package main

import (
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func LookupIP(domain string, server string) ([]string, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	msg.RecursionDesired = true

	var err error
	var reply *dns.Msg
	reply, err = dns.Exchange(msg, server)
	if err != nil {
		return nil, err
	}

	// Check for errors in the response
	if reply.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("error in response: %s", reply.String())
	}

	// Extract IP addresses from answer section
	var ips []string
	for _, answer := range reply.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}

	return ips, nil
}

func handleIP(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	result := ""

	lookupIp := getEnv("LOOKUP_SERVER", "8.8.8.8:53")
	res, err := LookupIP(string(m.Payload), lookupIp)
	if err != nil {
		fmt.Println(err)
		result = "NXDOMAIN"
	}
	if len(res) > 0 {
		result = res[0]
	}

	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte(result),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		return res
	}
	return nil
}

func main() {
	mux := coap.NewServeMux()
	mux.Handle("/ip", coap.FuncHandler(handleIP))
	log.Println("Starting CoAP server")
	log.Fatal(coap.ListenAndServe("udp", ":5688", mux))
}
