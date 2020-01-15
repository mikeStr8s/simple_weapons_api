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

func PostData(dataset string, byteData []byte) interface{} {
	if dataset == util.MONSTER {

	} else if lookupDataset, ok := util.RELATED[dataset]; ok {
		return CreateRelational(dataset, lookupDataset, byteData)
	}
	println(dataset)
	return byteData
}

func WriteJSONData(dataset string, byteData []byte) {
	ioutil.WriteFile(dataset+".json", byteData, 0644)
}

func CreateRelational(dataset string, lookupDataset string, byteData []byte) RelationalRecord {
	var rawRelationalRecord RawRelationalRecord
	json.Unmarshal(byteData, &rawRelationalRecord)

	lookupData := ParseLookup(util.ReadJSONFile(lookupDataset))
	for _, lookupRecord := range lookupData {
		if lookupRecord.ID == rawRelationalRecord.Related {
			relationalData := ParseRelational(util.ReadJSONFile(dataset), lookupDataset)
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
