<hr>
<h3 align="center">ygp - YouTube Goes Podcast</h3>
<hr>

🎞 👉🎙 Put youtube channel get podcast audio feed 🎞 👉🎙
This API is created mainly to receive youtube channel and return audio podcast feed that you can use in your favorite podcast app.

## Features/Roadmap
* [x] Return trending channels
* [x] Redirect to audio stream (if available)
* [x] Get trending channels (base on trending videos)
* [x] Get channel feed
* [x] Test channel feed

## Development

### Requirements
1. Go >1.9
1. `now` - for deployment
1. Docker
1. Realize https://github.com/oxequa/realize

### Install dependencies
1. `make dependencies`

### Build
1. `make build`

### Test
1. `make test`

### Development run
1. `make dev`

## Credits
This project uses big part of https://github.com/rylio/ytdl. I couldn't use it as dependencies because there was conflict with `logrus`. Will use it as dependency as soon as it will fix that problem.
