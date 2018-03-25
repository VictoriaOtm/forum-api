package post_profile

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/user"
	"log"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/database"
)

// Изменение данных о пользователе

func Profile(ctx *fasthttp.RequestCtx) {
	usr := model_pool.UserPool.Get().(*user.User)
	defer model_pool.UserPool.Put(usr)

	usr.Nickname = ctx.UserValue("nickname").(string)

	usrNewInfo := &user.Update{}
	err := usrNewInfo.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}

	err = usr.Update(usrNewInfo)

	var resp []byte

	switch err {
	case nil:
		resp, err = usr.MarshalJSON()

	case database.ErrorUserNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	case database.ErrorUserConflict:
		resp = error_m.CommonError
		ctx.SetStatusCode(409)

	default:
		log.Fatalln(err)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
