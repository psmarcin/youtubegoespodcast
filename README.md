<hr>
<h2 align="center">YouTube Goes Podcast</h2>
<h4 align="center">🎞 👉🎙 Put youtube channel get podcast audio feed 🎞 👉🎙</h4>
<hr>

Youtube Goes Podcast makes podcast feed from YouTube channel.

1. open https://yt.psmarcin.dev 
1. type channel name
1. select from results
1. copy generated feed url
1. and use it in your favourite app

It's that simple!

## Features/Roadmap
* [x] **UI** for friendly usage
* [x] Find channel using **search** field
* [x] Automatically generate unique url for YouTube Channel
* [x] **Podcast app agnostics**. Works well in Apple Podcast, Plex Podcasts and others!
* [x] Live updates, you will get up-to-date list of the latest items immediately! 
* [x] Support more than latest 15 videos
* [x] Daily updates
* [x] Listen any video on your phone in background
* [x] **Works only for videos with "embed" enabled** 

### Examples
Use this url `https://yt.psmarcin.dev/feed/channel/UCblfuW_4rakIf2h6aqANefA` in your favorite podcast app. It works on desktop and mobile too. Tested on:
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
1. Modd (auto restart), more https://github.com/cortesi/modd
1. Google Cloud Console or mocks 
1. golangci-lint: https://golangci-lint.run/usage/install/

### Environment variables
Example environment variables
```bash
APP_ENV=development
GOOGLE_API_KEY=<YOUR_YOUTUBE_API_KEY>
PORT=8080
API_URL=http://localhost:8080/
```

### Build
1. `make build`

### Test
1. `make test`

### Develop
1. `make dev`

### Debug
1. `make debug`
