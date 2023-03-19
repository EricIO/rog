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
		msgs := m[key]
		sort.Slice(msgs, func(i, j int) bool {
			return msgs[i].Answer[0].Header().Rrtype < msgs[j].Answer[0].Header().Rrtype
		})

		for _, msg := range msgs {
			for _, answer := range msg.Answer {
				fmt.Println(answer.String())
			}

			fmt.Println()
		}
	}
}
