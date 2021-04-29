package tokenizer

import (
	//"bufio"
	//"strings"
	//"fmt"
	"os"
	//"unicode/utf8"
	"github.com/golang-collections/collections/stack"
)

const(
	Keyword = -(iota + 1)
	Symbol
	IntConst
	StrConst
	Identifier
	None
)

const(
	SPACE byte = 32
	F_SLASH byte = 47
	NEW_LINE byte = 10
	TAB byte = 9
	INT_ZERO byte = 48
	INT_NINE byte = 56
	POINT byte = 46
)

func check(e error){
	if e != nil {
		panic(e)
	}
}

type Token struct{
	Class int
	Value []byte
}

type Reader struct{
	Tokens []Token
	File  *os.File
	CurrentByte int
}

func (r *Reader) Advance() ( Token, error){
	var t Token = Token{None, make([]byte, 0)}
	var valueStack *stack.Stack = stack.New()
	return t, rAdvance(r, &t, valueStack)
}

func rAdvance(r *Reader, token *Token, stack *stack.Stack) ( error ){
	b :=  make([]byte, 1) // create byte array to store n=1 amonut of bytes
	_, err := r.File.ReadAt(b, int64(r.CurrentByte))
	if err != nil { return err }
	// byte read from file
	nextByte := b[0]
	// DEBUG
	//fmt.Println(nextByte)
	// Debug
	//ddinput, _ := utf8.DecodeRune(b)
	//fmt.Printf("%x", b[0])
	//fmt.Print("\n\t")
	// Base Case
	peeked := stack.Peek()
	//fmt.Println(stack)
	if ( peeked == nil ){ // empty stack
		stack.Push(nextByte)
		r.CurrentByte += 1
		return rAdvance(r, token, stack)
	} else if ( peeked.(byte) >= INT_ZERO && peeked.(byte) <= INT_NINE) { // peeked == number
		if( nextByte >=  INT_ZERO && nextByte <= INT_NINE || nextByte == POINT){ //nextByte == number || point
			stack.Push(nextByte)
			r.CurrentByte += 1
			return rAdvance(r, token, stack)
		} else { // create a return token
			temp := stack.Pop()
			for temp != nil {
				token.Value = append(token.Value, temp.(byte))
				temp = stack.Pop()
			}
			r.Tokens = append(r.Tokens, *token)
			//Debug
			//fmt.Println()
			//return nil
		}
	} else if (peeked.(byte) == POINT) {
		if( nextByte == POINT || (nextByte > INT_NINE && nextByte < INT_ZERO)){
			panic("syntax error: double or trailling decimal point")
		} else if (nextByte >= INT_ZERO && nextByte <= INT_NINE){
			stack.Push(nextByte)
			r.CurrentByte += 1
			return rAdvance(r, token, stack)
		} else {
			temp := stack.Pop()
			for temp != nil {
				token.Value = append(token.Value, temp.(byte))
				temp = stack.Pop()
			}
			r.Tokens = append(r.Tokens, *token)
		}
	}
	// Debug
	//fmt.Println()
	return nil
}

/*
func (r *Reader) rAdvance() ( error ){
	b := make([]byte, 1)
	s :=  stack.New()
	_, err := r.File.Read(b)
	if err != nil{
		return err
	}
	/*v, _ := utf8.DecodeRune(b)
	fmt.Print(v)
	fmt.Print(" : ")
	// check space
	val := b[0]
	if val == SPACE || val == NEW_LINE || val == TAB {
		return r.Advance()
	}
	// check number
	if ((val >= INT_ZERO && val <= INT_NINE) || val == POINT){
		s.Push(val)

	}
	/*else if val == F_SLASH {
		temp := make([]byte, 1)
		_, err := r.File.ReadAt(temp, int64(r.CurrentByte + 1))
		if err != nil {
			return err
		}
		if(temp[0] == F_SLASH){
			for val != NEW_LINE {
				_, err := r.File.Read(b)
				if err != nil {
					return err
				}
				val = b[0]
				r.CurrentByte += 1
			}
			return r.Advance()
		}
	}
	//classify(b[0])
	r.Tokens = append(r.Tokens, Token{0, b[0]})
	r.CurrentByte += 1
	return nil
}*/



