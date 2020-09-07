package app

import (
	"context"
	"errors"
	"github.com/rylio/ytdl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type ytdlRepo struct {
	mock.Mock
}

func (m *ytdlRepo) GetDownloadURL(ctx context.Context, info *ytdl.VideoInfo, format *ytdl.Format) (*url.URL, error) {
	args := m.Called(ctx, info, format)
	return args.Get(0).(*url.URL), args.Error(1)
}

func (m *ytdlRepo) GetVideoInfo(ctx context.Context, value interface{}) (*ytdl.VideoInfo, error) {
	args := m.Called(ctx, value)
	return args.Get(0).(*ytdl.VideoInfo), args.Error(1)
}

func TestVideo_GetFileURL(t *testing.T) {
	ctx := context.Background()
	// video-JZAunPKoHL0
	id1 := "JZAunPKoHL0"
	u1, _ := url.Parse("https://www.youtube.com/watch?v=" + id1)
	f1 := &ytdl.Format{
		Itag: ytdl.Itag{
			Number: 32,
		},
	}
	d1 := &ytdl.VideoInfo{
		ID: id1,
		Formats: ytdl.FormatList{
			f1,
		},
	}

	// video-JZAunPKoHL0 no format
	d2 := &ytdl.VideoInfo{
		ID:      id1,
		Formats: nil,
	}
	ytdlR := new(ytdlRepo)
	ytdlS := NewYTDLService(ytdlR)

	type arguments struct {
		videoInfo *ytdl.VideoInfo
	}
	tests := []struct {
		name      string
		arguments arguments
		want      url.URL
		wantErr   bool
		err       error
		before    func()
	}{
		{
			name: "should throw error on url not found",
			arguments: arguments{
				videoInfo: d2,
			},
			want:    url.URL{},
			wantErr: true,
			err:     errors.New(FormatsNotFound),
			before: func() {
				ytdlR.On("GetDownloadURL", ctx, d1, f1).Return(u1, nil)
			},
		},
		{
			name: "should return url http://google.com",
			arguments: arguments{
				videoInfo: d1,
			},
			want:    *u1,
			wantErr: false,
			err:     nil,
			before: func() {
				ytdlR.On("GetDownloadURL", ctx, d1, f1).Return(u1, nil)
			},
		},
		//	TODO: test with fixed format and get real url
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			got, err := ytdlS.GetFileURL(tt.arguments.videoInfo)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("GetFileURL() error = %v, want %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_GetFileInformation_ShouldReturnErrorOnVideoInfoCall(t *testing.T) {
	ctx := context.Background()
	u := "https://www.youtube.com/watch?v=123"
	fileInformation := &ytdl.VideoInfo{}
	ytdlR := new(ytdlRepo)
	ytdlS := NewYTDLService(ytdlR)
	errExpected := errors.New("weird error")
	ytdlR.On("GetVideoInfo", ctx, u).Return(fileInformation, errExpected)

	_, err := ytdlS.GetFileInformation(u)
	if err != errExpected {
		t.Errorf("GetFileInformation() error = %v, wantErr %v", err, errExpected)
		return
	}
}

func TestVideo_GetFileInformation_ShouldVideoWithInformation(t *testing.T) {
	ctx := context.Background()
	id := "JZAunPKoHL0"
	u, _ := url.Parse("https://www.youtube.com/watch?v=" + id)
	n := time.Now()
	format := &ytdl.Format{
		Itag: ytdl.Itag{
			Number: 32,
		},
	}
	details := &ytdl.VideoInfo{
		ID:          id,
		Title:       "Wojciech Szczęsny's Most Incredible Saves! | The Best of Tek | Juventus",
		Description: "In honour of Wojciech Szczęsny upcoming 30th birthday, we put together all of his best saves at Juventus so far! ",
		Formats: ytdl.FormatList{
			format,
		},
		Duration:      time.Hour / 2,
		Uploader:      "Juventus",
		DatePublished: n,
		Keywords:      []string{"1", "2"},
	}
	ytdlR := new(ytdlRepo)
	ytdlR.On("GetDownloadURL", ctx, details, format).Return(u, nil)
	ytdlR.On("GetVideoInfo", ctx, id).Return(details, nil)

	ytdlS := NewYTDLService(ytdlR)

	result, err := ytdlS.GetFileInformation(id)
	if err != nil {
		t.Errorf("GetFileInformation() error = %v, wantErr %v", err, nil)
		return
	}

	assert.Equal(t, result.ID, id)
	assert.Equal(t, result.FileUrl, u)
	assert.Equal(t, result.Title, details.Title)
	assert.Equal(t, result.Description, details.Description)
	assert.Equal(t, result.Duration, details.Duration)
	assert.Equal(t, result.Author, details.Uploader)
	assert.Equal(t, result.DatePublished, details.DatePublished)
	assert.Equal(t, result.Keywords, details.Keywords)
	assert.Equal(t, result.Keywords, details.Keywords)
	assert.Equal(t, result.ContentType, "")
	assert.Equal(t, result.ContentLength, int64(0))
}

func TestVideo_GetFileInformation_ShouldReturnErrorOnGetFileUrlError(t *testing.T) {
	ctx := context.Background()
	id := "JZAunPKoHL0"
	u, _ := url.Parse("https://www.youtube.com/watch?v=" + id)
	n := time.Now()
	format := &ytdl.Format{
		Itag: ytdl.Itag{
			Number: 32,
		},
	}
	details := &ytdl.VideoInfo{
		ID:          id,
		Title:       "Wojciech Szczęsny's Most Incredible Saves! | The Best of Tek | Juventus",
		Description: "In honour of Wojciech Szczęsny upcoming 30th birthday, we put together all of his best saves at Juventus so far! ",
		Formats: ytdl.FormatList{
			format,
		},
		Duration:      time.Hour / 2,
		Uploader:      "Juventus",
		DatePublished: n,
		Keywords:      []string{"1", "2"},
	}
	errWanted := errors.New("can't get file url for this video")
	ytdlR := new(ytdlRepo)
	ytdlR.On("GetDownloadURL", ctx, details, format).Return(u, errWanted)
	ytdlR.On("GetVideoInfo", ctx, id).Return(details, nil)

	ytdlS := NewYTDLService(ytdlR)
	_, err := ytdlS.GetFileInformation(id)
	if err != errWanted {
		t.Errorf("GetFileInformation() error = %v, wantErr %v", err, nil)
		return
	}
}
