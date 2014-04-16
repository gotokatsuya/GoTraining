package main

import (
	"fmt"
	"encoding/json"
)

const jsonStr = `
{
  "array": [
    1,
    2,
    3
  ],
  "boolean": true,
  "null": null,
  "number": 123,
  "object": {
    "a": "b",
    "c": "d",
    "e": "f"
  },
  "string": "Hello World"
}
`

type partT struct {
	O map[string]string `json:"object,omitempty"`
	S string `json:"string,omitempty"`
}

func main() {
	blob := []byte(jsonStr)

	// decode
	var whole interface{}
	json.Unmarshal(blob, &whole)
	fmt.Printf("whole => %#+v\n", whole)
	var part partT
	json.Unmarshal(blob, &part)
	fmt.Printf("part => %#+v\n", part)

	fmt.Println()

	// encode
	blob, _ = json.Marshal(whole)
	fmt.Printf("whole => %s\n", string(blob))
	blob, _ = json.MarshalIndent(part, "", "  ")
	fmt.Printf("part => %s\n", string(blob))
}
