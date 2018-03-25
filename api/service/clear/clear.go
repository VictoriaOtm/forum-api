package clear

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/database"
	"log"
)

// Очистка всех данных в базе

func Clear(ctx *fasthttp.RequestCtx) {
	_, err := database.DBConnPool.Exec("TRUNCATE t_forum, t_forum_user, t_posts, t_thread, t_user")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = database.DBConnPool.Exec("VACUUM FULL ANALYZE")
	if err != nil {
		log.Fatalln(err)
	}
}
