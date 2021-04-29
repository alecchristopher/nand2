package main

import (
	"fmt"
	"os"
	"p10/tokenizer"
)

func main() {
	src := "../lib/number.jack"
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	t := tokenizer.Reader{make([]tokenizer.Token, 0), file, 0}
	var hasMore = true
	for hasMore {
		_, err := t.Advance()
		if err != nil {
			hasMore = false
		}
	}
	for i := range t.Tokens {
		fmt.Print("<")
		fmt.Print(tokenizer.ClassMap[t.Tokens[i].Class])
		fmt.Print("> ")
		fmt.Print(t.Tokens[i].Value)
		fmt.Print("<")
		fmt.Print(tokenizer.ClassMap[t.Tokens[i].Class])
		fmt.Print("> ")
		fmt.Println()
	}
	//fmt.Println(t.Tokens)
}
