package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	currentDir, _ = os.Getwd()
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

// Index is the base landing for routing
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// SetResponse takes a context pointer and assignes the response context data for the API
func SetResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(200)
}

// Lookup is a dynamic handler for lookup datasets
func Lookup(ctx *fasthttp.RequestCtx) {
	SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(GetData(path.Base(string(ctx.Path())))); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// Create is a dynamic handler for POST requests
func Create(ctx *fasthttp.RequestCtx) {
	SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(PostData(path.Base(string(ctx.Path())), ctx.PostBody())); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// ReadJSONFile reads the contents of a JSON file and returns an
// array of bytes of file data to be parsed into data object
func ReadJSONFile(dataset string) []byte {
	file, err := os.Open(path.Join(currentDir, "data/"+dataset+".json"))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	return bytes
}

// GetData gets the data in the file represented by the supplied dataset
func GetData(dataset string) interface{} {
	fileData := ReadJSONFile(dataset) // Raw byte array of specified dataset from file
	if dataset == MONSTER {
		return ParseMonster(fileData)
	} else if lookupDataset, ok := RELATED[dataset]; ok {
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

	var relationalData []RelationalRecord                  // Instantiate empty final array of RelationalRecords
	lookupData := ParseLookup(ReadJSONFile(lookupDataset)) // Get array of LookupRecords for related dataset
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
			relationalData := ParseRelational(ReadJSONFile(relationalDataset), RELATED[relationalDataset])
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

// Contains checks a list of strings to see if the supplied
// term exists in the array. If it exists the index of the
// term is returned, otherwise -1
func Contains(list []string, term string) int {
	for idx, item := range list {
		if item == term {
			return idx
		}
	}
	return -1
}

func WriteJSONData(dataset string, byteData []byte) {
	ioutil.WriteFile("data/"+dataset+".json", byteData, 0644)
}

func PostData(dataset string, byteData []byte) interface{} {
	if dataset == MONSTER {

	} else if lookupDataset, ok := RELATED[dataset]; ok {
		return CreateRelational(dataset, lookupDataset, byteData)
	}
	println(dataset)
	return byteData
}

func CreateRelational(dataset string, lookupDataset string, byteData []byte) RelationalRecord {
	var rawRelationalRecord RawRelationalRecord
	json.Unmarshal(byteData, &rawRelationalRecord)

	lookupData := ParseLookup(ReadJSONFile(lookupDataset))
	for _, lookupRecord := range lookupData {
		if lookupRecord.ID == rawRelationalRecord.Related {
			relationalData := ParseRelational(ReadJSONFile(dataset), lookupDataset)
			sort.Slice(relationalData, func(i, j int) bool { return relationalData[i].ID < relationalData[j].ID })
			relationalRecord := RelationalRecord{relationalData[len(relationalData)-1].ID + 1, rawRelationalRecord.Value, lookupRecord}
			relationalData = append(relationalData, relationalRecord)
			relationalJSON, _ := json.Marshal(relationalData)
			WriteJSONData(dataset, relationalJSON)
			return relationalRecord
		}
	}
	return RelationalRecord{999999, "Problem encountered, no data was Added", LookupRecord{10, "Related ID does not reference existing entry"}}
}

func main() {
	router := router.New()
	router.GET("/", Index)

	api := router.Group("/api")
	api.GET("/abilityscore", Lookup)
	api.GET("/condition", Lookup)
	api.GET("/damage", Lookup)
	api.GET("/language", Lookup)
	api.GET("/movement", Lookup)
	api.GET("/sense", Lookup)
	api.GET("/skill", Lookup)
	api.GET("/movementspeed", Lookup)
	api.GET("/savingthrow", Lookup)
	api.GET("/sensevalue", Lookup)
	api.GET("/skillvalue", Lookup)
	api.GET("/monster", Lookup)

	api.POST("/movementspeed", Create)
	api.POST("/skillvalue", Create)
	api.POST("/savingthrow", Create)
	api.POST("/sensevalue", Create)
	api.POST("/monster", Create)

	log.Fatal(fasthttp.ListenAndServe(":1234", router.Handler))
}
