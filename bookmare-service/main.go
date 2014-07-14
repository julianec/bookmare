package main

import (
        "net/http"
        "log"
        "flag"
        "ancient-solutions.com/ancientauth"
)

type BookmarkSite struct {
    auth *ancientauth.Authenticator
    static_dir string
}

func (h *BookmarkSite) ServeHTTP(w http.ResponseWriter, req *http.Request) {
        var user string = h.auth.GetAuthenticatedUser(req)
        if user == "" {
            h.auth.RequestAuthorization(w, req)
            return
        }
        http.ServeFile(w, req, h.static_dir+"/static/static.html")
}

func main() {
        var app_name, cert_file, key_file, ca_bundle, authserver string
        var dbserver, keyspace string
        var bind string
        var static_dir string
        var err error
        var auth *ancientauth.Authenticator
        var bookmarkdb *BookmarkDB

        flag.StringVar(&bind, "bind", ":12345", "host:port pair to bind the http server to")
        flag.StringVar(&app_name, "app-name", "bookmare", "Application name to present to the authentication server")
        flag.StringVar(&cert_file, "cert", "bookmarks.crt", "Path to the X.509 application certificate")
        flag.StringVar(&key_file, "key", "bookmarks.der", "Path to the X.509 application private key")
        flag.StringVar(&ca_bundle, "ca", "cacert.pem", "Path to the X.509 certificate authority")
        flag.StringVar(&authserver, "auth-server", "login.ancient-solutions.com", "Server for handling authentication")
        flag.StringVar(&static_dir, "static-dir", ".", "Directory to consider for serving static files")
        flag.StringVar(&dbserver, "db-server", "localhost:9160", "Host:Port pair of the cassandra database server")
        flag.StringVar(&keyspace, "db-keyspace", "bookmare", "Name of database keyspace")
        flag.Parse()

        auth, err = ancientauth.NewAuthenticator(app_name, cert_file, key_file, ca_bundle, authserver)
        if err != nil {
                log.Fatal("Error setting up authenticator: ", err)
        }

        bookmarkdb, err = NewBookmarkDB(dbserver, keyspace)
        if err != nil {
                log.Fatal("Error connecting to database: ", err)
        }

        http.Handle("/", &BookmarkSite{auth: auth, static_dir: static_dir,})
        http.Handle("/css/", http.FileServer(http.Dir(static_dir)))
        http.Handle("/js/", http.FileServer(http.Dir(static_dir)))
        http.Handle("/fonts/", http.FileServer(http.Dir(static_dir)))
        http.Handle("/api/savelink", &SaveLink{
                auth: auth,
                db:   bookmarkdb,
        })


        err = http.ListenAndServe(bind, nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}

