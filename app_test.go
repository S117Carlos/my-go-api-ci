package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(false, DbParams{})
	if err != nil {
		log.Fatalln("Error occurred when initializing app")
	}
	m.Run()
}

func addProduct() {
	p := product{ID: 1, Name: "Toy", Quantity: 1, Price: 2.10}
	MockedProducts = append(MockedProducts, p)
}
func clearProducts() {
	MockedProducts = []product{}
}

func TestAddProduct(t *testing.T) {
	p := []byte(`{
		"name": "Toy",
		"quantity": 10,
		"price": 3.50
	}`)
	r, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(p))
	r.Header.Set("Content-Type", "application/json")
	response := sendRequest(r)
	checkStatusCode(t, 200, response.Code)

	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Toy" {
		t.Errorf("Expected name: %v Got: %v", "Toy", m["name"])
	}
	clearProducts()
}

func TestDeleteProduct(t *testing.T) {
	addProduct()
	p := []byte(`{
		"id": 1,
		"name": "Toy",
		"quantity": 10,
		"price": 3.50
	}`)
	r, _ := http.NewRequest("DELETE", "/product/1", bytes.NewBuffer(p))
	r.Header.Set("Content-Type", "application/json")
	response := sendRequest(r)
	checkStatusCode(t, 200, response.Code)
	r, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(r)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	checkStatusCode(t, 500, response.Code)
	if m["error"] != "not found in mocked data" {
		t.Errorf("Expected: %v Got: %v", 1, m["id"])
	}
	clearProducts()
}

func TestUpdateProduct(t *testing.T) {
	addProduct()
	p := []byte(`{
		"id": 1,
		"name": "Toy 2 updated",
		"quantity": 10,
		"price": 3.50
	}`)
	r, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(p))
	r.Header.Set("Content-Type", "application/json")
	response := sendRequest(r)
	checkStatusCode(t, 200, response.Code)
	r, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(r)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	checkStatusCode(t, 200, response.Code)
	if m["name"] != "Toy 2 updated" {
		t.Errorf("Expected: %v Got: %v", "Toy 2 updated", m["name"])
	}
	clearProducts()
}

func TestGetProduct(t *testing.T) {
	addProduct()
	r, _ := http.NewRequest("GET", "/product/1", nil)
	fmt.Printf(r.URL.Path)
	response := sendRequest(r)
	checkStatusCode(t, 200, response.Code)
	clearProducts()
}

func checkStatusCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected status %v actual %v", strconv.Itoa(expected), strconv.Itoa(actual))
	}
}

func sendRequest(r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, r)
	return recorder
}
