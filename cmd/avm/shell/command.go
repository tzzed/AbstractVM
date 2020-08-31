package shell

import (
	"fmt"
)

type Command struct {
	name string
	opts string
	help string
}

type Instructions struct {
	cmds []Command
}

var instructions Instructions

func registerCommands() {
	instructions.cmds = append(instructions.cmds, Command{name: "assert", opts: "value", help: "Verify that the value at the top of the stack is equal to the one passed as parameter in this instruction"})
	instructions.cmds = append(instructions.cmds, Command{name: "add", help: "Unstack the first two values in the stack, add them, and then stack the result."})
	instructions.cmds = append(instructions.cmds, Command{name: "push", opts: "value", help: "Stack the v value at the top."})
	instructions.cmds = append(instructions.cmds, Command{name: "pop", help: "Unstack the value at the top of the stack."})
	instructions.cmds = append(instructions.cmds, Command{name: "div", help: "Unstack the first two values in the stack, divide them."})
	instructions.cmds = append(instructions.cmds, Command{name: "mod", help: "Unstack the first two values in the stack, calculate their modulo."})
	instructions.cmds = append(instructions.cmds, Command{name: "mul", help: "Unstack the first two values in the stack, multiply them."})
	instructions.cmds = append(instructions.cmds, Command{name: "sub", help: "Unstack the first two values in the stack, substract them."})
}

func displayHelpCommand() error {
	commands := instructions.cmds
	for _, c := range commands {
		indent := 15 - len(c.opts) - len(c.name)
		fmt.Printf("%s %s", c.name, c.opts)
		fmt.Printf("%*s%s\n", indent, "", c.help)
	}

	return nil
}

func getCommands() []Command {
	return instructions.cmds
}

func getAllOperands() []string {
	return []string{"int8",
		"int16",
		"int32",
		"float",
		"double"}
}
