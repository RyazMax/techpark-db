package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

// Post method createas new user
func UserCreate(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	body := ctx.PostBody()
	newUser := &models.User{}
	json.Unmarshal(body, newUser)

	newUser.Nickname = nickname
	sameUsers := models.GetUserByNickOrEmail(nickname, newUser.Email, db)
	if len(sameUsers) > 0 {
		serveJson(ctx, http.StatusConflict, &sameUsers)
		return
	}

	models.UserAdd(*newUser, db)
	serveJson(ctx, http.StatusCreated, newUser)
}
