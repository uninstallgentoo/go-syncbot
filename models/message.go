package models

import (
	"regexp"
	"strings"
)

type Message struct {
	Username string                 `json:"username"`
	Image    string                 `json:"image"`
	Text     string                 `json:"msg"`
	Time     int64                  `json:"time"`
	Meta     map[string]interface{} `json:"meta"`
}

func (m *Message) Clean() *Message {
	s := strings.ReplaceAll(m.Text, "&#39;", "'")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#40;", "(")
	s = strings.ReplaceAll(s, "&#41;", ")")
	// clear links from tags
	re, _ := regexp.Compile("<[^>]*>")
	s = re.ReplaceAllString(s, "")
	return &Message{
		Username: m.Username,
		Image:    m.Image,
		Text:     s,
		Time:     m.Time,
		Meta:     m.Meta,
	}
}
