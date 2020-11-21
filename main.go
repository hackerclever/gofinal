package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func createCustomerHandler(c *gin.Context) {
	cus := Customer{}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stmt := "INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id"
	row := db.QueryRow(stmt, cus.Name, cus.Email, cus.Status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("Can't insert customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cus.ID = id
	c.JSON(http.StatusCreated, cus)
}

func getCustomerHandler(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cus := Customer{}
	stmtSl := "SELECT id, name, email, status FROM customers WHERE id=$1"
	stmt, err := db.Prepare(stmtSl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cus)
}

func getAllCustomerHandler(c *gin.Context) {
	stmtSlAll := "SELECT id, name, email, status FROM customers"
	stmt, err := db.Prepare(stmtSlAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cuss := []Customer{}
	for rows.Next() {
		cus := Customer{}
		err := rows.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		cuss = append(cuss, cus)
	}

	c.JSON(http.StatusOK, cuss)
}

func updateCustomerHandler(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Param("id"))
	cus := Customer{}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cus.ID = id
	stmtUp := "UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1"
	stmt, err := db.Prepare(stmtUp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = stmt.Exec(id, cus.Name, cus.Email, cus.Status)
	// err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cus)
}

func deleteCustomerHandler(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Param("id"))

	stmtDel := "DELETE FROM customers WHERE id=$1"
	stmt, err := db.Prepare(stmtDel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func main() {
	fmt.Println("customer service")
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("Can't create table", err)
	}

	r := gin.Default()
	r.POST("/customers", createCustomerHandler)
	r.GET("/customers/:id", getCustomerHandler)
	r.GET("/customers", getAllCustomerHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerHandler)
	r.Run(":2009")
}
