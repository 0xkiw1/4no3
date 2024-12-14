package main

import (
	"bufio"
	"log"
	"os"

	"4no3/pkg/encode"
	"4no3/pkg/header"
	"4no3/pkg/httpclient"
	"4no3/pkg/method"
	"4no3/pkg/options"
	"4no3/pkg/path"
	"4no3/pkg/util"
)

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	util.PrintASCIIArt()
	options := options.ParseFlags()

	for _, bypass := range options.BypassMethods {
		if bypass == "path" && options.PathWordlist == "" {
			log.Fatalf("Error: Path fuzzing method selected but no path wordlist (-pw) specified.")
		}
	}

	httpClient := httpclient.NewHttpClient(&options)

	for _, bypass := range options.BypassMethods {
		switch bypass {
		case "path":
			var wordlist []string
			file, err := os.Open(options.PathWordlist)
			if err != nil {
				log.Fatalf("Failed to open path wordlist: %v", err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				wordlist = append(wordlist, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				log.Fatalf("Error reading path wordlist: %v", err)
			}
			util.PrintBypassName("Path bypass")
			path.Fuzz(httpClient, wordlist)
			util.PrintBypassDelimeter()
		case "method":
			util.PrintBypassName("HTTP Method bypass")
			method.Fuzz(httpClient)
			util.PrintBypassDelimeter()
		case "encode":
			util.PrintBypassName("Path encoding bypass")
			encode.Fuzz(httpClient)
			util.PrintBypassDelimeter()
		case "header":
			util.PrintBypassName("Header bypass")
			header.Fuzz(httpClient, false)
			util.PrintBypassDelimeter()
		case "connection":
			util.PrintBypassName("Connection header bypass")
			header.Fuzz(httpClient, true)
			util.PrintBypassDelimeter()
		}
	}
}
