package cli

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func PrettyPrintJSON(v any) {
	prettyjson, err := prettyjson.Marshal(v)
	if err != nil {
		fmt.Println("Error marshaling to pretty JSON:", err)
		return
	}
	fmt.Println(string(prettyjson))
}
