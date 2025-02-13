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

	server := flag.String("server", "", "ICAP server address")
	file := flag.String("file", "", "File to be scanned")
	debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

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
	req.Header.Set("X-File-Name", *file)

	if *debug {
		dump, err := ic.DumpRequest(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dump))
	}

	if err != nil {
		log.Fatal(err)
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

	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("No virus found")
	} else {
		fmt.Println("Virus found")
		if *debug {
			body, _ := io.ReadAll(resp.ContentResponse.Body)
			fmt.Println(string(body))
		}
	}

}
