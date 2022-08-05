package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cloudflare/roughtime"
)

func main() {
	cfg := flag.String("c", "", "Configuration file with list of servers")
	single := flag.Bool("s", false, "Only test against a single server")
	flag.Parse()

	if *cfg == "" {
		fmt.Println("No config file provided")
		os.Exit(1)
	}
	servers, skipped, err := roughtime.LoadConfig(*cfg)
	if skipped > 0 {
		fmt.Printf("skipped %d servers\n", skipped)
	}
	if err != nil {
		panic(err)
	}

	var prev *roughtime.Roughtime
	var rt *roughtime.Roughtime
	for _, s := range servers {
		fmt.Printf("Contacting %s@%s\n", s.Name, s.Addresses[0].Address)
		rt, err = roughtime.Get(&s, 3, time.Second, prev)
		if err != nil {
			fmt.Println(err)
		}

		if *single {
			break
		}
	}
	if rt == nil {
		fmt.Println("Could not get time")
		os.Exit(1)
	}

	now, delta := rt.Now()
	fmt.Printf("Current time is %s +/- %s\n", now, delta)
}
