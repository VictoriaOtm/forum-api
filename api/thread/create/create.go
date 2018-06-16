package create

import (
	"strings"
	"sync"

	"time"

	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/poststore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	"github.com/VictoriaOtm/forum-api/helpers"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/valyala/fasthttp"
)

var (
	errorThreadNotFound      = e.MakeError("error: thread not found")
	errorParentInWrongThread = e.MakeError("error: parent in wrong thread")
	errorUserNotFount        = e.MakeError("error: user not found")
)

var timeCreated = time.Unix(time.Now().Unix(), 0)

func Create(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	thr := threadstore.Pool.Acquire()
	defer threadstore.Pool.Utilize(thr)

	var err error
	sOrID := ctx.UserValue("slug_or_id").(string)
	if helpers.IsNumber(sOrID) {
		err = thr.GetById(sOrID)
	} else {
		err = thr.GetBySlug(sOrID)
	}

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(errorThreadNotFound)
		return
	}

	ps := poststore.PoolPostSlice.Acquire()
	defer poststore.PoolPostSlice.Utilize(ps)
	ps.MustUnmarshalJSON(ctx.PostBody())

	wg := sync.WaitGroup{}

	wg.Add(1)
	var err2 error
	go func() {
		err2 = ps.ValidateParentsAndThread(thr.Id)
		wg.Done()
	}()

	usersMap := treemap.NewWithStringComparator()

	wg.Add(1)
	var err3 error
	go func() {
		for _, post := range ps {
			usersMap.Put(strings.ToLower(post.Author), nil)
		}

		err3 = userstore.GetByNicknames(usersMap)
		wg.Done()
	}()
	wg.Wait()

	if err2 != nil {
		ctx.SetStatusCode(409)
		ctx.Write(errorParentInWrongThread)
		return
	}

	if err3 != nil {
		ctx.SetStatusCode(404)
		ctx.Write(errorUserNotFount)
		return
	}

	wg.Add(1)
	go func() {
		for i, post := range ps {
			u, _ := usersMap.Get(strings.ToLower(post.Author))
			ps[i].Author = u.(userstore.User).Nickname

			ps[i].Thread = thr.Id
			ps[i].Forum = thr.Forum
			ps[i].Created = timeCreated
		}

		ps.Insert()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		userstore.StoreInForumUserTable(thr.Forum, usersMap.Values())
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		forumstore.UpdateForumPosts(thr.Forum, int64(len(ps)))
		wg.Done()
	}()

	wg.Wait()
	ctx.SetStatusCode(201)
	ctx.Write(ps.MustMarshalJSON())
}
