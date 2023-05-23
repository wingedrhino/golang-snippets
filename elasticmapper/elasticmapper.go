package elasticmapper

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
)

const (
	structTag = "elasticmapper"
)

// GetElasticMapping returns an ElasticSearch mapping from a struct
// Arguments: input, which could be any struct and typeName, which is the name
// of the type you wish to set in ElasticSearch
func GetElasticMapping(input interface{}, typeName string) (mapping string, err error) {
	mappingTree := gabs.New()
	rootTree, err := mappingTree.Object("mappings", typeName, "properties")
	if err != nil {
		return
	}
	err = getElasticMappingImpl(input, rootTree, typeName, false)
	mapping = mappingTree.StringIndent("", "  ")
	return
}

func getElasticMappingImpl(input interface{}, jsonTree *gabs.Container, structName string, isText bool) error {
	// Handle special case of us defining time.Time
	if t, ok := (input).(time.Time); ok {
		return getElasticMappingImpl(t.Unix(), jsonTree, structName, false)
	}
	inputType := reflect.TypeOf(input)

	switch inputType.Kind() {
	case reflect.Array, reflect.Slice:
	case reflect.Struct:
		for i := 0; i < inputType.NumField(); i++ {
			fieldType := inputType.Field(i)
			tagString := fieldType.Tag.Get(structTag)
			tagLookup := getTagLookup(tagString)
			if tagLookup["-"] {
				log.Printf("Encountered field %s with '-'set. Ignoring field!.\n", fieldType.Name)
				continue
			}
			fieldValue := reflect.Zero(fieldType.Type)
			fieldKind := fieldValue.Type().Kind()
			innerValue := fieldValue.Interface()
			// Handle special case of us defining time.Time
			if _, ok := (innerValue).(time.Time); ok {
				_, err := jsonTree.Set("date", fieldType.Name, "type")
				if err != nil {
					return err
				}
				continue
			}
			if fieldKind == reflect.Struct { // cannot support interfaces
				childTreeContainer, err := jsonTree.Object(fieldType.Name)
				if err != nil {
					log.Println(err)
					return err
				}
				_, err = childTreeContainer.Set("object", "type")
				if err != nil {
					log.Println(err)
					return err
				}
				childTree, err := childTreeContainer.Object("properties")
				if err != nil {
					log.Println(err)
					return err
				}
				// do note that fieldType.Name is useless if innerValue is a struct
				err = getElasticMappingImpl(innerValue, childTree, "asfdADAdASD", false)
				if err != nil {
					log.Println(err)
					return err
				}
			} else if fieldKind == reflect.Array || fieldKind == reflect.Slice {
				// TODO we don't yet process timestamps properly
				innerType := fieldType.Type.Elem()
				zeroVal := reflect.Zero(innerType).Interface()
				if innerType.Kind() == reflect.Struct {
					// Handle special case of us defining time.Time
					if _, ok := (zeroVal).(time.Time); ok {
						_, err := jsonTree.Set("date", fieldType.Name, "type")
						if err != nil {
							return err
						}
						goto endArray
					}
					childTreeContainer, err := jsonTree.Object(fieldType.Name)
					if err != nil {
						log.Println(err)
						return err
					}
					// Store arrays of structs as nested objects in ElasticSearch
					childTreeContainer.Set("nested", "type")
					childTree, err := childTreeContainer.Object("properties")
					if err != nil {
						log.Println(err)
						return err
					}
					err = getElasticMappingImpl(zeroVal, childTree, "QWWEW5afasfsddfsf", false)
					if err != nil {
						log.Println(err)
						return err
					}
				} else {
					// Pass along
					isText = tagLookup["text"]
					childTreeContainer, err := jsonTree.Object(fieldType.Name)
					if err != nil {
						log.Println(err)
						return err
					}
					err = getElasticMappingImpl(zeroVal, childTreeContainer, fieldType.Name, isText)
					if err != nil {
						log.Println(err)
						return err
					}
				}
			endArray:
			} else {
				isText = tagLookup["text"]
				childTreeContainer, err := jsonTree.Object(fieldType.Name)
				if err != nil {
					log.Println(err)
					return err
				}
				err = getElasticMappingImpl(innerValue, childTreeContainer, fieldType.Name, isText)
				if err != nil {
					log.Println(err)
					return err
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
	return nil
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
