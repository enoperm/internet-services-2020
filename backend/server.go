package main

import (
    "context"
    "database/sql"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

    api "github.com/enoperm/internet-services-2020/api"
    appdb "github.com/enoperm/internet-services-2020/db"
)

func main() {
	db := configureDb()
    urlRouter := mux.NewRouter()
    configureRouter(urlRouter, db)

    server := http.Server{
        Addr: "0.0.0.0:2000",
        Handler: urlRouter,

        WriteTimeout: 10 * time.Second,
        ReadTimeout:  10 * time.Second,
        IdleTimeout:  10 * time.Second,
    }

    go func() {
        err := server.ListenAndServe()
        if err != nil {
            log.Fatalln(err)
        }
    }()

    awaitShutdown(&server)
}

func configureDb() *appdb.ApplicationDatabase {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalln("configure-db:", err)
	}
	return &appdb.ApplicationDatabase{db}
}

func configureRouter(r *mux.Router, db *appdb.ApplicationDatabase) {
	{
		sr := r.PathPrefix("/register").Subrouter()
		api.NewRegisterApi(sr, db)
	}
    r.PathPrefix("/session").Handler(&api.Todo{Name: "session"})
    r.PathPrefix("/profile").Handler(&api.Todo{Name: "profile"})
    r.PathPrefix("/savings").Handler(&api.Todo{Name: "savings"})
    r.PathPrefix("/rankings").Handler(&api.Todo{Name: "rankings"})
    r.PathPrefix("/motd").Handler(&api.Todo{Name: "motd"})

    r.PathPrefix("/").HandlerFunc(logUnhandledPrefix)

    r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		log.Println(route.GetPathTemplate())
		return nil
	})
}

func logUnhandledPrefix(rw http.ResponseWriter, req *http.Request) {
    log.Printf("request to unhandled prefix '%s' from '%s'", req.URL.String(), req.RemoteAddr)
}

func awaitShutdown(server *http.Server) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    <-sigChan

    ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
    defer cancel()

    server.Shutdown(ctx)
}
