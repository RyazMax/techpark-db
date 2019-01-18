package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

// Get method returns information about user
func UserGetProfile(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user := models.User{}
	exist := user.GetUserByNick(nickname, db)

	if exist {
		serveJson(ctx, http.StatusOK, &user)
	} else {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not found user"})
	}
}

// Post method returns updates information about user
func UserUpdate(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	body := ctx.PostBody()

	newUser := models.User{}
	json.Unmarshal(body, &newUser)

	newUser.Nickname = nickname
	sameUsers := newUser.GetLike(db)

	if len(sameUsers) == 0 || (len(sameUsers) == 1 && sameUsers[0].Nickname != newUser.Nickname) {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "User not found"})
	} else if len(sameUsers) == 1 {
		newUser.Update(db)
		newUser.GetUserByNick(nickname, db)
		serveJson(ctx, http.StatusOK, &newUser)
	} else {
		serveJson(ctx, http.StatusConflict, &models.Message{Message: "Can not update user"})
	}
}
