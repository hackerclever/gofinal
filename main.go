package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hackerclever/gofinal/database"
	"github.com/hackerclever/gofinal/handlers"
)

func main() {
	fmt.Println("customer service")
	var myDB = database.MyDB{}
	myDB.Conn(os.Getenv("DATABASE_URL"))
	myDB.CreateCustomersTb()

	r := gin.Default()
	r.POST("/customers", handlers.CreateCustomer)
	r.GET("/customers/:id", handlers.GetCustomer)
	r.GET("/customers", handlers.GetAllCustomer)
	r.PUT("/customers/:id", handlers.UpdateCustomer)
	r.DELETE("/customers/:id", handlers.DeleteCustomer)
	err := r.Run(":2009")
	if err != nil {
		log.Fatal("Can't run server!", err)
	}
	fmt.Println("Customer service running.")
}
