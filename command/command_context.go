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

func (context *CommandContext) ReadOption(option string) (string, bool) {
	index, present := context.StringifiedArgs[option]
	if !present || index >= len(context.Args)-1 {
		return "", false
	}

	return string(context.Args[index+1]), true
}
