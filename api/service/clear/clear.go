package clear

import (
	"github.com/VictoriaOtm/forum-api/database"
	"github.com/valyala/fasthttp"
)

// Очистка всех данных в базе

func Clear(ctx *fasthttp.RequestCtx) {
	_, err := database.DB.Exec("TRUNCATE t_forum, t_forum_user, t_posts, t_thread, t_user")
	if err != nil {
		panic(err)
	}

	_, err = database.DB.Exec("VACUUM FULL ANALYZE")
	if err != nil {
		panic(err)
	}
}
