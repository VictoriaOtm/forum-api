package create

import (
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	"github.com/valyala/fasthttp"
)

// Создание нового пользователя
func Create(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	usr := userstore.Pool.Acquire()
	defer userstore.Pool.Utilize(usr)
	usr.MustUnmarshalJSON(ctx.PostBody())
	usr.Nickname = ctx.UserValue("nickname").(string)

	us, err := usr.GetByNicknameOrEmail(usr.Nickname, usr.Email)
	if len(us) != 0 {
		ctx.SetStatusCode(409)
		ctx.Write(us.MustMarshalJSON())
		return
	}

	err = usr.Insert()

	if err != nil {
		panic(err)
	}

	ctx.SetStatusCode(201)
	ctx.Write(usr.MustMarshalJSON())
}
