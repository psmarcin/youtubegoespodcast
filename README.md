# YgP - API
It's simple API helper for yt.psmarcin.dev.

## Features
* [x] Return trending channels

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
