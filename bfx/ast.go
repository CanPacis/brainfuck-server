package bfx

import (
	"strings"
)

type op_kind string

const (
	move_op_kind    op_kind = "op:move"
	mutate_op_kind  op_kind = "op:mutate"
	io_op_kind      op_kind = "op:io"
	loop_op_kind    op_kind = "op:loop"
	target_op_kind  op_kind = "op:target"
	comment_op_kind op_kind = "op:comment"
	debug_op_kind   op_kind = "op:debug"
)

type op interface {
	kind() op_kind
	String() string
}

type ast struct {
	body []op
}

func (ast ast) String() string {
	builder := strings.Builder{}

	for _, op := range ast.body {
		builder.WriteString(op.String())
	}

	return builder.String()
}

type dir string

const (
	move_left  dir = "move:left"
	move_right dir = "move:right"
)

type move_op struct {
	dir dir
}

func (op move_op) kind() op_kind {
	return move_op_kind
}

func (op move_op) String() string {
	if op.dir == move_right {
		return ">"
	}

	return "<"
}

type mutation string

const (
	increment mutation = "mutation:increment"
	decrement mutation = "mutation:decrement"
)

type mutate_op struct {
	mutation mutation
}

func (op mutate_op) kind() op_kind {
	return mutate_op_kind
}

func (op mutate_op) String() string {
	if op.mutation == increment {
		return "+"
	}

	return "-"
}

type io_kind string

const (
	write io_kind = "io:write"
	read  io_kind = "io:read"
)

type io_op struct {
	io io_kind
}

func (op io_op) kind() op_kind {
	return io_op_kind
}

func (op io_op) String() string {
	if op.io == write {
		return "."
	}

	return ","
}

type target_op struct{}

func (op target_op) kind() op_kind {
	return target_op_kind
}

func (op target_op) String() string {
	return "*"
}

type loop_op struct {
	body []op
}

func (op loop_op) kind() op_kind {
	return loop_op_kind
}

func (op loop_op) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")

	for _, op := range op.body {
		builder.WriteString(op.String())
	}

	builder.WriteString("]")
	return builder.String()
}

type comment_op struct {
	raw string
}

func (op comment_op) kind() op_kind {
	return comment_op_kind
}

func (op comment_op) String() string {
	return op.raw
}

type debug_op struct{}

func (op debug_op) kind() op_kind {
	return debug_op_kind
}

func (op debug_op) String() string {
	return "debug"
}
