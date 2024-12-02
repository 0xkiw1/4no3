package main

import (
	"log"

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
	httpClient := httpclient.NewHttpClient(&options)

	for _, bypass := range options.BypassMethods {
		switch bypass {
		case "path":
			path.Fuzz(httpClient)
		case "method":
			method.Fuzz(httpClient)
		case "encode":
			encode.Fuzz(httpClient)
		case "header":
			header.Fuzz(httpClient)
		}
	}
}
