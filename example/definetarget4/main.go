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
	// YANG path arguments; defaults to "yangpaths.json"
	ypath := flag.String("ypath", "../input/yangpaths.json", "YANG path arguments")
	flag.Parse()

	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	var output string

	// Manually specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.1:57344"),
		xr.WithCert("ems.pem"),
		xr.WithTimeout(5),
	)
	if err != nil {
		log.Fatalf("target parameters are incorrect: %s", err)
	}

	// Setup a connection to the target.
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
