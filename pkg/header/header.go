package header

import (
	"fmt"
	"log"
	"sync"

	"4no3/pkg/httpclient"
	"4no3/pkg/util"
)

func applyIPHeader(httpClient *httpclient.HttpClient, header string, ip string, sem chan struct{}) {
	defer func() { <-sem }()

	baseOptions := httpClient.GetOptions()
	customOptions := &httpclient.ClientOptions{
		Host:             baseOptions.Host,
		Path:             baseOptions.Path,
		Headers:          copyHeaders(baseOptions.Headers),
		Method:           baseOptions.Method,
		Body:             baseOptions.Body,
		Timeout:          baseOptions.Timeout,
		ForceHttpVersion: baseOptions.ForceHttpVersion,
		ThreadLimit:      baseOptions.ThreadLimit,
	}
	customOptions.Headers[header] = ip

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for header %s with IP %s: %v", header, ip, err)
	} else {
		defer response.Body.Close()
	}

	util.LogResponseDetails(fmt.Sprintf("Header: %s: %s", header, ip), response)
}

func applyPathHeader(httpClient *httpclient.HttpClient, header string, sem chan struct{}) {
	defer func() { <-sem }()

	baseOptions := httpClient.GetOptions()
	customOptions := &httpclient.ClientOptions{
		Host:             baseOptions.Host,
		Path:             "/",
		Headers:          copyHeaders(baseOptions.Headers),
		Method:           baseOptions.Method,
		Body:             baseOptions.Body,
		Timeout:          baseOptions.Timeout,
		ForceHttpVersion: baseOptions.ForceHttpVersion,
		ThreadLimit:      baseOptions.ThreadLimit,
	}
	customOptions.Headers[header] = baseOptions.Path

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for header %s: %v", header, err)
		return
	}
	defer response.Body.Close()

	util.LogResponseDetails(fmt.Sprintf("Header: %s: %s, Path: %s,", header, baseOptions.Path, "/"), response)
}

func copyHeaders(headers map[string]string) map[string]string {
	newHeaders := make(map[string]string)
	for k, v := range headers {
		newHeaders[k] = v
	}
	return newHeaders
}

func Fuzz(httpClient *httpclient.HttpClient) {
	var wg sync.WaitGroup

	originalOptions := httpClient.GetOptions()
	threadLimit := originalOptions.ThreadLimit
	if threadLimit <= 0 {
		threadLimit = 1
	}
	sem := make(chan struct{}, threadLimit)

	for _, header := range ipHeaders {
		for _, ip := range ipValues {
			wg.Add(1)
			sem <- struct{}{}
			go func(header, ip string) {
				defer wg.Done()
				applyIPHeader(httpClient, header, ip, sem)
			}(header, ip)
		}
	}

	for _, header := range pathHeaders {
		wg.Add(1)
		sem <- struct{}{}
		go func(header string) {
			defer wg.Done()
			applyPathHeader(httpClient, header, sem)
		}(header)
	}

	wg.Wait()
}
