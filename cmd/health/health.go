package main

import (
	"fmt"
	"log"
	"net/http"
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
	fmt.Println(resp.Status)
	defer resp.Body.Close()
}
