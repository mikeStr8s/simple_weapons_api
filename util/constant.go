package util


var (
	DIRNAME = "/root/go/src/github.com/mikeStr8s/simple_weapons_api"
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
