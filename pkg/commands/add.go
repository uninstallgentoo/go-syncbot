package commands

import (
	"regexp"
	"strings"
)

type Media struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Video struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Position string `json:"pos"`
	Duration uint   `json:"duration"`
	Temp     bool   `json:"temp"`
}

type addCommand struct{}

func NewAddCommand() CommandExecutor {
	return &addCommand{}
}

func (c *addCommand) Validate(args []string) error {
	return nil
}

func (c *addCommand) GetMinRequiredRank() float64 {
	return 1
}

func (c *addCommand) Exec(args []string) (*CommandResult, error) {
	payloads := make([]*Event, 0, len(args))
	for _, url := range args {
		media := parseMediaLink(url)
		video := Video{
			media.Id,
			media.Type,
			"end",
			0,
			true,
		}
		payloads = append(payloads, &Event{
			Method:  "queue",
			Message: video,
		})
	}
	return NewCommandResult(payloads), nil
}

func lookupMediaId(pattern, url string) string {
	re, _ := regexp.Compile(pattern)
	matches := re.FindStringSubmatch(url)
	if matches != nil {
		return matches[len(matches)-1]
	}
	return ""
}

func newMedia(url, pattern, vtype string) Media {
	media := Media{}
	id := lookupMediaId(pattern, url)
	if id != "" {
		return Media{id, vtype}
	}
	return media
}

func parseMediaLink(url string) Media {
	switch {
	case strings.HasPrefix(url, "jw:"):
		return Media{url[3:], "jw"}
	case strings.HasPrefix(url, "rtmp://"):
		return Media{url, "rt"}
	case strings.Contains(url, "soundcloud.com"):
		return Media{url, "sc"}
	case strings.HasSuffix(url, ".m3u8"):
		return Media{url, "hl"}
	case strings.Contains(url, "youtube.com"):
		pattern := "^http(?:s?):\\//(?:www.)?(youtube.com/watch\\?v=)([^#]+)"
		return newMedia(url, pattern, "yt")
	case strings.Contains(url, "google.com/file"):
		pattern := "^http(?:s?):\\//(?:www.)?(docs.google.com|drive.google.com)/(file/d)/([^/]*)"
		return newMedia(url, pattern, "gd")
	case strings.Contains(url, "twitch.tv"):
		pattern := "^http(?:s?):\\//(?:www.)?(twitch.tv)/([^&#]+)"
		return newMedia(url, pattern, "tw")
	case strings.Contains(url, "vimeo.com"):
		pattern := "^http(?:s?):\\//(?:www.)?(vimeo.com)/([^&#]+)"
		return newMedia(url, pattern, "vi")
	case strings.Contains(url, "dailymotion.com"):
		pattern := "^http(?:s?):\\//(?:www.)?(dailymotion.com/video)/([^&#]+)"
		return newMedia(url, pattern, "dm")
	case strings.Contains(url, "streamable.com"):
		pattern := "^http(?:s?):\\//(?:www.)?(streamable.com)/([^&#]+)"
		return newMedia(url, pattern, "sb")
	default:
		return Media{}
	}
}
