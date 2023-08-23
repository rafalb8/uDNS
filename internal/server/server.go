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

func Start(port, httpPort, dir string, noResolv bool) error {
	if dir == "" {
		dir = internal.ConfigPath
	}
	os.MkdirAll(dir, os.ModePerm)

	// Check if file exists, if not create it
	file := path.Join(dir, "hosts")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			return err
		}
	}

	if !noResolv {
		internal.ModifyResolv()
	}

	// Load hosts
	refresh(dir)

	// track changes
	close := internal.Track(dir, func() { refresh(dir) })
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
				ip, ok := hostMapping[q.Name]
				if !ok {
					continue
				}
				fmt.Printf("Query for %s\n", q.Name)

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
