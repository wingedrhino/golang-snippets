package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
)

const (
	structTag = "elasticmapper"
)

// MyType is a sample type for testing purposes
type MyType struct {
	A string `json:"a" xml:"AElement" elasticmapper:"-"`
	B int64
	M MyType2
	N MyType3 `elasticmapper:"world,hello"`
	O bool    `elasticmapper:"hello,     world" json:"-"`
	P MyType3
}

// MyType2 is a sample type for testing purposes
type MyType2 struct {
	C string `elasticmapper:"text"`
}

// MyType3 is a sample type for testing purposes
type MyType3 struct {
	C string `elasticmapper:"hello"`
	D int64  `elasticmapper:"idonothing"`
}

func main() {
	myVar := MyType{A: "Hello World!", B: 123123123, M: MyType2{"Ouch"}, N: MyType3{"Yieks!", 13}}
	log.Printf("Variable myVar: %v\n", myVar)
	jsonTree := gabs.New()
	GetElasticMapping(myVar, jsonTree, "myVar", false)
	fmt.Println(jsonTree.StringIndent("", "  "))
}

// GetElasticMapping returns an ElasticSearch mapping that's properly populated
// in the provided *gabs.Contiainer.
// TODO: Follow `json` tags and handle pointers.
func GetElasticMapping(input interface{}, jsonTree *gabs.Container, structName string, isText bool) {
	// Handle special case of us defining time.Time
	if t, ok := (input).(time.Time); ok {
		GetElasticMapping(t.Unix(), jsonTree, structName, false)
		return
	}
	inputType := reflect.TypeOf(input)

	switch inputType.Kind() {
	case reflect.Array, reflect.Slice:
		innerType := inputType.Elem()
		zeroVal := reflect.Zero(innerType).Interface()
		if innerType.Kind() == reflect.Struct {
			childTreeContainer, err := jsonTree.Object(structName)
			if err != nil {
				log.Fatal(err)
			}
			// Store arrays of structs as nested objects in ElasticSearch
			childTreeContainer.Set("nested", "type")
			childTree, err := childTreeContainer.Object("properties")
			if err != nil {
				log.Fatal(err)
			}
			GetElasticMapping(zeroVal, childTree, "QWWEW5afasfsddfsf", false)

		} else {
			// Pass along
			childTreeContainer, err := jsonTree.Object(structName)
			if err != nil {
				log.Fatalf("Error creating container for inner struct. Exiting... %v\n", err)
				return
			}
			GetElasticMapping(zeroVal, childTreeContainer, structName, false)
		}
	case reflect.Struct:
		for i := 0; i < inputType.NumField(); i++ {
			fieldType := inputType.Field(i)
			tagString := fieldType.Tag.Get(structTag)
			tagLookup := GetTagLookup(tagString)
			if tagLookup["-"] {
				log.Printf("Encountered field %s with '-'set. Ignoring field!.\n", fieldType.Name)
				continue
			}
			fieldValue := reflect.Zero(fieldType.Type)
			fieldKind := fieldValue.Type().Kind()
			childTreeContainer, err := jsonTree.Object(fieldType.Name)
			if err != nil {
				log.Fatalf("Error creating container for inner struct. Exiting... %v\n", err)
				return
			}
			innerValue := fieldValue.Interface()
			if fieldKind == reflect.Struct { // cannot support interfaces
				_, err = childTreeContainer.Set("object", "type")
				if err != nil {
					log.Fatalf("Error setting type to object. Exiting... %v\n", err)
				}
				childTree, err := childTreeContainer.Object("properties")
				if err != nil {
					log.Fatalf("Error creating properties holder for inner struct. Exiting...%v\n", err)
				}
				// do note that fieldType.Name is useless if innerValue is a struct
				GetElasticMapping(innerValue, childTree, "asfdADAdASD", false)
			} else {
				if tagLookup["text"] {
					GetElasticMapping(innerValue, childTreeContainer, fieldType.Name, true)
				} else {
					GetElasticMapping(innerValue, childTreeContainer, fieldType.Name, false)
				}
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		jsonTree.Set("long", "type")
	case reflect.String:
		if isText {
			jsonTree.Set("text", "type")
		} else {
			jsonTree.Set("keyword", "type")
		}
	case reflect.Bool:
		jsonTree.Set("boolean", "type")
	case reflect.Float32, reflect.Float64:
		jsonTree.Set("double", "type")
	case reflect.Ptr:
		log.Printf("This is a pointer. WTF? '%v'\n", input)
	default:
		// Will have to handle edge cases properly later
		log.Printf("Unsupported type. reflect.Type value is : %d\n", inputType.Kind())
	}
}

// GetTagLookup returns a lookup table for a custom struct tag's options.
// For example, if you set mytag="hello,world", the options "hello" and "world"
// are active according to convention.
func GetTagLookup(tagString string) map[string]bool {
	lookup := make(map[string]bool)
	tags := strings.Split(tagString, ",")
	for _, tag := range tags {
		trimmedTag := strings.Trim(tag, " \t")
		lookup[trimmedTag] = true
	}
	return lookup
}
