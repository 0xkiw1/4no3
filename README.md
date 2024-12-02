# 4no3 - Golang based 403 & 401 Bypass Testing Tool

4no3 is a golang based tool to help identify common 403 & 401 bypasses. Currently supports 4 bypass methods: headers, paths, methods, and encodings to analyze server responses and identify potential 403 bypasses.

## Features
- **Header Bypass Testing**: Sends headers with various ip addresses.
- **Path Bypass Testing**: Modifies URL paths to test different variations.
- **Method Bypass Testing**: Uses various HTTP methods for request manipulation.
- **Encoding Bypass Testing**: Encodes parts of the path to test server behavior.

## Usage
```
Usage of 4no3:
  -H value
        Custom header in key:value format
  -b string
        Comma-separated bypass methods (default "header,path,method,encode")
  -h string
        Target host
  -n int
        Number of threads (default 20)
  -p string
        Full request path (starting with /) (default "/")
  -t float
        Request timeout in seconds (default 10)
```

## Example usage
`./4no3 -h https://example.com -p /admin -b header,method -H "host:admin.example.com"`

`./4no3 -h https://example.com -p /admin/console -n 5 -t 20 -H "Authorization:test" -H "Header:value"`

## Version 1.0
TODO:
- add wordlist support for headers and their values
- add wordlist support for custom path fuzzing
