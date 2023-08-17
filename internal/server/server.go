package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/miekg/dns"
	"github.com/rafalb8/uDNS/internal"
)

var (
	// Mapping between ip addresses and hostnames
	hostMapping = map[string]string{}
	mu          sync.Mutex
)

func Start(port, httpPort, path string) error {
	if path == "" {
		path = internal.ConfigPath
	}
	os.MkdirAll(path, os.ModePerm)

	// Load hosts
	refresh(path)

	// track changes
	close := internal.Track(path, func() { refresh(path) })
	defer close()

	// Start http server
	go http.ListenAndServe("127.0.0.1:"+httpPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Hello from µDNS\n")))
	}))

	// Start DNS server
	return dns.ListenAndServe("127.0.0.1:"+port, "udp", dns.HandlerFunc(dnsRequest))
}

func dnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range m.Question {
			switch q.Qtype {
			case dns.TypeA:
				fmt.Printf("Query for %s\n", q.Name)
				ip, ok := hostMapping[q.Name]
				if !ok {
					continue
				}

				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}

	w.WriteMsg(m)
}

// Refresh hosts
func refresh(dir string) {
	mu.Lock()
	defer mu.Unlock()

	hostMapping = map[string]string{
		"udns.local.": "127.0.0.1",
	}

	content, err := os.ReadFile(path.Join(dir, "hosts"))
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		ip := parts[0]
		hosts := parts[1:]

		for _, host := range hosts {
			host = strings.TrimSpace(host)
			hostMapping[host+"."] = ip
		}
	}
	fmt.Println("refreshed", hostMapping)
}
