package datastore

import "github.com/justhackit/go-webapp-template/entities"

type Animal interface {
	Get(id int) ([]entities.Animal, error)
	Create(entities.Animal) (entities.Animal, error)
}
