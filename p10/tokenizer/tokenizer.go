package tokenizer

/*
	Check the Top of the Stack.
	case stack = empty:
		push the nextByte
		recur
	case top of stack is number:
		if nextByte == number or decimalpoint
			push
			recur
		else
			append intConst token
			end
	case top of stack is deicmal point:
		if nextByte == number
			push
			recur
		else
			append intConst token
			end
	case top of stack is quotion mark:
		if nextByte != quotation mark
			push
			recur
		else
			append StringConst token
			end
	case top of stack is letter:
*/
import (
	//"bufio"
	//"strings"
	//"fmt"

	"fmt"
	"io"
	"os"

	//"unicode/utf8"
	"github.com/golang-collections/collections/stack"
)

const (
	Keyword = -(iota + 1)
	Symbol
	IntConst
	StrConst
	Identifier
	None
)

var knownSymols []byte = []byte{
	123, 125, // {, }
	40, 41, // (, )
	91, 93, // [, ]
	46, 44, 59, 43, 45, 42, 47, 38, // ., ',', ;, +, -, *, /, &
	124, 60, 62, 61, 126, // |, <, >, =, ~
}

var ClassMap map[int]string = map[int]string{
	Keyword:    "Keyword",
	Symbol:     "Symbol",
	IntConst:   "Integer Constant",
	StrConst:   "String Constant",
	Identifier: "Identifier",
	None:       "None",
}

const (
	NULL     byte = 0
	SPACE    byte = 32
	F_SLASH  byte = 47
	NEW_LINE byte = 10
	TAB      byte = 9
	INT_ZERO byte = 48
	INT_NINE byte = 57
	POINT    byte = 46
	QUOTE    byte = 34
	CAP_A    byte = 65
	CAP_Z    byte = 90
	LOW_A    byte = 97
	LOW_Z    byte = 122
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Token struct {
	Class int
	Value []byte
}

type Reader struct {
	Tokens     []Token
	File       *os.File
	ByteOffset int64 // the ByteOffset used, for os.ReadAt, tracks current position in file
	HasMore    bool
}

func (r *Reader) Advance() error {
	var stack *stack.Stack = stack.New()
	var token Token
	var buff []byte = make([]byte, 1)
	var tokenFound bool = false
	fmt.Println("Token Search Loop")
	for !tokenFound {

		// read from stream
		_, err := r.File.ReadAt(buff, r.ByteOffset)
		if err != nil {
			if err == io.EOF {
				r.HasMore = false
				fmt.Println()
				return nil
			}
			panic(err)
		}
		// label input from stream
		var nextByte byte = buff[0]
		var nextByteIsNumeric bool = nextByte >= INT_ZERO && nextByte <= INT_NINE
		peek := stack.Peek() //Note: I can't assign peek to byte becuase an empty stack returns nil and not 0. This mean I need to assert the byte type in all other comparisons
		// loop through input for a known class
		if peek == nil {
			fmt.Println("\t\tstack is empty")
			// dont push SPACE onto stack
			if nextByte != SPACE && nextByte != NEW_LINE && nextByte != TAB {
				fmt.Printf("\t\t\tpush %d\n", nextByte)
				stack.Push(nextByte)
			}
			r.ByteOffset += 1
		} else {
			fmt.Println("\t\tstack is not-empty")
			var peekIsNumeric bool = peek.(byte) >= INT_ZERO && peek.(byte) <= INT_NINE
			var peekIsQuote bool = peek.(byte) == QUOTE
			var peekIsSymbol bool = false
			var peekIsLetter bool = peek.(byte) >= CAP_A && peek.(byte) <= CAP_Z || peek.(byte) >= LOW_A && peek.(byte) <= LOW_Z
			for i := range knownSymols {
				if knownSymols[i] == peek.(byte) {
					peekIsSymbol = true
				}
			}

			if peekIsNumeric {
				fmt.Println("\t\t\tpeek == numeric")
				if nextByteIsNumeric {
					fmt.Printf("\t\t\tpush %d\n", nextByte)
					stack.Push(nextByte)
					r.ByteOffset += 1
				} else { // populate the token
					fmt.Println("\t\t\tload token")
					count := stack.Len()
					for count > 0 {
						val := stack.Pop().(byte)
						fmt.Printf("\t\t\tpop %d\n", val)
						token.Value = append(token.Value, val)
						count--
					}
					token.Class = IntConst
					tokenFound = true
				}
			} else if peekIsQuote {
				fmt.Println("\t\t\tpeek == quote")
				if nextByte == QUOTE {
					//empty string
					r.ByteOffset += 1
					fmt.Printf("\t\t\t\tpop %d\n", stack.Pop())
					token.Class = StrConst
					token.Value = append(token.Value, NULL)
				} else {
					for nextByte != QUOTE {
						// read from stream
						_, err := r.File.ReadAt(buff, r.ByteOffset)
						if err != nil {
							if err == io.EOF {
								r.HasMore = false
								return nil
							}
							panic(err)
						}
						// label input from stream
						nextByte = buff[0]
						fmt.Printf("\t\t\t\tpush %d\n", nextByte)
						stack.Push(nextByte)
						r.ByteOffset += 1
					}
					// load token
					fmt.Println("\t\t\tload token")
					fmt.Printf("\t\t\t\tpop %d\n", stack.Pop()) // pop starting quote
					count := stack.Len() - 1                    // go until quote
					for count > 0 {
						val := stack.Pop().(byte)
						fmt.Printf("\t\t\t\tpop %d\n", val)
						token.Value = append(token.Value, val)
						count--
					}
					fmt.Printf("\t\t\t\tpop %d\n", stack.Pop()) // Pop last quote
					token.Class = StrConst
					tokenFound = true
				}
			} else if peekIsSymbol {
				fmt.Println("\t\t\tpeek == symbol")
				fmt.Println("\t\t\tload token")
				val := stack.Pop()
				fmt.Printf("\t\t\t\tpop %d\n", val)
				token.Value = append(token.Value, val.(byte))
				token.Class = Symbol
				tokenFound = true
			} else if peekIsLetter {
				var nextByteIsLetter bool = nextByte >= CAP_A && nextByte <= CAP_Z || nextByte >= LOW_A && nextByte <= LOW_Z
				for nextByteIsLetter { // look for a keyword
					fmt.Printf("\t\t\t\tpush %d\n", nextByte)
					stack.Push(nextByte)
					// read from stream
					_, err := r.File.ReadAt(buff, r.ByteOffset)
					if err != nil {
						if err == io.EOF {
							r.HasMore = false
							return nil
						}
						panic(err)
					}
					// label input from stream
					nextByte = buff[0]
					nextByteIsLetter = nextByte >= CAP_A && nextByte <= CAP_Z || nextByte >= LOW_A && nextByte <= LOW_Z
					r.ByteOffset += 1
				}
			}
		}
	}
	fmt.Print("Adding token to list : ")
	fmt.Println(token)
	fmt.Println()
	r.Tokens = append(r.Tokens, token)
	return nil
}
