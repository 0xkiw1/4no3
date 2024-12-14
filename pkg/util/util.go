package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func PrintASCIIArt() {
	asciiArt := "\033[34m" + `

	_  _              _____ 
	| || |  _ __   ___|___ / 
	| || |_| '_ \ / _ \ |_ \ 
	|__   _| | | | (_) |__) |
	   |_| |_| |_|\___/____/ 
							 
					by 0xkiw1
` + "\033[0m"
	fmt.Println(asciiArt)
}

func PrintBypassName(text string) {
	colorBlue := "\033[34m"
	resetColor := "\033[0m"
	message := colorBlue + text + resetColor
	fmt.Println(message)
	fmt.Print("\n")
}

func PrintBypassDelimeter() {
	fmt.Print("\n")
	fmt.Println("----------------------------------------------------------")
	fmt.Print("\n")
}

func colorizeStatusCode(statusCode int) string {
	color := ""
	reset := "\033[0m"

	switch statusCode / 100 {
	case 1:
		color = "\033[33m"
	case 2:
		color = "\033[32m"
	case 3:
		color = "\033[34m"
	case 4:
		color = "\033[38;5;214m"
	case 5:
		color = "\033[31m"
	default:
		color = reset
	}

	return fmt.Sprintf("%s%d%s", color, statusCode, reset)
}

func LogResponseDetails(info string, response *http.Response) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %v [%s]", err, info)
		return
	}

	status := colorizeStatusCode(response.StatusCode)
	log.Printf("%s [%s, %d]",
		info, status, len(body))
}

func SplitDir(path string) (string, string) {
	path = strings.TrimSuffix(path, "/")

	if !strings.Contains(path, "/") {
		return path, ""
	}

	lastSlashIndex := strings.LastIndex(path, "/")
	dir1 := path[:lastSlashIndex]
	dir2 := path[lastSlashIndex+1:]

	return dir1, dir2
}

func ContainsString(slice *[]string, item string) bool {
	for _, v := range *slice {
		if v == item {
			return true
		}
	}

	return false
}

func URLEncodeString(input string) string {
	result := ""

	for i := 0; i < len(input); i++ {
		result += fmt.Sprintf("%%%02X", input[i])
	}

	return result
}

func ReplaceUnicode(input string) string {
	var builder strings.Builder

	for _, char := range input {
		if translatedChar, exists := TranslationTable[char]; exists {
			builder.WriteRune(translatedChar)
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func ApplyToLastDir(inputPath string, transform func(string) string) string {
	parts := strings.Split(inputPath, "/")

	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			parts[i] = transform(parts[i])
			break
		}
	}

	return strings.Join(parts, "/")
}
