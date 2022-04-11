package main

import (
	"context"
	// "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"app/handlers"
	mydb "app/mydb"
)

func main() {

	l := log.New(os.Stdout, "app.com", log.LstdFlags)
	db, err := mydb.GetDb()
	if err != nil {
		l.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/createuser", handlers.SaveUser)
	http.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
		// hash, err := u.HashPassword("secret12345")
		// if err != nil {
		// 	panic(err)
		// }

		// err = u.ComparePassword("secret12345", []byte(hash))
		// if err != nil {
		// 	rw.Write([]byte("not logged in!"))
		// 	log.Default().Println(err.Error())
		// }
		rw.Write([]byte("logged in!"))

	})
	s := &http.Server{
		Addr: ":9090",
		// Handler:      mycors(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	gracefulShutdown(s, l)

}

func gracefulShutdown(s *http.Server, l *log.Logger) {
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Received terminate shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
