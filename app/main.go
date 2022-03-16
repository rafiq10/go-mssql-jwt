package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"app/mydb"
)

func main() {
	l := log.New(os.Stdout, "app.com", log.LstdFlags)
	db, err := mydb.GetDb()
	if err != nil {
		l.Fatal(err)
	}
	defer db.Close()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				generateKeys()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Print("Authorized")
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

func generateKeys() {
	err := mydb.RunSQL("use [users-db] delete from Keys")
	if err != nil {
		log.Fatalf("Error generating random key: GenerateRandomString(15)=%v", err)
	}
	for i := 0; i < 10; i++ {
		key, err := mydb.GenerateRandomString(15)
		if err != nil {
			log.Fatalf("Error generating random key: GenerateRandomString(15)=%v", err)
		}
		err = mydb.RunSQL("use [users-db] insert into Keys values ('" + key + "') ")
		if err != nil {
			log.Fatalf("Error generating random key: GenerateRandomString(15)=%v", err)
		}
	}

}
