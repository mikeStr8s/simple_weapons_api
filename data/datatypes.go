package data

// LookupRecord is a struct that represents the common shape of the JSON data used for simple lookup datasets.
type LookupRecord struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RawRelationalRecord is a struct that represents the common shape of the JSON data for relational datasets.
type RawRelationalRecord struct {
	ID      int    `json:"id"`
	Value   string `json:"value"`
	Related int    `json:"related"`
}

// RelationalRecord is a struct that represents the common shape of the JSON data for relational datasets
// with a parsed lookup related reference
type RelationalRecord struct {
	ID      int          `json:"id"`
	Value   string       `json:"value"`
	Related LookupRecord `json:"related"`
}

// RawMonsterRecord is a struct representing the data structure of Monsters
type RawMonsterRecord struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	HitPoints        string           `json:"hitpoints"`
	ArmorClass       string           `json:"armorclass"`
	STR              int              `json:"STR"`
	DEX              int              `json:"DEX"`
	CON              int              `json:"CON"`
	INT              int              `json:"INT"`
	WIS              int              `json:"WIS"`
	CHA              int              `json:"CHA"`
	Challenge        int              `json:"challenge"`
	Traits           string           `json:"traits"`
	Actions          string           `json:"actions"`
	LegendaryActions string           `json:"legendaryactions"`
	Reactions        string           `json:"reactions"`
	Related          map[string][]int `json:"related"`
}

// MonsterRecord is a struct representing the data structure of Monsters
type MonsterRecord struct {
	ID               int                           `json:"id"`
	Name             string                        `json:"name"`
	HitPoints        string                        `json:"hitpoints"`
	ArmorClass       string                        `json:"armorclass"`
	STR              int                           `json:"STR"`
	DEX              int                           `json:"DEX"`
	CON              int                           `json:"CON"`
	INT              int                           `json:"INT"`
	WIS              int                           `json:"WIS"`
	CHA              int                           `json:"CHA"`
	Challenge        int                           `json:"challenge"`
	Traits           string                        `json:"traits"`
	Actions          string                        `json:"actions"`
	LegendaryActions string                        `json:"legendaryactions"`
	Reactions        string                        `json:"reactions"`
	Related          map[string][]RelationalRecord `json:"related"`
}
