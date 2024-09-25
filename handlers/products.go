package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}


func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	// d,err:=json.Marshal(lp)
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
	// w.Write(d)
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	prod:=r.Context().Value(KeyProduct{}).(data.Product)

	// prod := &data.Product{}
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	// 	return
	// }
	// p.l.Printf("Prod %#v",prod)
	data.AddProduct(&prod)
}


func (p *Products) UpdateProducts( w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")
	vars:=mux.Vars(r)
	idstring:=vars["id"]
	id,err:=strconv.Atoi(idstring)
	if err!=nil{
		http.Error(w,"can't convert idstring to int",http.StatusInternalServerError)
		return
	}

	prod:=r.Context().Value(KeyProduct{}).(data.Product)


	// prod := &data.Product{}
	// err = prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
	// 	return
	// }
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
	}

}

type KeyProduct struct{}

func (p *Products)MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ){
		prod := data.Product{}
	    err := prod.FromJSON(r.Body)
	    if err != nil {
		    http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
		    return
	    }

		err=prod.Validate()
		if err!=nil{
			p.l.Println("[ERROR] validating product",err)
			http.Error(w,fmt.Sprintf("Error validating product: %s",err),http.StatusBadRequest)
			return
		}

		ctx:=context.WithValue(r.Context(),KeyProduct{},prod)
		req:=r.WithContext(ctx)

		next.ServeHTTP(w,req)
	})
} 