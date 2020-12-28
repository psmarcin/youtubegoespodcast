package feed

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findCategory(t *testing.T) {
	type args struct {
		category string
	}
	tests := []struct {
		name string
		args args
		want Category
	}{
		{
			name: "should return Sports for sport",
			args: args{
				category: "sport",
			},
			want: Category{
				Name:     "Sports",
				Score:    1,
				Children: nil,
			},
		},
		{
			name: "should return Arts for arts",
			args: args{
				category: "arts",
			},
			want: Category{
				Name:     "Arts",
				Score:    0,
				Children: nil,
			},
		},
		{
			name: "should return Leisure for video_game_culture",
			args: args{
				category: "video_game_culture",
			},
			want: Category{
				Name:     "Leisure",
				Score:    17,
				Children: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findCategory(tt.args.category); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanupRawCategories(t *testing.T) {
	var n []string

	type args struct {
		categories []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should remove prefix from: " + rawCategoryPrefix + "123",
			args: args{
				categories: []string{rawCategoryPrefix + "123"},
			},
			want: []string{"123"},
		},
		{
			name: "should convert to lowercase TO_LOWER-CASE",
			args: args{
				categories: []string{"TO_LOWER-CASE"},
			},
			want: []string{"to_lower-case"},
		},
		{
			name: "should omit empty string",
			args: args{
				categories: []string{""},
			},
			want: n,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanupRawCategories(tt.args.categories)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestSelectCategory(t *testing.T) {
	type args struct {
		rawCategories []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should find Music for list of categories",
			args: args{
				rawCategories: []string{
					"https://en.wikipedia.org/wiki/Independent_music",
					"https://en.wikipedia.org/wiki/Pop_music",
					"https://en.wikipedia.org/wiki/Electronic_music",
					"https://en.wikipedia.org/wiki/Music",
				},
			},
			want: "Music",
		},
		{
			name: "should find empty string for empty list of categories",
			args: args{
				rawCategories: []string{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectCategory(tt.args.rawCategories)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
