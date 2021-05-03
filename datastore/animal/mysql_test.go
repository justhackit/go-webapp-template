package animal

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/justhackit/go-webapp-template/driver"
	"github.com/justhackit/go-webapp-template/entities"
)

func initializeMySQL(t *testing.T) *sql.DB {
	conf := driver.MySQLConfig{
		Host:     os.Getenv("SQL_HOST"),
		User:     os.Getenv("SQL_USER"),
		Password: os.Getenv("SQL_PASSWORD"),
		Port:     os.Getenv("SQL_PORT"),
		Db:       os.Getenv("SQL_DB"),
	}

	var err error
	db, err := driver.ConnectToMySQL(conf)
	if err != nil {
		t.Errorf("could not connect to sql, err:%v", err)
	}

	return db
}

func TestDatastore(t *testing.T) {
	db := initializeMySQL(t)
	a := New(db)
	testAnimalStorer_Get(t, a)
	testAnimalStorer_Create(t, a)

}

func testAnimalStorer_Create(t *testing.T, db AnimalDAO) {
	testcases := []struct {
		req      entities.Animal
		response entities.Animal
	}{
		{entities.Animal{Name: "Hen", Age: 1}, entities.Animal{ID: 3, Name: "Hen", Age: 1}},
		{entities.Animal{Name: "Pig", Age: 2}, entities.Animal{ID: 4, Name: "Pig", Age: 2}},
	}
	for i, v := range testcases {
		resp, _ := db.Create(v.req)

		if !reflect.DeepEqual(resp, v.response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.response)
		}
	}
}

func TestAnimalStorer_CreateMock(t *testing.T) {
	dbConn, mock, _ := sqlmock.New()
	db := New(dbConn)
	testcase := entities.Animal{Name: "Peppa Pig", Age: 2}

	mock.ExpectExec("INSERT INTO animal").WithArgs("Peppa Pig", 2).WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := db.Create(testcase)
	if err != nil {
		t.Errorf("Test Failed. Error occured : %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were NOT met : %v", err)
	}

}

func testAnimalStorer_Get(t *testing.T, db AnimalDAO) {
	testcases := []struct {
		id   int
		resp []entities.Animal
	}{
		{0, []entities.Animal{{1, "Hippo", 10}, {2, "Ele", 20}}},
		{1, []entities.Animal{{1, "Hippo", 10}}},
	}
	for i, v := range testcases {
		resp, _ := db.Get(v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.resp)
		}
	}
}
