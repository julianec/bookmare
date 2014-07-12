package main

import (
        "io"
        "net/http"
        "log"
        "flag"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
        io.WriteString(w, "hello, world!\n")
}

func main() {
        var bind string
        flag.StringVar(&bind, "bind", ":12345", "host:port pair to bind the http server to")
        flag.Parse()
        http.HandleFunc("/", HelloServer)
        err := http.ListenAndServe(bind, nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}

