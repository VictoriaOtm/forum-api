package post_details

import (
	"github.com/VictoriaOtm/forum-api/database/stores/poststore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorPostNotFound = e.MakeError("error: not found")

// Изменение сообщения
func Details(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	p := poststore.Post{}
	err := p.Get(ctx.UserValue("id"))
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorPostNotFound)
		return
	}

	pUpdate := poststore.PostUpdate{}
	pUpdate.MustUnmarshalJSON(ctx.PostBody())

	err = p.Update(pUpdate)
	if err != nil {
		panic(err)
	}

	ctx.Write(p.MustMarshalJSON())
}
