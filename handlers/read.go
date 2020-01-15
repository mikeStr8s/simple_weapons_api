package handlers

import (
	"encoding/json"
	"path"
	"sort"

	"github.com/mikeStr8s/simple_weapons_api/util"
	"github.com/valyala/fasthttp"
)

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

// Lookup is a dynamic handler for lookup datasets
func Lookup(ctx *fasthttp.RequestCtx) {
	util.SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(GetData(path.Base(string(ctx.Path())))); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// GetData gets the data in the file represented by the supplied dataset
func GetData(dataset string) interface{} {
	fileData := util.ReadJSONFile(dataset) // Raw byte array of specified dataset from file
	if dataset == util.MONSTER {
		return ParseMonster(fileData)
	} else if lookupDataset, ok := util.RELATED[dataset]; ok {
		return ParseRelational(fileData, lookupDataset)
	}
	return ParseLookup(fileData)
}

// ParseLookup takes a []byte from a JSON file and parses the
// data into an array of LookupRecords
func ParseLookup(byteData []byte) []LookupRecord {
	var lookupData []LookupRecord         // Instantiate empty array of LookupRecords
	json.Unmarshal(byteData, &lookupData) // Parse JSON data into empty array
	return lookupData
}

// ParseRelational takes []byte from a JSON file and its related
// lookupDataset to parse into a complete array of RelationalRecords
func ParseRelational(byteData []byte, lookupDataset string) []RelationalRecord {
	var rawRelationalData []RawRelationalRecord  // Instantiate empty array of RawRelationalRecords
	json.Unmarshal(byteData, &rawRelationalData) // Parse JSON data into empty array

	var relationalData []RelationalRecord                       // Instantiate empty final array of RelationalRecords
	lookupData := ParseLookup(util.ReadJSONFile(lookupDataset)) // Get array of LookupRecords for related dataset
	for _, rawRelationalRecord := range rawRelationalData {
		var lookupRecord LookupRecord
		for _, lr := range lookupData {
			if lr.ID == rawRelationalRecord.Related {
				lookupRecord = lr // Identify and save related LookupRecord
			}
		}

		// Create final RelationalRecord and add to final array
		var relationalRecord RelationalRecord
		relationalRecord.ID = rawRelationalRecord.ID
		relationalRecord.Value = rawRelationalRecord.Value
		relationalRecord.Related = lookupRecord
		relationalData = append(relationalData, relationalRecord)
	}
	return relationalData
}

// ParseMonster takes []byte from a JSON file and parses out
// monsters and all related data into a monster list.
func ParseMonster(byteData []byte) []MonsterRecord {
	var rawMonsterData []RawMonsterRecord
	json.Unmarshal(byteData, &rawMonsterData)

	var monsterData []MonsterRecord
	for _, rawMonsterRecord := range rawMonsterData {
		monsterRecord := MonsterRecord{
			rawMonsterRecord.ID,
			rawMonsterRecord.Name,
			rawMonsterRecord.HitPoints,
			rawMonsterRecord.ArmorClass,
			rawMonsterRecord.STR,
			rawMonsterRecord.DEX,
			rawMonsterRecord.CON,
			rawMonsterRecord.INT,
			rawMonsterRecord.WIS,
			rawMonsterRecord.CHA,
			rawMonsterRecord.Challenge,
			rawMonsterRecord.Traits,
			rawMonsterRecord.Actions,
			rawMonsterRecord.LegendaryActions,
			rawMonsterRecord.Reactions,
			map[string][]RelationalRecord{},
		}

		for relationalDataset, ids := range rawMonsterRecord.Related {
			relationalData := ParseRelational(util.ReadJSONFile(relationalDataset), util.RELATED[relationalDataset])
			rIDs := ids
			sort.Slice(rIDs, func(i int, j int) bool { return rIDs[i] < rIDs[j] })
			for _, relationalRecord := range relationalData {
				if len(rIDs) > 0 && rIDs[0] == relationalRecord.ID {
					monsterRecord.Related[relationalDataset] = append(monsterRecord.Related[relationalDataset], relationalRecord)
					_, rIDs = rIDs[0], rIDs[1:]
				}
			}
		}
		monsterData = append(monsterData, monsterRecord)
	}
	return monsterData
}
