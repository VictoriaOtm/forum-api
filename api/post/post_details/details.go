package post_details

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/models/post"
	"github.com/VictoriaOtm/forum-api/models/error_m"
)

// Изменение сообщения

func Details(ctx *fasthttp.RequestCtx) {
	pUpd := post.PostUpdate{}
	pUpd.UnmarshalJSON(ctx.Request.Body())

	p := post.Post{}
	err := p.Update(pUpd, ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(error_m.CommonError)
		ctx.Response.Header.Set("Content-type", "application/json")
		return
	}

	resp, _ := p.MarshalJSON()
	ctx.Write(resp)
	ctx.Response.Header.Set("Content-type", "application/json")
}
