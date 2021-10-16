package feed

import (
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const (
	rawCategoryPrefix = "https://en.wikipedia.org/wiki/"
)

type Category struct {
	Name     string
	Score    uint32
	Children []Category
}

var Categories = []Category{
	{Name: "Arts"},
	{Name: "Business"},
	{Name: "Comedy"},
	{Name: "Education"},
	{Name: "Fiction"},
	{Name: "Government"},
	{Name: "History"},
	{Name: "Health & Fitness"},
	{Name: "Kids & Family"},
	{Name: "Leisure"},
	{Name: "Music"},
	{Name: "News"},
	{Name: "Religion & Spirituality"},
	{Name: "Science"},
	{Name: "Society & Culture"},
	{Name: "Sports"},
	{Name: "Technology"},
	{Name: "True Crime"},
	{Name: "TV & Film"},
}

func SelectCategory(rawCategories []string) string {
	best := Category{}
	categories := cleanupRawCategories(rawCategories)

	for _, category := range categories {
		foundCategory := findCategory(category)
		if (best.Score > foundCategory.Score) || best.Name == "" {
			best = foundCategory
		}
		l.Infof("looking at category: %s and found: %s with score: %d", category, foundCategory.Name, foundCategory.Score)
	}

	l.Infof("best category: %s with score: %d", best.Name, best.Score)

	return best.Name
}

func cleanupRawCategories(categories []string) []string {
	var cleanCategories []string
	for _, cat := range categories {
		cleanString := strings.ToLower(strings.TrimPrefix(cat, rawCategoryPrefix))
		if cleanString != "" {
			cleanCategories = append(cleanCategories, cleanString)
		}
	}

	return cleanCategories
}

func findCategory(category string) Category {
	var found Category

	for _, cat := range Categories {
		distance := levenshtein.
			DistanceForStrings([]rune(strings.ToLower(cat.Name)), []rune(category), levenshtein.DefaultOptions)
		cat.Score = uint32(distance)

		if cat.Score < found.Score || found.Name == "" {
			found = cat
		}
	}

	return found
}
