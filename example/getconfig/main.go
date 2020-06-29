package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	xr "github.com/ljm625/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)

}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// YANG path arguments; defaults to "yangocpaths.json"
	ypath := flag.String("ypath", "../input/yangocpaths.json", "YANG path arguments")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	var output string

	// Manually specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername("cisco"),
		xr.WithPassword("cisco"),
		xr.WithHost("[2001:420:2cff:1204::5502:2]:57344"),
		xr.WithTimeout(5),
	)
	if err != nil {
		log.Fatalf("could not build a router, %v", err)
	}

	conn, ctx, err := xr.Connect(*router)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router.Host, err)
	}
	defer conn.Close()

	// Get config for the YANG paths specified on 'js'
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("could not read file: %v: %v\n", *ypath, err)
	}
	output, err = xr.GetConfig(ctx, conn, string(js), id)
	if err != nil {
		log.Fatalf("could not get the config from %s, %v", router.Host, err)
	}
	fmt.Printf("\nconfig from %s\n %s\n", router.Host, output)
}
