package get_profile

import (
	"log"

	"github.com/VictoriaOtm/forum-api/models/user"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/models/error_m"
)

// Получение данных о пользователе

func Profile(ctx *fasthttp.RequestCtx) {
	usr := model_pool.UserPool.Get().(*user.User)
	defer model_pool.UserPool.Put(usr)

	usr.Nickname = ctx.UserValue("nickname").(string)

	err := usr.GetProfile()

	var resp []byte
	switch err {
	case nil:
		resp, err = usr.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}

	default:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
