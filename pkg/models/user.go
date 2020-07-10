package models

type UserMeta struct {
	AFK     bool     `json:"afk"`
	Aliases []string `json:"aliases"`
	IP      string   `json:"ip"`
	Muted   bool     `json:"muted"`
	SMuted  bool     `json:"smuted"`
}

type UserProfile struct {
	Image string `json:"image"`
	Text  string `json:"text"`
}

type User struct {
	Meta    UserMeta    `json:"meta"`
	Name    string      `json:"name"`
	Profile UserProfile `json:"profile"`
	Rank    float64     `json:"rank"`
}

type UpdatedUser struct {
	Name string  `json:"name"`
	Rank float64 `json:"rank"`
}

type UserLeave struct {
	Name string `json:"name"`
}
