package command

import (
	"fmt"
	"strconv"
)

type CommandContext struct {
	Args            [][]byte
	StringifiedArgs map[string]int
}

func (context *CommandContext) StringArg(index int) string {
	return string(context.Args[index])
}

func (context *CommandContext) IntegerArg(index int) int {
	value, err := strconv.Atoi(context.StringArg(index))
	if err != nil {
		fmt.Println(err.Error())
	}

	return value
}

func (context *CommandContext) HasOption(option string) bool {
	_, present := context.StringifiedArgs[option]
	return present
}

func (context *CommandContext) ReadOption(option string) ([]byte, bool) {
	index, present := context.StringifiedArgs[option]
	if !present || index >= len(context.Args)-1 {
		return make([]byte, 0), false
	}

	return context.Args[index+1], true
}

func (context *CommandContext) ReadOptionAsString(option string) (string, bool) {
	value, present := context.ReadOption(option)
	return string(value), present
}

func (context *CommandContext) ReadOptionAsInt(option string) (int, bool) {
	value, present := context.ReadOptionAsString(option)
	integer, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(err.Error())
	}

	return integer, present
}
