package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/ahmetartuc/dighub/internal/searcher"
)

func main() {
	org := flag.String("org", "", "GitHub organization to scan")
	token := flag.String("token", "", "GitHub Personal Access Token")
	flag.Parse()

	if *org == "" || *token == "" {
		log.Fatal("[!] Usage: dighub -org <target-org> -token <ghp_token>")
	}

	err := searcher.Run(*org, *token)
	if err != nil {
		log.Fatalf("[!] Error during scan: %v", err)
	}

	fmt.Println("[+] Dighub scan complete.")
}
