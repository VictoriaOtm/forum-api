package posts

import (
	"bytes"

	"strconv"

	"strings"

	"github.com/VictoriaOtm/forum-api/database/stores/poststore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/helpers"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/VictoriaOtm/forum-api/helpers/unsafe_map"
	"github.com/valyala/fasthttp"
)

var responseErrorThreadNotExists = e.MakeError("error: thread not exists")

func getId(slugOrID string) int32 {
	if helpers.IsNumber(slugOrID) {
		id, _ := strconv.Atoi(slugOrID)
		return int32(id)
	} else {
		return unsafe_map.SlugIDMap.GetId(strings.ToLower(slugOrID))
	}
}

// Сообщение данной ветви обсуждения
func Posts(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")
	slugOrId := ctx.UserValue("slug_or_id").(string)
	var err error

	id := getId(slugOrId)
	if id == 0 {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorThreadNotExists)
		return
	}

	ps := poststore.PoolPostSlice.Acquire()
	defer poststore.PoolPostSlice.Utilize(ps)

	err = ps.Get(
		id,
		ctx.QueryArgs().Peek("sort"),
		ctx.QueryArgs().Peek("since"),
		ctx.QueryArgs().Peek("limit"),
		bytes.Equal([]byte("true"), ctx.QueryArgs().Peek("desc")),
	)

	if err != nil {
		panic(err)
	}

	if len(ps) == 0 {
		thr := threadstore.Pool.Acquire()
		defer threadstore.Pool.Utilize(thr)

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
	}

	ctx.Write(ps.MustMarshalJSON())
}
