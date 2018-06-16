package get_details

import (
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/helpers"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorThreadNotExists = e.MakeError("error: thread not exists")

// Получение информации о ветке
func Details(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	thr := threadstore.Pool.Acquire()
	defer threadstore.Pool.Utilize(thr)

	var err error
	slugOrId := ctx.UserValue("slug_or_id").(string)
	if helpers.IsNumber(slugOrId) {
		err = thr.GetById(slugOrId)
	} else {
		err = thr.GetBySlug(slugOrId)
	}

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorThreadNotExists)
		return
	}

	ctx.Write(thr.MustMarshalJSON())
}
