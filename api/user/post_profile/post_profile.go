package post_profile

import (
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var (
	responseErrorUserNotExists = e.MakeError("error: user not exists")
	responseErrorUserConflict  = e.MakeError("error: user conflict")
)

// Изменение данных о пользователе
func Profile(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	usr := userstore.Pool.Acquire()
	defer userstore.Pool.Utilize(usr)

	err := usr.Get(ctx.UserValue("nickname"))
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorUserNotExists)
		return
	}

	usrUpdateInfo := userstore.UserUpdate{}
	usrUpdateInfo.MustUnmarshalJSON(ctx.PostBody())

	err = usr.Update(usrUpdateInfo)
	if err != nil {
		ctx.SetStatusCode(409)
		ctx.Write(responseErrorUserConflict)
		return
	}

	ctx.Write(usr.MustMarshalJSON())
}
