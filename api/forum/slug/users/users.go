package users

import (
	"log"

	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/user"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/database"
)

// Пользователи данного форума

func Users(ctx *fasthttp.RequestCtx) {
	usrArr := model_pool.UserArrPool.Get().(user.Arr)
	defer model_pool.UserArrPool.Put(usrArr)

	err := usrArr.GetForumUsers(ctx.UserValue("slug").(string), ctx.QueryArgs().Peek("limit"),
		ctx.QueryArgs().Peek("since"), ctx.QueryArgs().Peek("desc"))

	var resp []byte
	switch err {
	case nil:
		resp, err = usrArr.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}

	case database.ErrorForumNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	default:
		log.Fatalln(err)
	}

	usrArr = usrArr[:0]
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
