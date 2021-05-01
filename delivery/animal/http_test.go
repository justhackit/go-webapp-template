package animal

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/justhackit/go-webapp-template/entities"
)

func TestAnimalHandler_Handler(t *testing.T) {
	testcases := []struct {
		method             string
		expectedStatusCode int
	}{
		{"GET", http.StatusOK},
		{"POST", http.StatusOK},
		{"DELETE", http.StatusMethodNotAllowed},
	}

	for _, v := range testcases {
		req := httptest.NewRequest(v.method, "/animal", nil)
		w := httptest.NewRecorder()

		a := New(mockDatastore{})
		a.Handler(w, req)

		if w.Code != v.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", v.expectedStatusCode, w.Code)
		}
	}
}

func TestAnimalGet(t *testing.T) {
	testcases := []struct {
		id       string
		response []byte
	}{
		{"1", []byte("could not retrieve animal")},
		{"1a", []byte("invalid parameter id")},
		{"2", []byte(`[{"ID":2,"Name":"Dog","Age":8}]`)},
		{"0", []byte(`[{"ID":1,"Name":"Ken","Age":23},{"ID":2,"Name":"Dog","Age":8}]`)},
	}

	for i, v := range testcases {
		req := httptest.NewRequest("GET", "/animal?id="+v.id, nil)
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.get(w, req)

		if !reflect.DeepEqual(w.Body, bytes.NewBuffer(v.response)) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Body.String(), string(v.response))
		}
	}
}

func TestAnimalPost(t *testing.T) {
	testcases := []struct {
		reqBody  []byte
		respBody []byte
	}{
		{[]byte(`{"Name":"Hen","Age":12}`), []byte(`could not create animal`)},
		{[]byte(`{"Name":"Maggie","Age":10}`), []byte(`{"ID":12,"Name":"Maggie","Age":10}`)},
		{[]byte(`{"Name":"Maggie","Age":"10"}`), []byte(`invalid body`)},
	}
	for i, v := range testcases {
		req := httptest.NewRequest("GET", "/animal", bytes.NewReader(v.reqBody))
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.create(w, req)

		if !reflect.DeepEqual(w.Body, bytes.NewBuffer(v.respBody)) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Body.String(), string(v.respBody))
		}
	}
}

type mockDatastore struct{}

func (m mockDatastore) Get(id int) ([]entities.Animal, error) {
	if id == 1 {
		return nil, errors.New("db error")
	} else if id == 2 {
		return []entities.Animal{{2, "Dog", 8}}, nil
	}

	return []entities.Animal{{1, "Ken", 23}, {2, "Dog", 8}}, nil
}

func (m mockDatastore) Create(animal entities.Animal) (entities.Animal, error) {
	if animal.Age == 12 {
		return entities.Animal{}, errors.New("db error")
	}

	return entities.Animal{12, "Maggie", 10}, nil
}
