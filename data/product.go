package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32  `json:"price" validate:"gt=0"`
	SKU         string   `json:"sku" validate:"required,sku"`
	CreatedOn   string   `json:"-"`
	UpdatedOn   string   `json:"-"`
	DeletedOn   string   `json:"-"`
}

func (p *Product)FromJSON(r io.Reader)error{
	e:=json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product)Validate()error{
	validate:=validator.New()
	validate.RegisterValidation("sku",validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel)bool{
	// sku is of  formate abc-absd-dfsdf
	reg:=regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	mathces:=reg.FindAllString(fl.Field().String(),-1)

	if len(mathces)!=1{
		return false
	}
	return true
}
type Products []*Product

func (p *Products)ToJSON(w io.Writer)error{
	e:=json.NewEncoder(w)
	return e.Encode(p)
}

func AddProduct(p *Product){
	p.ID=getNextID()
	prductList=append(prductList, p)
}

func UpdateProduct(id int,prod *Product)error{
	_,pos,err:=findProduct(id)
	if err!=nil{
		return err
	}
	prod.ID=id
	prductList[pos]=prod
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int)(*Product,int,error){
	for i,p:=range prductList{
		if p.ID==id{
			return p,i,nil
		}
	}
	return nil,-1, ErrProductNotFound
}

func getNextID()int{
	lp:=prductList[len(prductList)-1]
	return lp.ID+1
}



func GetProducts()Products{
	return prductList
}

var prductList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn: time.Now().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn: time.Now().String(),
	},
}