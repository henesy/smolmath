package main

import (
	"os"
	"io"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"errors"
	"unicode"
)

// The *kind* of token
type Type int
const (
	Number		Type = iota
	OpenParen
	CloseParen
	Multiply
	Divide
	Add
	Subtract
)

type Tokens []Token

// Represents a token
type Token struct {
	Type
	Value	interface{}
}

// Represents an expression tree
type Tree struct {
}

var (
	maxTokens	uint64
)

// Small infix arithmetic parser
func main() {
	flag.Uint64Var(&maxTokens, "M", 1024, "Maximum tokens per line")

	prompt := flag.String("p", "» ", "Prompt for input")

	flag.Parse()
	
	in := bufio.NewReader(os.Stdin)
	
	loop:
	for {
		fmt.Print(*prompt)
		
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break loop
			}
			
			fatal("err: could read line -", err)
		}
		
		if len(line) < 2 {
			fmt.Println("fail: specify an expression")
			continue loop
		}
		line = line[:len(line)-1]
		
		tokens, err := lex(line)
		if err != nil {
			fmt.Println("lex fail:", err)
			continue loop
		}
		
		fmt.Fprintln(os.Stderr, ">>> DEBUG tokens -", tokens)
		
		tree, err := parse(tokens)
		if err != nil {
			fmt.Println("parse fail:", err)
			continue loop
		}
		
		result := eval(tree)
		
		fmt.Println(result)
	}
}


// Lex input string into tokens
func lex(line string) (Tokens, error) {
	reader := strings.NewReader(line)
	
	parenCount := 0
	tokens := make([]Token, 0, maxTokens)
	num := 0
	inNum := false
	isNeg := false
	lastRune := 'β'
	
	push := func(t Type, dat interface{}) {
			tokens = append(tokens, Token{ t, dat })
		}

	loop:
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break loop
			}
			fatal("err: could not read rune -", err)
		}
		
		// If we leave number space, push what we had
		if inNum && !unicode.IsDigit(r) {
			if isNeg {
				num *= -1
			}
		
			push(Number, num)
			inNum = false
			isNeg = false
			num = 0
		}
		
		switch {
		case unicode.IsDigit(r):
			num *= 10
			num += int(r-'0')
			inNum = true
		
		case r == '-':
			if lastRune == '(' {
				// Only support (-1) to make a negative
				// ex. 1 + (-1) ≡ 0
				isNeg = true
			} else {
				push(Subtract, r)
			}
		
		case r == '+':
			push(Add, r)
		
		case r == '*':
			push(Multiply, r)
		
		case r == '/':
			push(Divide, r)
		
		case r == '(':
			parenCount++
			push(OpenParen, parenCount)
		
		case r == ')':
			parenCount--
			push(CloseParen, parenCount)
		
		case unicode.IsSpace(r):
			// Ignore whitespace
			continue loop
		
		default:
			return nil, errors.New(fmt.Sprintf("unknown token in input → %c", r))
		}
		
		lastRune = r
	}
	
	// Catch a trailing number
	if inNum {
		push(Number, num)
	}
	
	// The only thing allowed to be trailing is a Number or )
	if t := tokens[len(tokens)-1].Type; t != Number && t != CloseParen {
		return nil, errors.New(fmt.Sprintf("invalid trailing type #%d; valid is a Number or ')'", t))
	}
	
	// Maybe remove and let tokenizer give better results?
	if parenCount != 0 {
		abs := func(n int) int {
			if n < 0 {
				return -1*n
			}
			return n
		}
		var c rune
		if parenCount > 0 {
			c = '('
		}
		if parenCount < 0 {
			c = ')'
		}
		return nil, errors.New(fmt.Sprintf("there are %d unmatched '%c'", abs(parenCount), c))
	}

	return tokens, nil
}

// Parse tokens into an expression tree
func parse(tokens []Token) (Tree, error) {
	var tree Tree
	
	// Subtract should be implemented by expanding
	/*
		2-2 → 2 + -1*2
		
		-3 → + -1 * 3
	*/

	return tree, nil
}

// Evaluate an expression tree ­ final result is a string
func eval(tree Tree) string {

	return "nil"
}

// Fatal - end program with an error message and newline
func fatal(s ...interface{}) {
	panic(fmt.Sprintln(s...))
}
