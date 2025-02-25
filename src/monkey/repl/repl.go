// Package repl provides an interactive Read-Eval-Print Loop for the Monkey language.
package repl

import (
	"fmt"
	"go-compiler/src/monkey/compiler"
	"go-compiler/src/monkey/lexer"
	"go-compiler/src/monkey/object"
	"go-compiler/src/monkey/parser"
	"go-compiler/src/monkey/vm"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// Constants used in REPL
const (
	Prompt          = ">> "
	MultilinePrompt = "... "
	Exit            = "exit()"
	Interrupt       = "^C"

	HistoryPath = "/Users/anirudhlakkaraju/Programming/go-compiler/src/monkey/repl_history.txt"
)

// REPL starts the input output loop
func REPL(_ io.Reader, out io.Writer) {
	// reader := bufio.NewReader(in)
	// env := object.NewEnvironment()

	// Configure readline
	rl, err := readline.NewEx(&readline.Config{
		HistoryFile:     HistoryPath,
		InterruptPrompt: Interrupt,
		Prompt:          Prompt,
	})
	check(err)
	defer rl.Close()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	// History buffer
	history := make([]string, 0)

	for {
		// Read Input
		line, err := rl.Readline()
		check(err)

		line = strings.TrimSpace(line)
		if line == Exit {
			fmt.Println("Goodbye!")
			return
		}

		// Allow multiline input for block statements
		if isMultilineStart(line) {
			line, err = acceptUntil(rl, line, "\n\n")
			check(err)
		}

		history = append(history, line)
		processInput(line, constants, globals, symbolTable, out)
	}
}

func check(err error) {
	if err == readline.ErrInterrupt {
		fmt.Println("Goodbye!")
		os.Exit(0)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// processInput parses and executes Monkey Program
func processInput(input string, constants []object.Object, globals []object.Object, symbolTable *compiler.SymbolTable, out io.Writer) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	comp := compiler.NewWithState(symbolTable, constants)
	err := comp.Compile(program)
	if err != nil {
		fmt.Fprintf(out, "Whoops! Compilation failed: \n %s\n", err)
		return
	}

	code := comp.Bytecode()
	constants = code.Constants

	machine := vm.NewWithGlobalsStore(code, globals)
	err = machine.Run()
	if err != nil {
		fmt.Fprintf(out, "Whoops! Executing bytecode failed: \n %s\n", err)
		return
	}

	stackTop := machine.LastPoppedStackElem()
	io.WriteString(out, stackTop.Inspect())
	io.WriteString(out, "\n")

	// evaluated := evaluator.Eval(program, env)
	// if evaluated != nil {
	// 	io.WriteString(out, evaluated.Inspect())
	// 	io.WriteString(out, "\n")
	// }
}

// isMultilineStart checks if the line ends with an unclosed bracket
func isMultilineStart(line string) bool {
	stack := []rune{}
	bracketPairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range line {
		// If it's an opening bracket, push to the stack
		if char == '{' || char == '(' || char == '[' {
			stack = append(stack, char)
		}

		// If it's a closing bracket, check if it matches the top of the stack
		if char == '}' || char == ')' || char == ']' {
			if len(stack) > 0 && stack[len(stack)-1] == bracketPairs[char] {
				stack = stack[:len(stack)-1] // Pop the stack
			} else {
				// Mismatched closing bracket
				return false
			}
		}
	}

	// If the stack is not empty, there are unclosed brackets
	return len(stack) > 0
}

// acceptUntil accepts multiline input until end encountered
func acceptUntil(rl *readline.Instance, start, end string) (string, error) {
	var buf strings.Builder

	buf.WriteString(start)
	buf.WriteRune('\n')
	rl.SetPrompt(MultilinePrompt)

	for {
		line, err := rl.Readline()
		if err != nil {
			return "", err
		}

		line = strings.TrimRight(line, " \n")
		buf.WriteString(line)
		buf.WriteRune('\n')

		if s := buf.String(); len(s) > len(end) && s[len(s)-len(end):] == end {
			break
		}
	}

	rl.SetPrompt(Prompt)

	return buf.String(), nil
}

// MonkeyFace during oopsies
const MonkeyFace = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MonkeyFace)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
