package command

type CommandContext struct {
	Args            [][]byte
	StringifiedArgs map[string]int
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
