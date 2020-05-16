package api

import (
	"github.com/gofiber/fiber"
	"html/template"
	"ygp/pkg/youtube"
)

const (
	BaseFeedURL = "https://yt.psmarcin.dev/feed/channel/"
)

var (
	templates = template.Must(template.ParseGlob("./templates/*.tmpl"))
)

func rootHandler(ctx *fiber.Ctx) {
	var channels []youtube.Channel
	var err error

	ctx.Set("content-type", "text/html; charset=utf-8")

	channelId := ctx.FormValue("channelId")
	if channelId != "" {
		channelId = BaseFeedURL + channelId
	}

	q := ctx.FormValue("q")
	if q != "" {
		channels, err = youtube.Yt.ChannelsListFromCache(q)
	}
	if err != nil {
		ctx.Next(err)
		return
	}

	err = templates.ExecuteTemplate(ctx.Fasthttp.Response.BodyWriter(), "index.tmpl", map[string]interface{}{
		"Channels":  channels,
		"ChannelId": channelId,
	})

	if err != nil {
		l.WithError(err).Errorf("erron while rendering template")
		ctx.Next(err)
	}
}
