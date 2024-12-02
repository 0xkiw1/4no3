package encode

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"4no3/pkg/httpclient"
	"4no3/pkg/util"
)

func uppercaseLastDir(inputPath string) string {
	return util.ApplyToLastDir(inputPath, strings.ToUpper)
}

func uRLEncodeLastDirN(inputPath string, n int) string {
	path := inputPath
	for i := 0; i < n; i++ {
		path = util.ApplyToLastDir(path, util.URLEncodeString)
	}
	return path
}

func unicodeTranslateLastDir(inputPath string) string {
	return util.ApplyToLastDir(inputPath, util.ReplaceUnicode)
}

func applyPathModification(httpClient *httpclient.HttpClient, modifyFunc func(string) string, sem chan struct{}) {
	defer func() { <-sem }()

	baseOptions := httpClient.GetOptions()
	modifiedPath := modifyFunc(baseOptions.Path)

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

	util.LogResponseDetails(fmt.Sprintf("Encoded path: %s", modifiedPath), response)
}

func Fuzz(httpClient *httpclient.HttpClient) {
	var wg sync.WaitGroup

	originalOptions := httpClient.GetOptions()
	threadLimit := originalOptions.ThreadLimit
	if threadLimit <= 0 {
		threadLimit = 1
	}
	sem := make(chan struct{}, threadLimit)

	modifications := []func(string) string{
		uppercaseLastDir,
		unicodeTranslateLastDir,
	}

	for i := 1; i <= 3; i++ {
		n := i
		modifications = append(modifications, func(path string) string {
			return uRLEncodeLastDirN(path, n)
		})
	}

	for _, modifyFunc := range modifications {
		wg.Add(1)
		sem <- struct{}{}
		go func(modifyFunc func(string) string) {
			defer wg.Done()

			applyPathModification(httpClient, modifyFunc, sem)
		}(modifyFunc)
	}

	wg.Wait()

	httpClient.SetOptions(&originalOptions)
}
