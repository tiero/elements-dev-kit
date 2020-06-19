# HTTP utility

This package defines internal HTTP utilities.

## Methods

```go
NewHTTPRequest(method string, url string, bodyString string, header map[string]string) (string, error)
```

## Example

```go
method := "GET"
url := "https://api.blockcypher.com/v1/btc/test3/addrs/mtXWDB6k5yC5v7TcwKZHB89SUp85yCKshy?unspentOnly=true"
response, err := uhttp.NewHttpRequest(method, url, "", nil)
if err != nil {
  fmt.Println("An error occurred:", err)
  return
}

fmt.Println("Response:", response)
```