package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
)

// MyType is a sample type for testing purposes
type MyType struct {
	A string `json:"a" xml:"AElement" customtag:"noscan"`
	B int64
	M MyTypeInner
	N MyTypeInner `customtag:"world,hello"`
	O MyTypeInner `customtag:"hello,     world" json:"-"`
	P MyType3
}

// MyType2 is a sample type for testing purposes
type MyType2 struct {
	C string `customtag:"decorateme"`
}

// Foo is a sample method
func (m MyType2) Foo() {

}

// MyType3 is a sample type for testing purposes
type MyType3 struct {
	C string `customtag:"hello"`
	D int64  `customtag:"idonothing"`
}

// Foo is a sample method
func (m MyType3) Foo() {

}

// MyTypeInner represents an inner type
type MyTypeInner interface {
	Foo()
}

func main() {
	myVar := MyType{A: "Hello World!", B: 123123123, M: MyType2{"Ouch"}, N: MyType3{"Yieks!", 13}}
	log.Printf("Variable myVar: %v\n", myVar)
	jsonBytes, err := json.MarshalIndent(myVar, "", " ")
	if err != nil {
		log.Fatalf("Error in json.Marshal: %v\n", err)
	}
	log.Printf("JSON representation (for reference):\n%s\n\n", string(jsonBytes))
	printStruct(myVar)
}

func printStruct(input interface{}) {
	inputType := reflect.TypeOf(input)
	inputValue := reflect.ValueOf(input)
	log.Printf("Kind of '%v': %s\n", input, inputType.Kind())

	// In the actual software you write,
	switch inputType.Kind() {
	case reflect.Struct:
		// Debug code: This should never execute
		if inputValue.NumField() != inputType.NumField() {
			log.Fatalf("Mismatch between NumField() of ValueOf (%d) and TypeOf (%d). Exiting!\n", inputValue.NumField(), inputType.NumField())
		}
		for i := 0; i < inputType.NumField(); i++ {
			fieldType := inputType.Field(i)
			fieldValue := inputValue.Field(i)
			tagString := fieldType.Tag.Get("customtag")
			tagLookup := getTagLookup(tagString)
			if tagLookup["noscan"] {
				log.Printf("Encountered field %s with noscan set. Continuing.\n", fieldType.Name)
				continue
			}
			if tagLookup["world"] {
				log.Printf("Encountered field %s with 'world' tag. It's useless\n", fieldType.Name)
			}
			if tagLookup["hello"] {
				log.Printf("Encountered field %s with 'hello' tag. It's useless\n", fieldType.Name)
			}

			decorate := tagLookup["decorateme"]

			if decorate {
				log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
			}
			innerValue := fieldValue.Interface()
			// If current element is assigned via an interface, it could be nil
			// This helps.
			if innerValue != nil {
				printStruct(innerValue)
			}
			if decorate {
				log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
			}
		}
	case reflect.Int64:
		log.Printf("Encountered Int64: %d\n", input)
	case reflect.String:
		log.Printf("Encountered string: '%s'\n", input)
	default:
		// I don't have time to write a general-purpose parser here  so let off
		log.Printf("Unsupported type. reflect.Type value is : %d\n", inputType.Kind())
	}
}

// getTagLookup returns a lookup table for a custom struct tag's options.
// For example, if you set mytag="hello,world", the options "hello" and "world"
// are active according to convention.
func getTagLookup(tagString string) map[string]bool {
	lookup := make(map[string]bool)
	tags := strings.Split(tagString, ",")
	for _, tag := range tags {
		trimmedTag := strings.Trim(tag, " \t")
		lookup[trimmedTag] = true
	}
	return lookup
}
