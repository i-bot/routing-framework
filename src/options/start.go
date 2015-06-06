package options

import (
	"framework"
	"settings"
)

type Start struct {
	properties *settings.Settings
}

func (start Start) Applies(identifier string) bool {
	return identifier == "-s" || identifier == "--start"
}

func (start Start) Execute(args []string) {
	loadSettings(args[0])

	framework.Start(properties)
}
