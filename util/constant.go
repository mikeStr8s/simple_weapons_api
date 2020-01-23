package util

import "os"

var (
	// DIRNAME is the directory that the project is located in
	DIRNAME, _ = os.Getwd()
	// RELATED is a mapping of relational data to it's base lookup
	RELATED = map[string]string{
		"movementspeed": "movement",
		"savingthrow":   "abilityscore",
		"skillvalue":    "skill",
		"sensevalue":    "sense",
	}
)

const (
	// MONSTER is a constant with value Monster for global use to avoid naked strings
	MONSTER = "monster"
)
