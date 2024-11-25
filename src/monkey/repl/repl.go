package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"go-compiler/src/monkey/compiler"
	"go-compiler/src/monkey/lexer"
	"go-compiler/src/monkey/parser"
	"go-compiler/src/monkey/vm"
)

const PROMPT = ">> "

func REPL(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	// env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)

		input, err := reader.ReadString('\n')
		check(err)

		input = strings.TrimRight(input, " \n")

		// Exit REPL
		if input == "exit()" {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}

		// Allow multiline input for block statements
		if len(input) > 0 && input[len(input)-1] == '{' {
			input, err = acceptUntil(reader, input, "\n\n")
			check(err)
		}

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed: \n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed: \n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")

		// evaluated := evaluator.Eval(program, env)
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func acceptUntil(r *bufio.Reader, start, end string) (string, error) {
	var buf strings.Builder

	buf.WriteString(start)
	buf.WriteRune('\n')
	for {
		fmt.Print("... ")
		line, err := r.ReadString('\n')
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
