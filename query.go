package rog

import (
	"fmt"
	"os"
	"sync"

	"github.com/miekg/dns"
)

var TypeToInt map[string]uint16 = map[string]uint16{
	"None":       0,
	"A":          1,
	"NS":         2,
	"MD":         3,
	"MF":         4,
	"CNAME":      5,
	"SOA":        6,
	"MB":         7,
	"MG":         8,
	"MR":         9,
	"NULL":       10,
	"PTR":        12,
	"HINFO":      13,
	"MINFO":      14,
	"MX":         15,
	"TXT":        16,
	"RP":         17,
	"AFSDB":      18,
	"X25":        19,
	"ISDN":       20,
	"RT":         21,
	"NSAPPTR":    23,
	"SIG":        24,
	"KEY":        25,
	"PX":         26,
	"GPOS":       27,
	"AAAA":       28,
	"LOC":        29,
	"NXT":        30,
	"EID":        31,
	"NIMLOC":     32,
	"SRV":        33,
	"ATMA":       34,
	"NAPTR":      35,
	"KX":         36,
	"CERT":       37,
	"DNAME":      39,
	"OPT":        41,
	"APL":        42,
	"DS":         43,
	"SSHFP":      44,
	"RRSIG":      46,
	"NSEC":       47,
	"DNSKEY":     48,
	"DHCID":      49,
	"NSEC3":      50,
	"NSEC3PARAM": 51,
	"TLSA":       52,
	"SMIMEA":     53,
	"HIP":        55,
	"NINFO":      56,
	"RKEY":       57,
	"TALINK":     58,
	"CDS":        59,
	"CDNSKEY":    60,
	"OPENPGPKEY": 61,
	"CSYNC":      62,
	"ZONEMD":     63,
	"SVCB":       64,
	"HTTPS":      65,
	"SPF":        99,
	"UINFO":      100,
	"UID":        101,
	"GID":        102,
	"UNSPEC":     103,
	"NID":        104,
	"L32":        105,
	"L64":        106,
	"LP":         107,
	"EUI48":      108,
	"EUI64":      109,
	"URI":        256,
	"CAA":        257,
	"AVC":        258,
	"TKEY":       249,
	"TSIG":       250,
	"IXFR":       251,
	"AXFR":       252,
	"MAILB":      253,
	"MAILA":      254,
	"ANY":        255,
	"TA":         32768,
	"DLV":        32769,
	"Reserved":   65535,
}

var IntToType map[uint16]string = map[uint16]string{
	0:     "None",
	1:     "A",
	2:     "NS",
	3:     "MD",
	4:     "MF",
	5:     "CNAME",
	6:     "SOA",
	7:     "MB",
	8:     "MG",
	9:     "MR",
	10:    "NULL",
	12:    "PTR",
	13:    "HINFO",
	14:    "MINFO",
	15:    "MX",
	16:    "TXT",
	17:    "RP",
	18:    "AFSDB",
	19:    "X25",
	20:    "ISDN",
	21:    "RT",
	23:    "NSAPPTR",
	24:    "SIG",
	25:    "KEY",
	26:    "PX",
	27:    "GPOS",
	28:    "AAAA",
	29:    "LOC",
	30:    "NXT",
	31:    "EID",
	32:    "NIMLOC",
	33:    "SRV",
	34:    "ATMA",
	35:    "NAPTR",
	36:    "KX",
	37:    "CERT",
	39:    "DNAME",
	41:    "OPT",
	42:    "APL",
	43:    "DS",
	44:    "SSHFP",
	46:    "RRSIG",
	47:    "NSEC",
	48:    "DNSKEY",
	49:    "DHCID",
	50:    "NSEC3",
	51:    "NSEC3PARAM",
	52:    "TLSA",
	53:    "SMIMEA",
	55:    "HIP",
	56:    "NINFO",
	57:    "RKEY",
	58:    "TALINK",
	59:    "CDS",
	60:    "CDNSKEY",
	61:    "OPENPGPKEY",
	62:    "CSYNC",
	63:    "ZONEMD",
	64:    "SVCB",
	65:    "HTTPS",
	99:    "SPF",
	100:   "UINFO",
	101:   "UID",
	102:   "GID",
	103:   "UNSPEC",
	104:   "NID",
	105:   "L32",
	106:   "L64",
	107:   "LP",
	108:   "EUI48",
	109:   "EUI64",
	256:   "URI",
	257:   "CAA",
	258:   "AVC",
	249:   "TKEY",
	250:   "TSIG",
	251:   "IXFR",
	252:   "AXFR",
	253:   "MAILB",
	254:   "MAILA",
	255:   "ANY",
	32768: "TA",
	32769: "DLV",
	65535: "Reserved",
}

var IntToClass map[uint16]string = map[uint16]string{
	1: "IN",
	2: "CS",
	3: "CH",
	4: "HS",
}

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
		qType, ok := TypeToInt[rr]
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
