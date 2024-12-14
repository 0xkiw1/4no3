# 4no3 - Golang based 403 & 401 Bypass Testing Tool

4no3 is a golang based tool to help identify common 403 & 401 bypasses. Currently supports 5 bypass methods: headers, connection header, paths, methods, and encodings to analyze server responses and identify potential 403 bypasses.

## Features
- **Header Bypass Testing**: Sends headers with various ip addresses.
- **Connection Header Bypass Testing**: Header bypass with headers being set to hop-by-hop via connection header for reverse proxy bypass.
- **Path Bypass Testing**: Modifies URL paths to test different variations.
- **Method Bypass Testing**: Uses various HTTP methods for request manipulation.
- **Encoding Bypass Testing**: Encodes parts of the path to test server behavior.

## Usage
```
Usage of 4no3:
  -H value
        Custom header in key:value format
  -b string
        Comma-separated bypass methods (default "header,connection,path,method,encode")
  -h string
        Target host
  -n int
        Number of threads (default 20)
  -p string
        Full request path (starting with /) (default "/")
  -pw string
        Path to the wordlist for path fuzzing
  -t float
        Request timeout in seconds (default 10)
```

## Example usage
`./4no3 -h https://example.com -p /admin -b header,method,connection -H "host:admin.example.com"`

`./4no3 -h https://example.com -p /admin/console -n 5 -t 20 -H "Authorization:test" -H "Header:value" -pw ~/wordlists/paths.txt`

## Path wordlist format
```
$1 - full path except for the last dir (/api/test/admin) -> $1 = /api/test
$2 - last dir only (/api/test/admin) -> $2 = admin

Please refer to wordlists/paths.txt for examples
```

## Version 1.3
TODO:
- add wordlist support for headers and their values
