package rog

import (
	"fmt"
	"os"
	"sync"

	"github.com/miekg/dns"
)

type QueryOption struct {
	TCP  bool
	RRs  []string
	Port int
	NS   string
}

type QueryAnswer struct {
	Qtype uint16
	Text  string
}

func Query(host string, qo QueryOption) ([]*dns.Msg, error) {
	answers := make([]*dns.Msg, 0)
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, rr := range qo.RRs {
		qType, ok := dns.StringToType[rr]
		if !ok {
			fmt.Println("invalid RR type")
			os.Exit(1)
		}

		wg.Add(1)
		go func(host string, qtype uint16) {
			defer wg.Done()
			m := new(dns.Msg)
			m.SetQuestion(dns.Fqdn(host), qType)
			answer, err := dns.Exchange(m, "8.8.8.8:53")
			if err != nil {
				panic(err)
			}

			lock.Lock()
			answers = append(answers, answer)
			lock.Unlock()
		}(host, qType)
	}

	wg.Wait()
	return answers, nil
}
