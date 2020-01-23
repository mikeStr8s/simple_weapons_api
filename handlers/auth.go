package handlers

import (
	"encoding/json"

	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/mikeStr8s/simple_weapons_api/util"
	"github.com/valyala/fasthttp"
)

// Register is the registration handler to create an account
// to allow the user to traverse the API through the auth middleware
func Register(ctx *fasthttp.RequestCtx) {
	// If there is unique username, create the new user and return the login token
	// Else return 409 because account with username already exists
	if userData, exists := userExists(ctx.PostBody()); !exists {
		newUser := createUser(userData)
		util.SetResponse(ctx)
		if err := json.NewEncoder(ctx).Encode(newUser); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
	} else {
		ctx.Response.SetStatusCode(409)
	}
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
