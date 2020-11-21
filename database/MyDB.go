package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hackerclever/gofinal/customer"
	_ "github.com/lib/pq"
)

// MyDB like a instance of DB
type MyDB struct{}

var db *sql.DB

var insertCustomerStmt = "INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id"
var selectCustomerStmt = "SELECT id, name, email, status FROM customers WHERE id=$1"
var selectCustomerAllstmt = "SELECT id, name, email, status FROM customers"
var updateCustomerStmt = "UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1"
var deleteStmt = "DELETE FROM customers WHERE id=$1"

// Conn Get DB Connection
func (myDB MyDB) Conn(url string) *sql.DB {
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	fmt.Println("Connect to database success.")
	return db
}

// CreateCustomersTb Create Table custormers when it's not exist.
func (myDB MyDB) CreateCustomersTb() error {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("Can't create table", err)
	}
	return nil
}

// CreateCustomer Insert new customer to DB
func (myDB MyDB) CreateCustomer(c customer.Customer) (customer.Customer, error) {
	row := db.QueryRow(insertCustomerStmt, c.Name, c.Email, c.Status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return customer.Customer{}, err
	}
	c.ID = id
	return c, nil
}

// GetCustomer Get customer by id
func (myDB MyDB) GetCustomer(id int) (*customer.Customer, error) {
	c := customer.Customer{}
	stmt, err := db.Prepare(selectCustomerStmt)
	if err != nil {
		return &customer.Customer{}, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
	if err != nil {
		fmt.Println("Test in")
		fmt.Println(err)
		return &customer.Customer{}, err
	}

	return &c, nil
}

// GetAllCustomer Get all customer
func (myDB MyDB) GetAllCustomer() (*[]customer.Customer, error) {
	stmt, err := db.Prepare(selectCustomerAllstmt)
	if err != nil {
		return &[]customer.Customer{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return &[]customer.Customer{}, err
	}

	cuss := []customer.Customer{}
	for rows.Next() {
		cus := customer.Customer{}
		err := rows.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
		if err != nil {
			return &[]customer.Customer{}, err
		}

		cuss = append(cuss, cus)
	}

	return &cuss, nil
}

// UpdateCustomer Update customer detail
func (myDB MyDB) UpdateCustomer(c customer.Customer) (*customer.Customer, error) {
	stmt, err := db.Prepare(updateCustomerStmt)
	if err != nil {
		return &customer.Customer{}, err
	}

	_, err = stmt.Exec(c.ID, c.Name, c.Email, c.Status)
	if err != nil {
		return &customer.Customer{}, err
	}

	return &c, nil
}

// DeleteCustomer Delete custermer by id
func (myDB MyDB) DeleteCustomer(id int) error {
	stmt, err := db.Prepare(deleteStmt)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
