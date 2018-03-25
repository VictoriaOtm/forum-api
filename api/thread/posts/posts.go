package posts

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/post"
	"bytes"
	"log"
	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/models/error_m"
)

// Сообщение данной ветви обсуждения

func Posts(ctx *fasthttp.RequestCtx) {
	postsArr := model_pool.PostArrPool.Get().(post.PostsArr)
	defer func() {
		postsArr = postsArr[:0]
		model_pool.PostArrPool.Put(postsArr)
	}()

	err := postsArr.Get(
		ctx.UserValue("slug_or_id").(string),
		string(ctx.QueryArgs().Peek("sort")),
		ctx.QueryArgs().Peek("limit"),
		ctx.QueryArgs().Peek("since"),
		bytes.Equal([]byte("true"), ctx.QueryArgs().Peek("desc")),
	)

	var resp []byte
	switch err {
	case nil:
		resp, err = postsArr.MarshalJSON()

	case database.ErrorThreadNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	default:
		log.Fatalln(err)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
