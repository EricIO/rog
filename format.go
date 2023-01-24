package rog

import (
	"fmt"
	"sort"

	"github.com/miekg/dns"
)

func FormatOutput(m map[string][]*dns.Msg) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, msg := range m[key] {
			for _, answer := range msg.Answer {
				// TODO: we also want to sort by RR type
				fmt.Println(answer.String())
			}
		}

		fmt.Println()
	}
}
