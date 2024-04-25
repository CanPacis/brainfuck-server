package bfx

import (
	"fmt"
	"strings"
	"syscall"
)

type Interpreter struct {
	tape    tape
	pointer int
	ast     ast
	fd      int
}

type Cell uint32

func (c Cell) String() string {
	if c == 0 {
		return "0"
	}

	as_string := strings.TrimSpace(string(rune(c)))
	as_byte := fmt.Sprintf("%d", c)

	if len(as_string) > 0 {
		return fmt.Sprintf("%s(%s)", as_byte, as_string)
	}

	return as_byte
}

type tape [30000]Cell

func (t tape) String() string {
	result := []string{}

	for _, item := range t[:50] {
		result = append(result, item.String())
	}

	return strings.Join(result, " ")
}

func (i *Interpreter) run(add_new_line bool) {
	i.run_context(i.ast.body)

	if add_new_line {
		syscall.Write(i.fd, []byte{10})
	}
}

func (i *Interpreter) SetTape(index int, value Cell) {
	i.tape[index] = value
}

func (i *Interpreter) run_context(body []op) {
	for _, op := range body {
		switch op := op.(type) {
		case mutate_op:
			if op.mutation == increment {
				i.tape[i.pointer]++
			} else {
				i.tape[i.pointer]--
			}
		case move_op:
			if op.dir == move_right {
				i.pointer++
			} else {
				i.pointer--
			}

			if i.pointer < 0 {
				i.pointer = 0
			}
			if i.pointer > 30000 {
				i.pointer = 30000
			}
		case loop_op:
			for i.tape[i.pointer] != 0 {
				i.run_context(op.body)
			}
		case io_op:
			if op.io == write {
				syscall.Write(i.fd, []byte{byte(i.tape[i.pointer])})
			} else {
				b := [1]byte{}
				syscall.Read(i.fd, b[:])
				i.tape[i.pointer] = Cell(b[0])
			}
		case target_op:
			i.fd = int(i.tape[i.pointer])
		case debug_op:
			fmt.Printf("Current cell %d = %s\n", i.pointer, i.tape[i.pointer].String())
			fmt.Printf("File descriptor %d\n", i.fd)
			fmt.Println(i.tape)
		}
	}
}

func new_interpreter(program ast) Interpreter {
	return Interpreter{
		tape:    tape{},
		pointer: 0,
		ast:     program,
		fd:      1,
	}
}
