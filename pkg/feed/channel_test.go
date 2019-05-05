package feed

import "testing"

func TestFeed_addItem(t *testing.T) {
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
	type args struct {
		item Item
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:   "fail to add item - no title and enclosure",
			fields: fields{},
			args: args{
				item: Item{
					GUID: "1",
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "fail to add item - no title",
			fields: fields{},
			args: args{
				item: Item{
					Enclosure: Enclosure{
						URL: "url",
					},
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "add item",
			fields: fields{},
			args: args{
				item: Item{
					Enclosure: Enclosure{
						URL: "url",
					},
					Title: "Example title",
				},
			},
			wantCount: 1,
			wantErr:   false,
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
			if err := f.addItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Feed.addItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(f.Items) != tt.wantCount {
				t.Errorf("Feed.addItem() error = %v, want %v", len(f.Items), tt.wantCount)
			}
		})
	}
}
