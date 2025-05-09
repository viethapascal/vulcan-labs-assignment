package db

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"testing"
)

type QueryableData struct {
	Value int `json:"value"`
}

func (q *QueryableData) Bytes() ([]byte, error) {
	marshal, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
func (q *QueryableData) FromBytes(input []byte) error {
	return json.Unmarshal(input, q)
}
func TestInitDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		panic(err)
	}
	dat := &QueryableData{rand.Int()}
	err = SetKey("data", dat)
	if err != nil {
		panic(err)
	}
	// Get data
	result := new(QueryableData)
	err = GetKey("data", result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Query data: %+v\n", result)

}
