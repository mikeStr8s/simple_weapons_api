package handlers

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/mikeStr8s/simple_weapons_api/util"

	"github.com/valyala/fasthttp"
)

func basicAuth(ctx *fasthttp.RequestCtx) bool {
	auth := ctx.Request.Header.Peek("Authorization")
	if auth == nil {
		credentials := ctx.Request.Header.Peek("credentials")
		if credentials != nil {
			creds := strings.Split(string(credentials), ":")
			userAuth, exists := getUserAuth(creds[0])
			if exists {
				username, password, hasAuth := parseBasicAuth(userAuth)
				if creds[0] == username && creds[1] == password {
					return hasAuth
				}
			}
		}
		return false
	}
	username, _, hasAuth := parseBasicAuth(string(auth))
	_, exists := getUserAuth(username)
	if exists {
		return hasAuth
	}
	return false
}

func getUserAuth(username string) (auth string, ok bool) {
	byteData := util.ReadJSONFile("user")
	var userData map[string]string
	json.Unmarshal(byteData, &userData)

	if val, ok := userData[username]; ok {
		return val, ok
	}
	return
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func BasicAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		hasAuth := basicAuth(ctx)
		if hasAuth {
			h(ctx)
			return
		}
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	})
}
