package api

import (
    "log"
    "net/http"
)

type Todo struct {
    Name string
}

var _ http.Handler = &Todo{""}

func (todo *Todo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    log.Printf("TODO: %s -> %s", todo.Name, req.URL.EscapedPath())
}
