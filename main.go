package main

import (
	"fmt"
    "log"
	"net/http"
    "context"
    "os"
    "os/signal"
    "time"
    "./database"
    "./handler"
    "./config"
)


func main() {
    var err error

    //Open connection to the database
	database.DB, err = database.GetDB(config.Settings.User, config.Settings.Password, config.Settings.Database)

    if err != nil {
		log.Fatal(err)
	}
    defer database.DB.Close()

    //Create the server
    server := &http.Server {
        Addr: fmt.Sprintf("%s:%s", config.Settings.Host, config.Settings.Port),
        Handler: handler.Handlers(),
    }

    //Run the server / prepare for possible shutdown
    stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

    go func() {
        fmt.Println("Listening on ", server.Addr)
        if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
    }()

    <-stop

    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

    server.Shutdown(ctx)
    fmt.Println("Graceful shutdown.")
}
