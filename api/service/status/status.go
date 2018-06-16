package status

import (
	"log"

	"github.com/VictoriaOtm/forum-api/helpers/status"
	"github.com/valyala/fasthttp"
)

// Получение информации о базе данных

func Status(ctx *fasthttp.RequestCtx) {
	s := status.Status{}
	err := s.Get()

	if err != nil {
		log.Fatalln(err)
	}

	resp, _ := s.MarshalJSON()
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
