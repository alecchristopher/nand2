package main

import (
	"os"
	"p10/tokenizer"
	"fmt"
)

func main() {
	src := "../lib/number.jack"
	file, err := os.Open(src)
	if err != nil{
		panic(err)
	}
	t := tokenizer.Reader{make([]tokenizer.Token, 0), file, 0}
	var hasMore = true
	for hasMore{
		v, err := t.Advance()
		fmt.Println(v)
		if err != nil {
			hasMore = false
		}
	}
	fmt.Println(t.Tokens)
}
