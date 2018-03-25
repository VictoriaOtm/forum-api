package create

import (
	"log"

	"github.com/VictoriaOtm/forum-api/models/user"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/database"
)

// Создание нового пользователя
func Create(ctx *fasthttp.RequestCtx) {
	var resp []byte

	usr := model_pool.UserPool.Get().(*user.User)
	defer model_pool.UserPool.Put(usr)

	err := usr.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}
	usr.Nickname = ctx.UserValue("nickname").(string)

	err, userDuplicates := usr.Create()

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		resp, err = usr.MarshalJSON()
	case database.ErrorUserConflict:
		ctx.SetStatusCode(409)
		resp, err = userDuplicates.MarshalJSON()
	default:
		log.Fatalln(err)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
