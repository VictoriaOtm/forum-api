package status

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/models/status"
	"log"
)

// Получение информации о базе данных

func Status(ctx *fasthttp.RequestCtx) {
	s := status.Status{}
	err := s.Get()

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := s.MarshalJSON()
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
