package header

import (
	"fmt"
	"log"
	"sync"

	"4no3/pkg/httpclient"
	"4no3/pkg/util"
)

func applyIPHeader(httpClient *httpclient.HttpClient, header string, ip string, sem chan struct{}, connection bool) {
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

	if connection {
		customOptions.Headers["Connection"] = fmt.Sprintf("close, %s", header)
	}

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for header %s (connection header set) with IP %s: %v", header, ip, err)
		return
	} else {
		defer response.Body.Close()
	}

	if connection {
		util.LogResponseDetails(fmt.Sprintf("Header (connection): %s: %s", header, ip), response)
	} else {
		util.LogResponseDetails(fmt.Sprintf("Header: %s: %s", header, ip), response)
	}
}

func applyPathHeader(httpClient *httpclient.HttpClient, header string, sem chan struct{}, connection bool) {
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

	if connection {
		customOptions.Headers["Connection"] = fmt.Sprintf("close, %s", header)
	}

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for header (connection header set) %s: %v", header, err)
		return
	}
	defer response.Body.Close()

	if connection {
		util.LogResponseDetails(fmt.Sprintf("Header (connection): %s: %s, Path: %s,", header, baseOptions.Path, "/"), response)
	} else {
		util.LogResponseDetails(fmt.Sprintf("Header: %s: %s, Path: %s,", header, baseOptions.Path, "/"), response)
	}
}

func copyHeaders(headers map[string]string) map[string]string {
	newHeaders := make(map[string]string)
	for k, v := range headers {
		newHeaders[k] = v
	}
	return newHeaders
}

func Fuzz(httpClient *httpclient.HttpClient, connection bool) {
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
				applyIPHeader(httpClient, header, ip, sem, connection)
			}(header, ip)
		}
	}

	for _, header := range pathHeaders {
		wg.Add(1)
		sem <- struct{}{}
		go func(header string) {
			defer wg.Done()
			applyPathHeader(httpClient, header, sem, connection)
		}(header)
	}

	wg.Wait()
}
