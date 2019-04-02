package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func PrintJSon(v interface{}) {
	b, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
