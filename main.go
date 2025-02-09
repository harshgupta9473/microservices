package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/microservices/handlers"
)

func main() {
	l:=log.New(os.Stdout,"product-api",log.LstdFlags)

	// hh:=handlers.NewHello(l)
	// gh:=handlers.NewGoodbye(l)
	ph:=handlers.NewProducts(l)

	// sm:=http.NewServeMux()
	sm:=mux.NewRouter()
	getRouter:=sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)

	putRouter:=sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}",ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter:=sm.Methods(http.MethodPost).Subrouter()

	postRouter.HandleFunc("/",ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)
	// sm.Handle("/products",ph)
	// sm.Handle("/goodbye",gh)

	s:=&http.Server{
		Addr: ":9090",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func ()  {
		err:=s.ListenAndServe()
		if err!=nil{
			l.Fatal(err)
		}
	}()



	sigChan:=make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt)
	signal.Notify(sigChan,os.Kill)

	sig:=<-sigChan
	l.Println("Recieved terminate, graceful shutdown",sig)


	tc,_:=context.WithTimeout(context.Background(),30*time.Second)

	s.Shutdown(tc)

	// http.ListenAndServe(":9090",sm)
}