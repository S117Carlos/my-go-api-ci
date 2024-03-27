package main

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

var MockedProducts = []product{
	// {ID: 1, Name: "Toy", Quantity: 1, Price: 2.10},
}

func getProductsList(db *sql.DB) ([]product, error) {
	if db != nil {
		query := "SELECT id, name, quantity, price FROM products"
		rows, err := db.Query(query)
		if err != nil {
			return nil, err
		}
		products := []product{}
		for rows.Next() {
			var p product
			// err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
			err := rows.Scan(&p)
			if err != nil {
				return nil, err
			}
			products = append(products, p)
		}
		return products, nil
	} else {
		return MockedProducts, nil
	}
}

func (p *product) getProductId(db *sql.DB) error {
	if db != nil {
		query := fmt.Sprintf("SELECT id, name, quantity, price FROM products WHERE id = %v", p.ID)
		row := db.QueryRow(query)
		var p product
		err := row.Scan(&p)
		if err != nil {
			return err
		}
		return nil
	} else {
		var find product
		for _, item := range MockedProducts {
			if item.ID == p.ID {
				find = item
				break
			}
		}
		if find.ID == 0 {
			return errors.New("not found in mocked data")
		}
		updateFields(p, find)
		return nil
	}
}

func (p *product) createProduct(db *sql.DB) error {
	if db != nil {
		query := fmt.Sprintf("INSERT INTO (name, quantity, price) products VALUES(%v, %v, %v)", p.Name, p.Quantity, p.Price)
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
		return nil
	} else {
		var maxId int
		if len(MockedProducts) > 0 {
			maxId = MockedProducts[0].ID
		}
		for _, item := range MockedProducts {
			if item.ID == p.ID {
				return errors.New("product already exists")
			}
			if item.ID > maxId {
				maxId = item.ID
			}
		}
		maxId += 1
		p.ID = maxId
		MockedProducts = append(MockedProducts, *p)
		return nil
	}
}

func (p *product) updateProduct(db *sql.DB) error {
	if db != nil {
		query := fmt.Sprintf("UPDATE products SET name = %v, quantity = %v, price = %v WHERE id = %v", p.Name, p.Quantity, p.Price, p.ID)
		res, err := db.Exec(query)
		if err != nil {
			return err
		}
		rowsffected, err := res.RowsAffected()
		if rowsffected == 0 {
			return errors.New("invalid product")
		}
		return err
	} else {
		var found bool
		for i, item := range MockedProducts {
			if item.ID == p.ID {
				MockedProducts[i] = *p
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid product")
		}
		return nil
	}
}

func (p *product) deleteProduct(db *sql.DB) error {
	if db != nil {
		query := fmt.Sprintf("DELETE FROM products WHERE id = %v", p.ID)
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
		return nil
	} else {
		var index int
		var found bool
		for i, item := range MockedProducts {
			if item.ID == p.ID {
				index = i
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid product")
		}
		MockedProducts[index] = MockedProducts[len(MockedProducts)-1]
		MockedProducts = MockedProducts[:len(MockedProducts)-1]
		return nil
	}
}

func updateFields(p1 *product, p2 product) {
	// .Elem() called to dereference the pointer
	aVal := reflect.ValueOf(p1).Elem()
	// aTyp := aVal.Type()

	// no .Elem() called here because it's not a pointer
	bVal := reflect.ValueOf(p2)

	for i := 0; i < aVal.NumField(); i++ {
		// skip the "Age" field:
		// if aTyp.Field(i).Name == "Age" {
		// 	continue
		// }
		// you might want to add some checks here,
		// eg stuff like .CanSet(), to avoid panics
		aVal.Field(i).Set(bVal.Field(i))
	}
}
