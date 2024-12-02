package method

import (
	"fmt"
	"log"
	"sync"

	"4no3/pkg/httpclient"
	"4no3/pkg/util"
)

func applyMethod(httpClient *httpclient.HttpClient, method string, sem chan struct{}) {
	defer func() { <-sem }()

	baseOptions := httpClient.GetOptions()
	customOptions := &httpclient.ClientOptions{
		Host:             baseOptions.Host,
		Path:             baseOptions.Path,
		Headers:          baseOptions.Headers,
		Method:           method,
		Body:             baseOptions.Body,
		Timeout:          baseOptions.Timeout,
		ForceHttpVersion: baseOptions.ForceHttpVersion,
		ThreadLimit:      baseOptions.ThreadLimit,
	}

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for method %s: %v", method, err)
		return
	}
	defer response.Body.Close()

	util.LogResponseDetails(fmt.Sprintf("Method: %s", method), response)
}

func testProtocolDowngrade(httpClient *httpclient.HttpClient, version string) {
	httpClient.UpdateHTTPVersion(version)
	response, err := httpClient.Do()
	if err != nil {
		log.Printf("Protocol downgrade to %s failed: %v", version, err)
		return
	}
	defer response.Body.Close()
}

func Fuzz(httpClient *httpclient.HttpClient) {
	var wg sync.WaitGroup

	originalOptions := httpClient.GetOptions()

	threadLimit := originalOptions.ThreadLimit
	if threadLimit <= 0 {
		threadLimit = 1
	}
	sem := make(chan struct{}, threadLimit)

	for _, method := range httpMethods {
		wg.Add(1)
		sem <- struct{}{}
		go func(method string) {
			defer wg.Done()
			applyMethod(httpClient, method, sem)
		}(method)
	}

	wg.Wait()

	versions := []string{"0.9", "1.0"}
	for _, version := range versions {
		testProtocolDowngrade(httpClient, version)
	}

	httpClient.UpdateHTTPVersion("1.1")
}
