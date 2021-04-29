package main

import (
	"fmt"
	"os"
	"p10/tokenizer"
)

func main() {
	src := "../lib/symbol.jack"
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}

	t := tokenizer.Reader{make([]tokenizer.Token, 0), file, 0, true}

	for t.HasMore {
		err := t.Advance()
		if err != nil {
			panic(err)
		}
	}

	/*var hasMore = true
	for hasMore {
		_, err := t.Advance()
		if err != nil {
			hasMore = false
		}
	}*/

	for i := range t.Tokens {
		fmt.Printf("<%s> ", tokenizer.ClassMap[t.Tokens[i].Class])
		j := len(t.Tokens[i].Value) - 1
		for j >= 0 {
			fmt.Print(string(t.Tokens[i].Value[j]))
			j--
		}
		fmt.Println()
	}
}
