package main

import (
	"fmt"
	"time"

	"github.com/wingedrhino/golang-snippets/elasticmapper"
)

// MyType is a sample type for testing purposes
type MyType struct {
	A string `json:"a" xml:"AElement" elasticmapper:"-"`
	B int64
	M MyType2
	N MyType3 `elasticmapper:"world,hello"`
	O bool    `elasticmapper:"hello,     world" json:"-"`
	P []MyType3
	Q []uint64
	R time.Time
	S []time.Time
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
	fmt.Printf("Variable myVar: %v\n", myVar)
	mapping, err := elasticmapper.GetElasticMapping(myVar, "MyType")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated Mapping:\n\n%s\n\n", mapping)
}
