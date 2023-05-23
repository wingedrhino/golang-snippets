package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// Program should be able to print want to print this:
////////////////////////
// <b>
//    <c>
//       <d>TEXT</d>
//    </c>
// </b>
////////////////////////
const data = `
<a>
	<b>
        <c>
        	<d>TEXT</d>             
        </c>
	</b>
</a>
`

// InnerXML Holds InnerXML as string
type InnerXML struct {
	Value string `xml:",innerxml"`
}

func main() {
	buf := bytes.NewBuffer([]byte(data))
	dec := xml.NewDecoder(buf)
	tagCount := 0
	for {
		t, err := dec.Token()

		if t == nil {
			if err == nil {
				continue
			} else if err == io.EOF {
				break
			} else {
				fmt.Printf("Encountered weird error: %v\n", err)
				break
			}
		}
		if err != nil {
			fmt.Printf("Encountered error while parsing file: %v\n", err)
			break
		}
		switch e := t.(type) {
		case xml.StartElement:
			tagCount++
			// see what element this is
			if e.Name.Local == "a" {
				fmt.Printf("Encountered 'a' tag\n")
				ix := InnerXML{""}
				dec.DecodeElement(&ix, &e)
				fmt.Printf("Decoded Element contents:\n\n %v \n\n", ix)
				err := dec.Skip() // skip until end of token
				if err != nil {
					if err != io.EOF {
						fmt.Printf("Encountered error calling dec.Skio(): %v\n", err)
					}
				}

			}
		}
	}
	fmt.Printf("Number of opening tags encountered: %d\n", tagCount)
}
