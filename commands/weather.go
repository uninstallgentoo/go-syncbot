package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var weatherAPIError = errors.New("Error has occured during request. Try again later.")
var locationBadRequestError = errors.New("You should provide one location to get current weather report.")

type location struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type condition struct {
	Text string `json:"text"`
}

type currentWeather struct {
	TempC     float32   `json:"temp_c"`
	TempF     float32   `json:"temp_f"`
	Condition condition `json:"condition"`
}

type response struct {
	Location location       `json:"location"`
	Current  currentWeather `json:"current"`
}

var Weather = &command.Command{
	Name:        "weather",
	Description: "shows current weather in specified location",
	Rank:        1,
	ExecFunc: func(args []string, cmd *command.Command) (models.CommandResult, error) {
		url := cmd.Config.ExternalServices.Weather.URL
		if val := cmd.Cache.Get(fmt.Sprintf("weather_%s", args[0])); val != nil {
			if str, ok := val.(string); ok {
				return models.NewCommandResult(
					models.NewChatMessage(str),
				), nil
			}
		}
		weatherClient := http.Client{
			Timeout: time.Second * 2,
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return models.CommandResult{}, weatherAPIError
		}

		q := req.URL.Query()
		q.Add("key", cmd.Config.ExternalServices.Weather.Token)
		q.Add("q", args[0])
		q.Add("aqi", "no")
		req.URL.RawQuery = q.Encode()

		res, getErr := weatherClient.Do(req)
		if getErr != nil {
			return models.CommandResult{}, weatherAPIError
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return models.CommandResult{}, weatherAPIError
		}

		resp := response{}
		jsonErr := json.Unmarshal(body, &resp)
		if jsonErr != nil {
			return models.CommandResult{}, weatherAPIError
		}

		msg := fmt.Sprintf(
			"Current temperature in %s, %s: %.0f°C|%.0f°F %s",
			resp.Location.Name,
			resp.Location.Region,
			resp.Current.TempC,
			resp.Current.TempF,
			resp.Current.Condition.Text,
		)
		defer cmd.Cache.Set(
			fmt.Sprintf("weather_%s", args[0]),
			msg,
			time.Duration(time.Hour),
		)
		return models.NewCommandResult(
			models.NewChatMessage(msg),
		), nil
	},
	ValidateFunc: func(args []string) error {
		if len(args) != 1 {
			return locationBadRequestError
		}
		return nil
	},
}
