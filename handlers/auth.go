package handlers

import (
	"encoding/json"

	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/mikeStr8s/simple_weapons_api/util"
	"github.com/valyala/fasthttp"
)

func Register(ctx *fasthttp.RequestCtx) {
	util.SetResponse(ctx)

	tokenData := map[string]string{}
	if userData, exists := userExists(ctx.PostBody()); !exists {
		newUser := createUser(userData)
		tokenData["token"] = newUser["token"]
	}

	ctx.Response.Header.Set("simple-weapons", tokenData["token"]) // TODO: add auth token instead of ""
}

func userExists(byteData []byte) (postBody map[string]string, exists bool) {
	var postData map[string]string
	json.Unmarshal(byteData, &postData)

	userByteData := util.ReadJSONFile("user")
	var userData map[string]string
	json.Unmarshal(userByteData, &userData)

	if val, ok := userData[postData["username"]]; ok {
		return map[string]string{"token": val}, true
	}
	return postData, false
}

func createUser(postData map[string]string) map[string]string {
	username := postData["username"]
	password, _ := scrypt.GenerateFromPassword([]byte(postData["password"]), scrypt.DefaultParams)

	userByteData := util.ReadJSONFile("user")
	var userData map[string]string
	json.Unmarshal(userByteData, &userData)

	userData[username] = string(password)
	userJSON, _ := json.Marshal(userData)
	WriteJSONData("user", userJSON)

	return map[string]string{"token": string(password)}
}
