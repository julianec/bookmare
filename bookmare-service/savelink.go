package main

import (
	"ancient-solutions.com/ancientauth"
	"code.google.com/p/goprotobuf/proto"
	"github.com/julianec/bookmare"
	"net/http"
	"net/url"
        "log"
)

type SaveLink struct {
	db   *BookmarkDB
	auth *ancientauth.Authenticator
}

func (s *SaveLink) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var user string = s.auth.GetAuthenticatedUser(req)
	var bookmark *bookmare.Bookmark = new(bookmare.Bookmark)
        var err error
        var link *url.URL

	if user == "" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Not authenticated."))
		return
	}

	bookmark.Url = proto.String(req.FormValue("url"))
	bookmark.Owner = proto.String(user)
	bookmark.Title = proto.String(req.FormValue("title"))
	bookmark.Description = proto.String(req.FormValue("description"))

        if *bookmark.Url == "" || *bookmark.Title == "" {
                rw.WriteHeader(http.StatusBadRequest)
                rw.Write([]byte("Empty URL or title."))
                return
        }

        link, err = url.Parse(*bookmark.Url)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                rw.Write([]byte(err.Error()))
                return
        }
        if !link.IsAbs() || link.Host == "" {
                rw.WriteHeader(http.StatusBadRequest)
                rw.Write([]byte("Bookmark URL must be absolute."))
                return
        }

        err = s.db.SaveBookmark(bookmark)
        if err != nil {
                rw.WriteHeader(http.StatusInternalServerError)
                rw.Write([]byte(err.Error()))
                log.Print("Error saving to database: ", err)
                return
        }
        rw.WriteHeader(http.StatusOK)
}
