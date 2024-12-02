package options

import (
	"flag"
	"strings"
)

type Options struct {
	Host          string
	Path          string
	BypassMethods []string
	Headers       map[string]string
	Timeout       float64
	Threads       int
}

type headerFlag map[string]string

func (h headerFlag) String() string {
	var headers []string
	for key, value := range h {
		headers = append(headers, key+":"+value)
	}
	return strings.Join(headers, ", ")
}

func (h headerFlag) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return flag.ErrHelp
	}
	h[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	return nil
}

func ParseFlags() Options {
	host := flag.String("h", "", "Target host")
	path := flag.String("p", "/", "Full request path (starting with /)")
	bypassMethods := flag.String("b", "header,path,method,encode", "Comma-separated bypass methods")
	timeout := flag.Float64("t", 10, "Request timeout in seconds")
	threads := flag.Int("n", 20, "Number of threads")

	headers := headerFlag{}
	flag.Var(headers, "H", "Custom header in key:value format")

	flag.Parse()

	selectedMethods := strings.Split(*bypassMethods, ",")

	return Options{
		Host:          *host,
		Path:          *path,
		BypassMethods: selectedMethods,
		Headers:       headers,
		Timeout:       *timeout,
		Threads:       *threads,
	}
}
