package main

import (
	"fmt"
	"os"
	"sync"

	"git.sr.ht/~tephra/rog"
	"github.com/miekg/dns"
	"github.com/pborman/getopt/v2"
)

const (
	MAJOR = 0
	MINOR = 1
	PATCH = 0

	EXIT_OK     = 0
	EXIT_ERR    = 1
	EXIT_USAGE  = 64
	EXIT_NOHOST = 68
)

var answerMap = map[string][]*dns.Msg{}
var mapLock sync.Mutex

func version() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR, MINOR, PATCH)
}

func main() {
	versionFlag := getopt.BoolLong("version", 'v', "Print out version information and exit")
	verboseFlag := getopt.BoolLong("verbose", 'V', "Set the verbose flag which makes other commands print out more verbose output")

	udp := true
	tcp := false

	getopt.FlagLong(&udp, "udp", '\x00', "Set rog to use UDP for queries, default: true. Mutually exclusive with the tcp option").SetGroup("protocol")
	getopt.FlagLong(&tcp, "tcp", '\x00', "Set rog to use TCP for queries, default: false. Mutually exclusive with the udp option").SetGroup("protocol")

	rrTypesOptions := getopt.ListLong("types", 't', "Comma separated RR types that are to be queried. Defaults to A records. The special string \"all\" queries all RRs.")
	getopt.Parse()

	rrTypes := []string{}
	if len(*rrTypesOptions) == 0 {
		rrTypes = append(rrTypes, "A")
	} else {
		rrTypes = *rrTypesOptions
	}

	if *versionFlag {
		fmt.Printf("rog version %s\n", version())

		if *verboseFlag {
			// todo: What is verbose version? License, build info, git tag?
		}

		os.Exit(EXIT_OK)
	}

	wg := sync.WaitGroup{}

	for _, host := range getopt.Args() {
		wg.Add(1)

		go func(h string) {
			defer wg.Done()
			answer, err := rog.Query(h, rog.QueryOption{
				TCP: tcp,
				RRs: rrTypes,
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			mapLock.Lock()
			answerMap[h] = answer
			mapLock.Unlock()

			return
		}(host)
	}

	wg.Wait()
	rog.FormatOutput(answerMap)
	os.Exit(EXIT_OK)
}
