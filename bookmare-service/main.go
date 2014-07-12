package main

import (
        "io"
        "net/http"
        "log"
        "flag"
        "ancient-solutions.com/ancientauth"
)

type HelloServer struct {
    auth *ancientauth.Authenticator
}

// hello world, the web server
func (h *HelloServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
        var user string = h.auth.GetAuthenticatedUser(req)
        if user == "" {
            h.auth.RequestAuthorization(w, req)
            return
        }
        io.WriteString(w, "hello, "+user+"!\n")
}

func main() {
        var app_name, cert_file, key_file, ca_bundle, authserver string
        var bind string
        var err error
        var auth *ancientauth.Authenticator

        flag.StringVar(&bind, "bind", ":12345", "host:port pair to bind the http server to")
        flag.StringVar(&app_name, "app-name", "bookmare", "Application name to present to the authentication server")
        flag.StringVar(&cert_file, "cert", "bookmarks.crt", "Path to the X.509 application certificate")
        flag.StringVar(&key_file, "key", "bookmarks.der", "Path to the X.509 application private key")
        flag.StringVar(&ca_bundle, "ca", "cacert.pem", "Path to the X.509 certificate authority")
        flag.StringVar(&authserver, "auth-server", "login.ancient-solutions.com", "Server for handling authentication")
        flag.Parse()

        auth, err = ancientauth.NewAuthenticator(app_name, cert_file, key_file, ca_bundle, authserver)
        if err != nil {
                log.Fatal("Error setting up authenticator: ", err)
        }
        http.Handle("/", &HelloServer{auth:auth,})
        err = http.ListenAndServe(bind, nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}

