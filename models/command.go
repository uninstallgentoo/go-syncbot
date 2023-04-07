package models

type Command struct {
	Command string
	Rank    float64
}

type CommandResult struct {
	Results []Event
}

func NewCommandResult(events ...Event) CommandResult {
	return CommandResult{Results: events}
}
