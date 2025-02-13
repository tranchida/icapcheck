package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	ic "github.com/egirna/icap-client"
)

func main() {

	server := flag.String("server", "", "ICAP server address (mandatory, format host:port)")
	file := flag.String("file", "", "File to be scanned (mandatory)")
	debug := flag.Bool("debug", false, "Enable debug mode (optional, default false)")

	flag.Parse()

	if *server == "" || *file == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	scanFile, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	/* preparing the http request required for the REQMOD */
	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8000/dummy", bytes.NewBuffer(scanFile))
	if err != nil {
		log.Fatal(err)
	}
	httpReq.Header.Set("Content-Type", "application/octet-stream")

	/* making a icap request with REQMOD method */
	req, err := ic.NewRequest(ic.MethodREQMOD, fmt.Sprintf("icap://%s/reqmod", *server), httpReq, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-File-Name", *file)

	if *debug {
		dump, err := ic.DumpRequest(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dump))
	}

	/* making the icap client responsible for making the requests */
	client := &ic.Client{
		Timeout: 10 * time.Second,
	}

	/* making the REQMOD request call */
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		fmt.Println(resp.Status)
		fmt.Println(resp.StatusCode)
	}

	switch resp.StatusCode {

	case http.StatusNoContent:
		fmt.Println("No virus found")
	case http.StatusOK:
		fmt.Println("Virus found")
		if *debug {
			body, _ := io.ReadAll(resp.ContentResponse.Body)
			fmt.Println(string(body))
		}
	default:
		fmt.Printf("Unexpected response %d %s\n", resp.StatusCode, resp.Status)

	}
}
