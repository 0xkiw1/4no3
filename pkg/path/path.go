package path

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"4no3/pkg/httpclient"
	"4no3/pkg/util"
)

func applyPathModification(httpClient *httpclient.HttpClient, pathTemplate string, dir1 string, dir2 string) {
	baseOptions := httpClient.GetOptions()

	modifiedPath := strings.ReplaceAll(pathTemplate, "$1", dir1)
	modifiedPath = strings.ReplaceAll(modifiedPath, "$2", dir2)

	customOptions := &httpclient.ClientOptions{
		Host:             baseOptions.Host,
		Path:             modifiedPath,
		Headers:          baseOptions.Headers,
		Method:           baseOptions.Method,
		Body:             baseOptions.Body,
		Timeout:          baseOptions.Timeout,
		ForceHttpVersion: baseOptions.ForceHttpVersion,
		ThreadLimit:      baseOptions.ThreadLimit,
	}

	response, err := httpClient.Do(customOptions)
	if err != nil {
		log.Printf("Request failed for path %s: %v", customOptions.Path, err)
		return
	}
	defer response.Body.Close()

	util.LogResponseDetails(fmt.Sprintf("Path: %s", modifiedPath), response)
}

func Fuzz(httpClient *httpclient.HttpClient, wordlist []string) {
	var wg sync.WaitGroup

	originalOptions := httpClient.GetOptions()

	threadLimit := originalOptions.ThreadLimit
	if threadLimit <= 0 {
		threadLimit = 1
	}
	sem := make(chan struct{}, threadLimit)

	dir1, dir2 := util.SplitDir(originalOptions.Path)

	for _, pathTemplate := range wordlist {
		wg.Add(1)
		sem <- struct{}{}
		go func(pathTemplate string) {
			defer func() { <-sem }()
			defer wg.Done()

			applyPathModification(httpClient, pathTemplate, dir1, dir2)
		}(pathTemplate)
	}

	wg.Wait()
}
