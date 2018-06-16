package get_details

import (
	"bytes"
	"strconv"

	"sync"

	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorPostNotFound = e.MakeError("error: post not found")

// Получение информации о ветке обсуждения
func Details(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	pd := pdPool.Acquire()
	defer pdPool.Utilize(pd)

	err := pd.Post.Get(ctx.UserValue("id"))
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorPostNotFound)
		return
	}

	related := ctx.QueryArgs().Peek("related")
	relatedArr := bytes.Split(related, []byte(","))

	wg := sync.WaitGroup{}

	for _, r := range relatedArr {
		switch true {
		case bytes.Equal(r, []byte("user")):
			wg.Add(1)
			go func() {
				pd.User = userstore.Pool.Acquire()
				pd.User.Get(pd.Post.Author)
				wg.Done()
			}()

		case bytes.Equal(r, []byte("forum")):
			wg.Add(1)
			go func() {
				pd.Forum = forumstore.Pool.Acquire()
				pd.Forum.Get(pd.Post.Forum)
				wg.Done()
			}()

		case bytes.Equal(r, []byte("thread")):
			wg.Add(1)
			go func() {
				pd.Thread = threadstore.Pool.Acquire()
				pd.Thread.GetById(strconv.Itoa(int(pd.Post.Thread)))
				wg.Done()
			}()
		}
	}

	wg.Wait()
	ctx.Write(pd.MustMarshalJSON())
	if pd.Thread != nil {
		threadstore.Pool.Utilize(pd.Thread)
	}

	if pd.Forum != nil {
		forumstore.Pool.Utilize(pd.Forum)
	}

	if pd.User != nil {
		userstore.Pool.Utilize(pd.User)
	}
}
