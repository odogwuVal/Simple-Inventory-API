package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPass, "test")
	if err != nil {
		log.Fatal("Error occured while initializing the database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10, 7),
		PRIMARY KEY(id)
	);`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	log.Println("ClearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500.00)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateTable(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chair", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	// if m["name"] != "chair" {
	// 	t.Errorf("Expected %v, Got: %v", "chair", m["name"])
	// }
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 300)

	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}
