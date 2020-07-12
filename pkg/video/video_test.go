package video

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

type MockFileUrlGetter struct {
	mock.Mock
}

func (m *MockFileUrlGetter) GetDownloadURL(ctx context.Context, info *ytdl.VideoInfo, format *ytdl.Format) (*url.URL, error) {
	args := m.Called(ctx, info, format)
	return args.Get(0).(*url.URL), args.Error(1)
}

type MockFileInformationGetter struct {
	mock.Mock
}

func (m *MockFileInformationGetter) GetVideoInfo(ctx context.Context, value interface{}) (*ytdl.VideoInfo, error) {
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
	mFileUrlGetter := new(MockFileUrlGetter)
	mFileUrlGetter.On("GetDownloadURL", ctx, d1, f1).Return(u1, nil)

	type fields struct {
		ID            string
		FileUrl       *url.URL
		Description   string
		Title         string
		DatePublished time.Time
		Keywords      []string
		Author        string
		Duration      time.Duration
		ContentType   string
		ContentLength int64
		details       *ytdl.VideoInfo
		FileUrlGetter interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		argument FileUrlGetter
		want     url.URL
		wantErr  bool
		err      error
	}{
		{
			name:     "should throw error on no rawInformation provided",
			fields:   fields{},
			argument: mFileUrlGetter,
			want:     url.URL{},
			wantErr:  true,
			err:      errors.New(DetailsNotFound),
		},
		{
			name: "should throw error on url not found",
			fields: fields{
				details: d2,
			},
			argument: mFileUrlGetter,
			want:     url.URL{},
			wantErr:  true,
			err:      errors.New(FormatsNotFound),
		},
		{
			name: "should return url http://google.com",
			fields: fields{
				details: d1,
			},
			argument: mFileUrlGetter,
			want:     *u1,
			wantErr:  false,
			err:      nil,
		},
		//	TODO: test with fixed format and get real url
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Video{
				ID:             tt.fields.ID,
				FileUrl:        tt.fields.FileUrl,
				Description:    tt.fields.Description,
				Title:          tt.fields.Title,
				DatePublished:  tt.fields.DatePublished,
				Keywords:       tt.fields.Keywords,
				Author:         tt.fields.Author,
				Duration:       tt.fields.Duration,
				ContentType:    tt.fields.ContentType,
				ContentLength:  tt.fields.ContentLength,
				rawInformation: tt.fields.details,
			}
			got, err := v.GetFileURL(tt.argument)

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
	id := "123"
	v := New(id)
	fileInformation := &ytdl.VideoInfo{}
	fileInformationGetter := new(MockFileInformationGetter)
	errExpected := errors.New("weird error")
	fileInformationGetter.On("GetVideoInfo", ctx, v.URL).Return(fileInformation, errExpected)

	_, err := v.GetFileInformation(fileInformationGetter, nil)
	if err != errExpected {
		t.Errorf("GetFileInformation() error = %v, wantErr %v", err, errExpected)
		return
	}
}

func TestVideo_GetFileInformation_ShouldVideoWithInformation(t *testing.T) {
	ctx := context.Background()
	id := "JZAunPKoHL0"
	v := New(id)
	u, _ := url.Parse("https://www.youtube.com/watch?v=" + v.ID)
	n := time.Now()
	format := &ytdl.Format{
		Itag: ytdl.Itag{
			Number: 32,
		},
	}
	details := &ytdl.VideoInfo{
		ID:          v.ID,
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
	fileUrlGetter := new(MockFileUrlGetter)
	fileUrlGetter.On("GetDownloadURL", ctx, details, format).Return(u, nil)

	fileInformationGetter := new(MockFileInformationGetter)
	fileInformationGetter.On("GetVideoInfo", ctx, v.URL).Return(details, nil)

	result, err := v.GetFileInformation(fileInformationGetter, fileUrlGetter)
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
	v := New(id)
	u, _ := url.Parse("https://www.youtube.com/watch?v=" + v.ID)
	n := time.Now()
	format := &ytdl.Format{
		Itag: ytdl.Itag{
			Number: 32,
		},
	}
	details := &ytdl.VideoInfo{
		ID:          v.ID,
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
	fileUrlGetter := new(MockFileUrlGetter)
	fileUrlGetter.On("GetDownloadURL", ctx, details, format).Return(u, errWanted)

	fileInformationGetter := new(MockFileInformationGetter)
	fileInformationGetter.On("GetVideoInfo", ctx, v.URL).Return(details, nil)

	_, err := v.GetFileInformation(fileInformationGetter, fileUrlGetter)
	if err != errWanted {
		t.Errorf("GetFileInformation() error = %v, wantErr %v", err, nil)
		return
	}
}
