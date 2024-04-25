package bfx

import (
	"slices"
	"strings"
)

func parse(input []byte) ast {
	lexer := new_lexer()
	stream := lexer.ParseBytes(input)
	defer stream.Close()

	program := ast{body: []op{}}
	ctx := [][]op{program.body}
	comment := strings.Builder{}

	for stream.IsValid() {
		field := stream.CurrentToken()

		switch field.Key() {
		case PlusToken:
			ctx = push(ctx, mutate_op{mutation: increment})
		case MinusToken:
			ctx = push(ctx, mutate_op{mutation: decrement})
		case ChevronRightToken:
			ctx = push(ctx, move_op{dir: move_right})
		case ChevronLeftToken:
			ctx = push(ctx, move_op{dir: move_left})
		case OpenBracketsToken:
			ctx = enter(ctx)
		case CloseBracketsToken:
			var popped []op
			ctx, popped = pop(ctx)
			ctx = push(ctx, loop_op{body: popped})
		case DotToken:
			ctx = push(ctx, io_op{io: write})
		case CommaToken:
			ctx = push(ctx, io_op{io: read})
		case StarToken:
			ctx = push(ctx, target_op{})
		case DebugToken:
			ctx = push(ctx, debug_op{})
		default:
			comment.WriteString(field.ValueString())
			comment.WriteString(" ")

			if slices.Contains(program_tokens, stream.NextToken().Key()) {
				ctx = push(ctx, comment_op{raw: comment.String()})
				comment.Reset()
			}
		}

		stream.GoNext()
	}

	last_comment := comment.String()
	if len(last_comment) > 0 {
		ctx = push(ctx, comment_op{last_comment})
		comment.Reset()
	}

	_, program.body = pop(ctx)

	return program
}

func push(ctx [][]op, value op) [][]op {
	current_ctx := ctx[len(ctx)-1]
	current_ctx = append(current_ctx, value)
	ctx[len(ctx)-1] = current_ctx
	return ctx
}

func enter(ctx [][]op) [][]op {
	ctx = append(ctx, []op{})
	return ctx
}

func pop(ctx [][]op) ([][]op, []op) {
	last := ctx[len(ctx)-1]
	ctx = ctx[:len(ctx)-1]

	return ctx, last
}
