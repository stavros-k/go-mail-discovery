package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	client := http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Get("http://localhost:8080/health")
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	fmt.Println(resp.Status)
	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
