package feed

import (
	"strings"
	"testing"
)

const (
	emptyTitle = `<title></title>`
)

func TestFeed_serialize(t *testing.T) {
	type fields struct {
		XMLName       string
		ChannelID     string
		Title         string
		Link          string
		Description   string
		Category      string
		Generator     string
		Language      string
		LastBuildDate string
		PubDate       string
		Image         Image
		ITAuthor      string
		ITSubtitle    string
		ITSummary     ITSummary
		ITImage       ITImage
		ITExplicit    string
		ITCategory    ITCategory
		Items         []Item
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "serialize success empty channel",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				XMLName:       tt.fields.XMLName,
				ChannelID:     tt.fields.ChannelID,
				Title:         tt.fields.Title,
				Link:          tt.fields.Link,
				Description:   tt.fields.Description,
				Category:      tt.fields.Category,
				Generator:     tt.fields.Generator,
				Language:      tt.fields.Language,
				LastBuildDate: tt.fields.LastBuildDate,
				PubDate:       tt.fields.PubDate,
				Image:         tt.fields.Image,
				ITAuthor:      tt.fields.ITAuthor,
				ITSubtitle:    tt.fields.ITSubtitle,
				ITSummary:     tt.fields.ITSummary,
				ITImage:       tt.fields.ITImage,
				ITExplicit:    tt.fields.ITExplicit,
				ITCategory:    tt.fields.ITCategory,
				Items:         tt.fields.Items,
			}
			got, err := f.serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), emptyTitle) {
				t.Errorf("Feed.serialize() = %v, want %v", got, emptyTitle)
			}
		})
	}
}
