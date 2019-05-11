<hr>
<h3 align="center">ygp - YouTube Goes Podcast</h3>
<hr>

🎞 👉🎙 Put youtube channel get podcast audio feed 🎞 👉🎙

This API is created mainly to receive youtube channel and return audio podcast feed that you can use in your favorite podcast app.

## Features/Roadmap
* [x] Generate podcast feed for youtube channel
* [x] Filter only wanted videos
* [x] Use audio file for videos
* [x] Get trending channels (base on trending videos)

### Examples
Use this url `https://ygp.psmarcin.dev/feed/channel/UCblfuW_4rakIf2h6aqANefA` in your favorite podcast app. It works on desktop and mobile too. Tested on:
* iPhone Podcast App
* iTunes MacOS App
* Plex Web App
* Plex iOS App

### Screens
![Tested apps](assets/iphone-podcast-app.png "Tested apps")

## Development

### Requirements
1. Go in version `>=1.9`, more: https://golang.org/dl/
1. Docker, more: https://docs.docker.com/install/
1. Realize (auto restart), more https://github.com/oxequa/realize

### Build
1. `make build`

### Test
1. `make test`

### Development run
1. `make dev`

## Credits
This project uses big part of https://github.com/rylio/ytdl. I couldn't use it as dependencies because there was conflict with `logrus`. Will use it as dependency as soon as it will fix that problem.
