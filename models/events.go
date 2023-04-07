package models

// socket.io events
const (
	JoinChannelMethod = "joinChannel"
	LoginMethod       = "login"
	ChatMessage       = "chatMsg"
	OnSetUserRank     = "setUserRank"
	OnAddUser         = "addUser"
	OnUserlist        = "userlist"
	OnUserLeave       = "userLeave"
	OnAfkStateChange  = "setAFK"
)

type Event struct {
	Method  string
	Message interface{}
}

type EventPayload struct {
	Message string      `json:"msg"`
	Meta    interface{} `json:"meta"`
}

func NewChatMessage(msg string) Event {
	return Event{
		Method: ChatMessage,
		Message: EventPayload{
			Message: msg,
			Meta:    struct{}{},
		},
	}
}
