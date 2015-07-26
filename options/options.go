package options

import (
	"settings"
	"strings"
)

type IOption interface {
	Applies(string) bool
	Execute([]string)
}

var (
	options    []IOption
	properties *settings.Settings
)

func Interpret(args []string) {
	initialize()

	identifiers, arguments := split(args)

	for _, option := range options {
		for _, identifier := range identifiers {
			if option.Applies(identifier) {
				option.Execute(arguments)
			}
		}
	}
}

func initialize() {
	options = append(options, &Start{})
}

func loadSettings(location string) {
	if properties == nil {
		properties = settings.LoadSettings(location)
	}
}

func split(args []string) (identifiers []string, arguments []string) {
	var splitIndex int

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			splitIndex = i + 1
			break
		}
	}

	return args[:splitIndex], args[splitIndex:]
}
