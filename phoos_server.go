package main

import (
  "fmt"
  "net/http"
)

func headers(w http.ResponseWriter, req *http.Request) {
  for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func main() {
    fs := http.FileServer(http.Dir("webpage/"))
	  http.Handle("/", fs)
    http.HandleFunc("/headers", headers)
    http.ListenAndServer(":3032", nil)
}
