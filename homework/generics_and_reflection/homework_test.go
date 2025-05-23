package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const (
	PropertiesTag   = "properties"
	OmitemptyTagVal = "omitempty"
)

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](object T) string {
	builder := strings.Builder{}
	objectType := reflect.TypeOf(object)

	for i := range objectType.NumField() {
		fieldValue := reflect.ValueOf(object).Field(i)
		tagValues := strings.Split(objectType.Field(i).Tag.Get(PropertiesTag), ",")

		omitempty := slices.Contains(tagValues, OmitemptyTagVal)
		if len(tagValues) == 0 || fieldValue.IsZero() && omitempty {
			continue
		}

		builder.WriteString(fmt.Sprintf("%s=%v\n", tagValues[0], fieldValue))
	}

	return strings.TrimRight(builder.String(), "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
