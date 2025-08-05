package manager

type ReservedCommand string

const (
	Start  ReservedCommand = "/start"
	Help   ReservedCommand = "/help"
	Debt   ReservedCommand = "/debt"
	Recipe ReservedCommand = "/recipe"
	Gym    ReservedCommand = "/gym"
	Task   ReservedCommand = "/task"
)

var reservedCommands = map[ReservedCommand]struct{}{
	Start:  {},
	Help:   {},
	Debt:   {},
	Recipe: {},
	Gym:    {},
	Task:   {},
}

func ParseCommand(text string) (ReservedCommand, bool) {
	cmd := ReservedCommand(text)
	_, exists := reservedCommands[cmd]
	return cmd, exists
}
