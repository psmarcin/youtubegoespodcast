package adapters

import "github.com/rylio/ytdl"

func NewYTDLRepository() *ytdl.Client {
	return ytdl.DefaultClient
}
