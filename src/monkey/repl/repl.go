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

const (
	PROMPT           = ">> "
	MULTILINE_PROMPT = "... "
	EXIT             = "exit()"
	INTERRUPT        = "^C"
	HISTORY_PATH     = "/Users/anirudhlakkaraju/Programming/go-compiler/src/monkey/repl_history.txt"
)

func REPL(in io.Reader, out io.Writer) {
	// reader := bufio.NewReader(in)
	// env := object.NewEnvironment()

	// Configure readline
	rl, err := readline.NewEx(&readline.Config{
		HistoryFile:     HISTORY_PATH,
		InterruptPrompt: INTERRUPT,
		Prompt:          PROMPT,
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
		if line == EXIT {
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

// isMultilineStart checks for start of multiline input
func isMultilineStart(line string) bool {
	// detect multiline start based on opening braces
	openers := []string{"{", "(", "["}
	for _, opener := range openers {
		if strings.Contains(line, opener) {
			return true
		}
	}
	return false
}

// acceptUntil accepts multiline input until end encountered
func acceptUntil(rl *readline.Instance, start, end string) (string, error) {
	var buf strings.Builder

	buf.WriteString(start)
	buf.WriteRune('\n')
	rl.SetPrompt(MULTILINE_PROMPT)

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

	rl.SetPrompt(PROMPT)

	return buf.String(), nil
}

const MONKEY_FACE = `            __,__
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
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
