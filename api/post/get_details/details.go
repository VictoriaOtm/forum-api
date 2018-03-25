package get_details

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/post"
	"log"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"strings"
)

// Получение информации о ветке обсуждения

func Details(ctx *fasthttp.RequestCtx) {
	p := model_pool.PostDetailsPool.Get().(*post.PostDetails)
	defer func() {
		p.User = nil
		p.Forum = nil
		p.Thread = nil

		model_pool.PostDetailsPool.Put(p)
	}()

	related := ctx.QueryArgs().Peek("related")
	relatedArr := strings.Split(string(related), ",")

	err := p.Get(ctx.UserValue("id").(string), relatedArr)

	var resp []byte
	switch err {
	case nil:
		resp, err = p.MarshalJSON()
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
