package httpclient

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"4no3/pkg/options"
)

type ClientOptions struct {
	Host             string
	Headers          map[string]string
	Method           string
	Body             io.Reader
	Timeout          time.Duration
	ForceHttpVersion string
	Path             string
	ThreadLimit      int
}

type HttpClient struct {
	client  *http.Client
	Options *ClientOptions
	mutex   sync.RWMutex
}

func (c *HttpClient) GetOptions() ClientOptions {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	optionsCopy := ClientOptions{
		Host:             c.Options.Host,
		Path:             c.Options.Path,
		Headers:          copyHeaders(c.Options.Headers),
		Method:           c.Options.Method,
		Body:             copyBody(c.Options.Body),
		Timeout:          c.Options.Timeout,
		ForceHttpVersion: c.Options.ForceHttpVersion,
		ThreadLimit:      c.Options.ThreadLimit,
	}

	return optionsCopy
}

func copyHeaders(headers map[string]string) map[string]string {
	copy := make(map[string]string, len(headers))
	for k, v := range headers {
		copy[k] = v
	}
	return copy
}

func copyBody(body io.Reader) io.Reader {
	if body == nil {
		return nil
	}

	buf := new(strings.Builder)
	_, err := io.Copy(buf, body)
	if err != nil {
		return nil
	}

	return strings.NewReader(buf.String())
}

func (c *HttpClient) SetOptions(options *ClientOptions) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Options = options
}

func getClientOptions(options *options.Options) *ClientOptions {
	clientOptions := ClientOptions{
		Host:             options.Host,
		Path:             options.Path,
		Headers:          options.Headers,
		Method:           "GET",
		Body:             nil,
		Timeout:          time.Duration(float64(options.Timeout) * float64(time.Second)),
		ForceHttpVersion: "",
		ThreadLimit:      options.Threads,
	}

	return &clientOptions
}

func NewHttpClient(customOptions *options.Options) *HttpClient {
	options := getClientOptions(customOptions)

	httpClient := &HttpClient{
		client: &http.Client{
			Timeout: options.Timeout,
		},
		Options: options,
	}

	httpClient.UpdateHTTPVersion("1.1")

	return httpClient
}

func (c *HttpClient) UpdateHTTPVersion(version string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Options.ForceHttpVersion = version

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if version == "1.0" {
		transport.ForceAttemptHTTP2 = false
		transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	} else if version == "1.1" {
		transport.ForceAttemptHTTP2 = false
	} else if version == "2" {
		transport.ForceAttemptHTTP2 = true
	}

	c.client.Transport = transport
}

func (c *HttpClient) replaceHostHeader(rawRequest string) (string, url.URL, error) {
	c.mutex.RLock()
	options := c.Options
	c.mutex.RUnlock()

	parsedURL, err := url.Parse(options.Host)
	if err != nil {
		return "", url.URL{}, err
	}

	hostHeader := ""
	key := "Host"
	value, exists := options.Headers[key]
	if exists {
		hostHeader = value
	} else {
		hostHeader = parsedURL.Hostname()
	}

	rawRequest = strings.ReplaceAll(rawRequest, fmt.Sprintf("Host: %s", options.Host), fmt.Sprintf("Host: %s", hostHeader))

	return rawRequest, *parsedURL, err
}

func (c *HttpClient) sendRawRequest(rawRequest string) (*http.Response, error) {
	rawRequest, parsedURL, err := c.replaceHostHeader(rawRequest)
	if err != nil {
		return nil, err
	}

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		if parsedURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	//fmt.Println(rawRequest)
	var conn net.Conn
	dialer := &net.Dialer{
		Timeout: c.Options.Timeout,
	}

	if parsedURL.Scheme == "https" {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		conn, err = tls.DialWithDialer(dialer, "tcp", host, tlsConfig)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, fmt.Errorf("connection timeout to %s", host)
			}
			return nil, fmt.Errorf("TLS connection error to %s", host)
		}
	} else {
		conn, err = dialer.Dial("tcp", host)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, fmt.Errorf("connection timeout to %s", host)
			}
			return nil, fmt.Errorf("TCP connection error to %s", host)
		}
	}
	defer conn.Close()

	_, err = conn.Write([]byte(rawRequest))
	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("response timeout from %s", host)
		}
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			resp.Body.Close()
			resp.Body = io.NopCloser(strings.NewReader(string(body)))
		}
	}

	return resp, nil
}

func (c *HttpClient) Do(customOptions ...*ClientOptions) (*http.Response, error) {
	options := c.GetOptions()
	if len(customOptions) > 0 && customOptions[0] != nil {
		options = *customOptions[0]
	}

	rawRequest := c.buildRawRequestBody(options)

	return c.sendRawRequest(rawRequest)
}

func (c *HttpClient) buildRawRequestBody(options ClientOptions) string {
	var rawRequest strings.Builder
	rawRequest.WriteString(fmt.Sprintf("%s %s HTTP/%s\r\n", options.Method, options.Path, options.ForceHttpVersion))

	for key, value := range options.Headers {
		if strings.ToLower(key) == "host" {
			continue
		}
		rawRequest.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	rawRequest.WriteString(fmt.Sprintf("Host: %s\r\n", options.Host))

	value, exists := options.Headers["Connection"]
	if exists {
		rawRequest.WriteString(fmt.Sprintf("Connection: %s\r\n", value))
	} else {
		rawRequest.WriteString("Connection: close\r\n")
	}
	rawRequest.WriteString("\r\n")

	if options.Body != nil {
		bodyBytes, _ := io.ReadAll(options.Body)
		rawRequest.WriteString(string(bodyBytes))
	}

	return rawRequest.String()
}

func PrintR(response *http.Response) {
	fmt.Printf("Response Status: %s\n", response.Status)

	fmt.Println("Response Headers:")
	for key, values := range response.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Println("Response Body:")
	fmt.Println(string(body))
}
