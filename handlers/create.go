package handlers

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sort"

	"github.com/mikeStr8s/simple_weapons_api/util"
	"github.com/valyala/fasthttp"
)

// Create is a dynamic handler for POST requests
func Create(ctx *fasthttp.RequestCtx) {
	util.SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(PostData(path.Base(string(ctx.Path())), ctx.PostBody())); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// PostData determines what dataset is being posted and dispatches
// the byte data to the correct handler
func PostData(dataset string, byteData []byte) interface{} {
	if dataset == util.MONSTER {

	}
	lookupDataset, _ := util.RELATED[dataset]
	return CreateRelational(dataset, lookupDataset, byteData)
}

// WriteJSONData writes the byte data to the specified dataset
func WriteJSONData(dataset string, byteData []byte) {
	ioutil.WriteFile(path.Join(util.DIRNAME, "data", dataset+".json"), byteData, 0644)
}

// CreateRelational creates a relational data entry to be written to a file
func CreateRelational(dataset string, lookupDataset string, byteData []byte) RelationalRecord {
	var rawRelationalRecord RawRelationalRecord
	json.Unmarshal(byteData, &rawRelationalRecord)

	var rawRelationalData []RawRelationalRecord
	json.Unmarshal(util.ReadJSONFile(dataset), &rawRelationalData)

	sort.Slice(rawRelationalData, func(i, j int) bool { return rawRelationalData[i].ID < rawRelationalData[j].ID })
	rawRelationalRecord.ID = rawRelationalData[len(rawRelationalData)-1].ID + 1
	relationalJSON, _ := json.Marshal(append(rawRelationalData, rawRelationalRecord))
	WriteJSONData(dataset, relationalJSON)

	lookupData := ParseLookup(util.ReadJSONFile(lookupDataset))
	for _, lookupRecord := range lookupData {
		if lookupRecord.ID == rawRelationalRecord.Related {
			return RelationalRecord{rawRelationalRecord.ID, rawRelationalRecord.Value, lookupRecord}
		}
	}
	return RelationalRecord{999999, "Problem encountered, no data was Added", LookupRecord{10, "Related ID does not reference existing entry"}}
}
