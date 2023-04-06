package command

type CommandContext struct {
	Args            [][]byte
	StringifiedArgs map[string]int
}

func (context *CommandContext) HasOption(option string) bool {
	_, result := context.StringifiedArgs[option]
	return result
}

func (context *CommandContext) GetOption(option string) (string, bool) {
	index, result := context.StringifiedArgs[option]
	if !result || index >= len(context.Args)-1 {
		return "", false
	}

	return string(context.Args[index+1]), true
}
