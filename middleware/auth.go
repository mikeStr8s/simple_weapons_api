package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/mikeStr8s/simple_weapons_api/util"
	"github.com/valyala/fasthttp"
)

// Auth is a handler middleware that checks that the user has an auth token
// in their header. If so it will then check to see that the auth
// is valid before sending the user to the desired end handler
func Auth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		authToken := ctx.Request.Header.Peek("Simple-Weapons")
		if authToken != nil {
			username, password, hasAuth := checkTokenValid(string(authToken))
			if hasAuth {
				err := scrypt.CompareHashAndPassword([]byte(getPasswordHash(username)), []byte(password))
				if err != nil {
					ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
				} else {
					handler(ctx)
					return
				}
			}
		}
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	})
}

func checkTokenValid(token string) (username, password string, ok bool) {
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}
	userpass := string(tokenBytes)
	sepIdx := strings.IndexByte(userpass, ':')
	if sepIdx < 0 {
		return
	}

	userByteData := util.ReadJSONFile("user")
	var userData map[string]string
	json.Unmarshal(userByteData, &userData)

	if _, ok := userData[userpass[:sepIdx]]; ok {
		return userpass[:sepIdx], userpass[sepIdx+1:], true
	}

	return
}

func getPasswordHash(username string) string {
	userByteData := util.ReadJSONFile("user")
	var userData map[string]string
	json.Unmarshal(userByteData, &userData)

	return userData[username]
}
