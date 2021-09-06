package models

type Command struct {
	Command string
	Rank    float64
}

type EventPayload struct {
	Message string      `json:"msg"`
	Meta    interface{} `json:"meta"`
}

type Event struct {
	Method  string
	Message interface{}
}

type CommandResult struct {
	Results []*Event
}
