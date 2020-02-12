package main

import (
	"os"
	"io"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"errors"
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
func lex(line string) ([]Token, error) {
	fmt.Println(line)
	reader := strings.NewReader(line)
	
	parenCount := 0
	tokens := make([]Token, 0, maxTokens)
	num := 0
	inNum := false
	
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
		if inNum && r < '9' && r > '0' {
			push(Number, num)
			inNum = false
		}
		
		switch {
		case r >= '0' && r <= '9':
			num *= 10
			num += int(r-'0')
			inNum = true
		
		case r == '-':
			push(Subtract, r)
		
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
		
		default:
			return nil, errors.New(fmt.Sprintf("unknown token in input → %c", r))
		}
	}
	
	// Catch a trailing number
	if inNum {
		push(Number, num)
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
