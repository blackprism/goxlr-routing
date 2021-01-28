package Config

import (
	"encoding/json"
	"fmt"
	"github.com/blackprism/goxlr-routing/GOXLR"
	"io/ioutil"
)

type Keybinding struct {
	Id       int16
	Name     string
	Triggers []Trigger
}

type Trigger struct {
	Type    string
	Actions []Action
}

type Action struct {
	Name   GOXLR.ActionType
	Input  GOXLR.Input
	Output GOXLR.Output
}

func Json(file string) []Keybinding {
	type Root []struct {
		Keybinding string `json:"keybinding"`
		Triggers   []struct {
			Type    string `json:"type"`
			Actions []struct {
				Name    string `json:"name"`
				Targets []struct {
					Inputs  []string `json:"inputs"`
					Outputs []string `json:"outputs"`
				} `json:"targets"`
			} `json:"actions,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"triggers"`
	}

	jsonConfig, err := ioutil.ReadFile(file)

	if err != nil {
		panic(fmt.Sprintf("File %s can't be opened", file))
	}

	var root Root
	json.Unmarshal(jsonConfig, &root)

	keybindings := []Keybinding{}

	for index, keybindingSetup := range root {
		var tmpTriggers []Trigger
		for _, trigger := range keybindingSetup.Triggers {
			if trigger.Type == "routing" {
				var tmpActions []Action
				for _, action := range trigger.Actions {
					for _, target := range action.Targets {
						for _, input := range target.Inputs {
							for _, output := range target.Outputs {
								tmpActions = append(tmpActions, Action{
									actionToGoXLR(action.Name),
									inputToGoXLR(input),
									outputToGoXLR(output),
								})
							}
						}
					}
				}
				tmpTriggers = append(tmpTriggers, Trigger{
					trigger.Type,
					tmpActions,
				})
			}
		}
		keybindings = append(keybindings, Keybinding{
			Id:       int16(index) + 1,
			Name:     keybindingSetup.Keybinding,
			Triggers: tmpTriggers,
		})
	}

	return keybindings
}

func actionToGoXLR(action string) GOXLR.ActionType {
	switch action {
	case "toggle":
		return GOXLR.Toggle
	}

	panic(fmt.Sprintf("Action %s can't be converted to GOXLR format", action))
}

func inputToGoXLR(input string) GOXLR.Input {
	switch input {
	case "Mic":
		return GOXLR.Mic
	case "Chat":
		return GOXLR.Chat
	case "Music":
		return GOXLR.Music
	case "Game":
		return GOXLR.Game
	case "Console":
		return GOXLR.Console
	case "Line In":
		return GOXLR.LineIn
	case "System":
		return GOXLR.System
	case "Samples":
		return GOXLR.Samples
	}

	panic(fmt.Sprintf("Input %s can't be converted to GOXLR format", input))
}

func outputToGoXLR(output string) GOXLR.Output {
	switch output {
	case "Headphones":
		return GOXLR.Headphones
	case "Broadcast Stream Mix":
		return GOXLR.BroadcastMix
	case "Line Out":
		return GOXLR.LineOut
	case "Chat Mic":
		return GOXLR.ChatMic
	}

	panic(fmt.Sprintf("Output %s can't be converted to GOXLR format", output))
}
