package api

import (
	"github.com/gofiber/fiber"
	"html/template"
	"ygp/pkg/youtube"
)

const (
	BaseFeedURL = "https://yt.psmarcin.dev/feed/channel/"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("./templates/*.tmpl")
	if err != nil {
		l.WithError(err).Errorf("can't find templates")
	}
}

func rootHandler(ctx *fiber.Ctx) {
	var err error
	var channels []youtube.Channel

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
		l.WithError(err).Errorf("error while rendering template")
		ctx.Next(err)
	}
}
