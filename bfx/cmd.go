package bfx

import "os"

func Run(program []byte, args ...Cell) error {
	interpreter := new_interpreter(parse(program))

	for i, arg := range args {
		interpreter.SetTape(i, arg)
	}

	interpreter.run(false)
	return nil
}

func RunFile(path string, args ...Cell) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return Run(file, args...)
}
