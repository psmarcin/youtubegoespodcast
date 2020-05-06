package youtube

import (
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestGetTrendings(t *testing.T) {
	tests := []struct {
		name    string
		want    YoutubeResponse
		wantErr bool
		status  int
	}{
		{
			name: "get trending",
			want: YoutubeResponse{
				Items: []Items{
					{
						Snippet: Snippet{
							ChannelID: "UCANLZYMidaCbLQFWXBC95Jg",
							Title:     "Taylor Swift - You Need To Calm Down",
						},
					},
					{
						Snippet: Snippet{
							ChannelID: "UClW4jraMKz6Qj69lJf-tODA",
							Title:     "NBA Youngboy - 4 Sons of a King (Official Audio)",
						},
					},
					{
						Snippet: Snippet{
							ChannelID: "UCuVHOs0H5hvAHGr8O4yIBNQ",
							Title:     "Swapping Boyfriends on Vacation",
						},
					},
					{
						Snippet: Snippet{
							ChannelID: "UCX6OQ3DkcsbYNE6H8uQQuVA",
							Title:     "I Opened The World's First FREE Store",
						},
					},
					{
						Snippet: Snippet{
							ChannelID: "UCpi8TJfiA4lKGkaXs__YdBA",
							Title:     "Why I'm Coming Out As Gay",
						},
					},
				},
			},
			wantErr: false,
			status:  http.StatusOK,
		},
		{
			name:    "should throw error 404",
			want:    YoutubeResponse{},
			wantErr: true,
			status:  http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			response := trending
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/videos").
				Reply(tt.status).
				BodyString(response)

			got, err := GetTrending()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrending(true) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, item := range got.Items {
				if item.Snippet.ChannelID != tt.want.Items[i].Snippet.ChannelID || item.Snippet.Title != tt.want.Items[i].Snippet.Title {
					t.Errorf("GetTrending(true) = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

const (
	trending = `{
 "kind": "youtube#videoListResponse",
 "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/5l1eMuCKeAZ60jh3ai3vUrlsJHQ\"",
 "nextPageToken": "CAUQAA",
 "pageInfo": {
  "totalResults": 200,
  "resultsPerPage": 5
 },
 "items": [
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/zZQZnMaiJCso_uFHS3-B_M-yjvc\"",
   "id": "Dkk9gvTmCXY",
   "snippet": {
    "publishedAt": "2019-06-17T12:13:12.000Z",
    "channelId": "UCANLZYMidaCbLQFWXBC95Jg",
    "title": "Taylor Swift - You Need To Calm Down",
    "description": "Music video by Taylor Swift performing “You Need To Calm Down” – off her upcoming new album “Lover” (out August 23).\n \n►Download “You Need To Calm Down\" here: https://TaylorSwift.lnk.to/YNTCDsu\n►Pre-Order “Lover” here: https://TaylorSwift.lnk.to/Loversu\n \n►Exclusive Merch: https://store.taylorswift.com\n \n►Follow Taylor Swift Online\nInstagram: http://www.instagram.com/taylorswift\nFacebook: http://www.facebook.com/taylorswift\nTumblr: http://taylorswift.tumblr.com\nTwitter: http://www.twitter.com/taylorswift13\nWebsite: http://www.taylorswift.com\n \n►Follow Taylor Nation Online\nInstagram: http://www.instagram.com/taylornation\nTumblr: http://taylornation.tumblr.com\nTwitter: http://www.twitter.com/taylornation13\n \nFeaturing appearances by (in alphabetical order):\nA'keria Davenport\nAdam Lambert\nAdam Rippon\nAdore Delano\nAntoni Porowski\nBilly Porter\nBobby Berk\nChester Lockhart\nCiara\nDelta Work\nDexter Mayfield\nEllen Degeneres\nHannah Hart\nHayley Kiyoko\nJade Jolie Jesse\nTyler Ferguson\nJonathan Van Ness\nJustin Mikita\nKaramo Brown\nKaty Perry\nLaverne Cox\nRiley Knox-Beyoncé\nRupaul\nRyan Reynolds\nTan France\nTatianna\nToderick Hall\nTrinity K Bonet\nTrinity Taylor\n \nDirectors: Drew Kirsch & Taylor Swift\nExecutive Producers: Todrick Hall & Taylor Swift\n \n© 2019 Taylor Swift\n \n#TaylorSwift #YouNeedToCalmDown #TaylorSwiftLover",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/Dkk9gvTmCXY/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/Dkk9gvTmCXY/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/Dkk9gvTmCXY/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/Dkk9gvTmCXY/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/Dkk9gvTmCXY/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "TaylorSwiftVEVO",
    "tags": [
     "Taylor",
     "Swift",
     "You",
     "Need",
     "To",
     "Calm",
     "Down",
     "Pop"
    ],
    "categoryId": "10",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "Taylor Swift - You Need To Calm Down",
     "description": "Music video by Taylor Swift performing “You Need To Calm Down” – off her upcoming new album “Lover” (out August 23).\n \n►Download “You Need To Calm Down\" here: https://TaylorSwift.lnk.to/YNTCDsu\n►Pre-Order “Lover” here: https://TaylorSwift.lnk.to/Loversu\n \n►Exclusive Merch: https://store.taylorswift.com\n \n►Follow Taylor Swift Online\nInstagram: http://www.instagram.com/taylorswift\nFacebook: http://www.facebook.com/taylorswift\nTumblr: http://taylorswift.tumblr.com\nTwitter: http://www.twitter.com/taylorswift13\nWebsite: http://www.taylorswift.com\n \n►Follow Taylor Nation Online\nInstagram: http://www.instagram.com/taylornation\nTumblr: http://taylornation.tumblr.com\nTwitter: http://www.twitter.com/taylornation13\n \nFeaturing appearances by (in alphabetical order):\nA'keria Davenport\nAdam Lambert\nAdam Rippon\nAdore Delano\nAntoni Porowski\nBilly Porter\nBobby Berk\nChester Lockhart\nCiara\nDelta Work\nDexter Mayfield\nEllen Degeneres\nHannah Hart\nHayley Kiyoko\nJade Jolie Jesse\nTyler Ferguson\nJonathan Van Ness\nJustin Mikita\nKaramo Brown\nKaty Perry\nLaverne Cox\nRiley Knox-Beyoncé\nRupaul\nRyan Reynolds\nTan France\nTatianna\nToderick Hall\nTrinity K Bonet\nTrinity Taylor\n \nDirectors: Drew Kirsch & Taylor Swift\nExecutive Producers: Todrick Hall & Taylor Swift\n \n© 2019 Taylor Swift\n \n#TaylorSwift #YouNeedToCalmDown #TaylorSwiftLover"
    }
   }
  },
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/z0SzB1oFmmcblPH8T57ijcoPrAk\"",
   "id": "fVVnQyRhtTk",
   "snippet": {
    "publishedAt": "2019-06-16T05:19:25.000Z",
    "channelId": "UClW4jraMKz6Qj69lJf-tODA",
    "title": "NBA Youngboy - 4 Sons of a King (Official Audio)",
    "description": "YoungBoy Never Broke Again",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/fVVnQyRhtTk/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/fVVnQyRhtTk/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/fVVnQyRhtTk/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/fVVnQyRhtTk/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/fVVnQyRhtTk/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "YoungBoy Never Broke Again",
    "tags": [
     "YoungBoy Never Broke Again",
     "YBNBA",
     "YoungBoy NBA",
     "NBA YoungBoy",
     "YB",
     "NBA YB",
     "Rap",
     "Hip Hop",
     "Atlantic Records",
     "Until Death Call My Name",
     "Master The Day Of Judgement"
    ],
    "categoryId": "10",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "NBA Youngboy - 4 Sons of a King (Official Audio)",
     "description": "YoungBoy Never Broke Again"
    }
   }
  },
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/-LfcB5Cz56kTQldEM234_t6m-CQ\"",
   "id": "paZaYFZbZtM",
   "snippet": {
    "publishedAt": "2019-06-16T21:03:38.000Z",
    "channelId": "UCuVHOs0H5hvAHGr8O4yIBNQ",
    "title": "Swapping Boyfriends on Vacation",
    "description": "Swapping Boyfriends on Vacation... (ahhhh)\nWe've done boyfriend swaps before, but nothing quite like this one. We swapped boyfriends for the nigh, but this time we're swapping boyfriends on vacation and going on romantic summer dates with each other's boyfriends.  Let the cringe begin. #swapping #boyfriends\nSubscribe here ➜ http://bit.ly/2vxi9ch\nmore swap videos below:\nopposite twins Swap apartments for a day:\nhttps://www.youtube.com/watch?v=3l1Dr8Ge7Go\nsisters swap clothes for a week: \nhttps://www.youtube.com/watch?v=taA4g...\nopposite twins swap lives for a day: https://www.youtube.com/watch?v=Mhzmc...\nopposite twins swap boyfriends for the weekend:\nhttps://www.youtube.com/watch?v=3iXq8...\nhilarious twin style swap & transformation:\nhttps://www.youtube.com/watch?v=pE5go...\n\n**NEW VIDEOS EVERY SUNDAY and SOMETIMES WEDNESDAY**\n\nIf you see this, comment \"Nate and Gabi getting physical woah\" \nonly those who watch up to that point will know what this means ;)\n\nvlog channels:\nniki demar\nhttps://www.youtube.com/user/nikidemar\nfancy vlogs by gab https://www.youtube.com/channel/UCLGe...\n\nSOCIAL MEDIA \nInstagram➜ @NIKI / @GABI\nTwitter➜ @nikidemar / @gabcake\nTumblr➜ nikidemar / breakfastatchanel-starringgabi\nSnapchat➜ nikidemarrr / fancysnapsbygab\n\nWe’re Niki and Gabi! We hope you enjoyed our swapping boyfriends on vacation video! We’re twin sisters who are different with opposite fashion and styles, but we come together to make videos like challenges, swaps, shopping challenges, 24 hour challenges, diy, style, beauty, lifestyle, fashion, comedy, types of girls, music, and more!",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/paZaYFZbZtM/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/paZaYFZbZtM/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/paZaYFZbZtM/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/paZaYFZbZtM/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/paZaYFZbZtM/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "Niki and Gabi",
    "tags": [
     "swapping boyfriends on vacation",
     "niki and gabi",
     "swapping boyfriends",
     "swapping",
     "swap",
     "swap boyfriends",
     "boyfriend swap",
     "sisters swap boyfriends",
     "niki and gabi boyfriends"
    ],
    "categoryId": "26",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "Swapping Boyfriends on Vacation",
     "description": "Swapping Boyfriends on Vacation... (ahhhh)\nWe've done boyfriend swaps before, but nothing quite like this one. We swapped boyfriends for the nigh, but this time we're swapping boyfriends on vacation and going on romantic summer dates with each other's boyfriends.  Let the cringe begin. #swapping #boyfriends\nSubscribe here ➜ http://bit.ly/2vxi9ch\nmore swap videos below:\nopposite twins Swap apartments for a day:\nhttps://www.youtube.com/watch?v=3l1Dr8Ge7Go\nsisters swap clothes for a week: \nhttps://www.youtube.com/watch?v=taA4g...\nopposite twins swap lives for a day: https://www.youtube.com/watch?v=Mhzmc...\nopposite twins swap boyfriends for the weekend:\nhttps://www.youtube.com/watch?v=3iXq8...\nhilarious twin style swap & transformation:\nhttps://www.youtube.com/watch?v=pE5go...\n\n**NEW VIDEOS EVERY SUNDAY and SOMETIMES WEDNESDAY**\n\nIf you see this, comment \"Nate and Gabi getting physical woah\" \nonly those who watch up to that point will know what this means ;)\n\nvlog channels:\nniki demar\nhttps://www.youtube.com/user/nikidemar\nfancy vlogs by gab https://www.youtube.com/channel/UCLGe...\n\nSOCIAL MEDIA \nInstagram➜ @NIKI / @GABI\nTwitter➜ @nikidemar / @gabcake\nTumblr➜ nikidemar / breakfastatchanel-starringgabi\nSnapchat➜ nikidemarrr / fancysnapsbygab\n\nWe’re Niki and Gabi! We hope you enjoyed our swapping boyfriends on vacation video! We’re twin sisters who are different with opposite fashion and styles, but we come together to make videos like challenges, swaps, shopping challenges, 24 hour challenges, diy, style, beauty, lifestyle, fashion, comedy, types of girls, music, and more!"
    },
    "defaultAudioLanguage": "en"
   }
  },
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/yMwK_HgUJQODgHhv20ENc6GKDAI\"",
   "id": "n6qc4LHN2KQ",
   "snippet": {
    "publishedAt": "2019-06-15T20:00:00.000Z",
    "channelId": "UCX6OQ3DkcsbYNE6H8uQQuVA",
    "title": "I Opened The World's First FREE Store",
    "description": "EVERYTHING IN THE STORE WAS %100 FREE!!!!!\n\nNew Merch - https://shopmrbeast.com/\n\nSUBSCRIBE OR I TAKE YOUR DOG\n\n\n----------------------------------------------------------------\nfollow all of these or i will kick you\n• Facebook - https://www.facebook.com/MrBeast6000/\n• Twitter - https://twitter.com/MrBeastYT\n•  Instagram - https://www.instagram.com/mrbeast\n--------------------------------------------------------------------",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/n6qc4LHN2KQ/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/n6qc4LHN2KQ/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/n6qc4LHN2KQ/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/n6qc4LHN2KQ/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/n6qc4LHN2KQ/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "MrBeast",
    "categoryId": "24",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "I Opened The World's First FREE Store",
     "description": "EVERYTHING IN THE STORE WAS %100 FREE!!!!!\n\nNew Merch - https://shopmrbeast.com/\n\nSUBSCRIBE OR I TAKE YOUR DOG\n\n\n----------------------------------------------------------------\nfollow all of these or i will kick you\n• Facebook - https://www.facebook.com/MrBeast6000/\n• Twitter - https://twitter.com/MrBeastYT\n•  Instagram - https://www.instagram.com/mrbeast\n--------------------------------------------------------------------"
    }
   }
  },
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/SrYEfi9xxTHv53hbKwfvugh_nZY\"",
   "id": "QruHsyt8paY",
   "snippet": {
    "publishedAt": "2019-06-17T15:04:17.000Z",
    "channelId": "UCpi8TJfiA4lKGkaXs__YdBA",
    "title": "Why I'm Coming Out As Gay",
    "description": "Behind the scenes of Eugene Lee Yang's coming out video, \"I'm Gay:\" a never-before-seen look inside the mind of a young artist and his complex, dark struggles with his identity, his work, his public persona, and his private life. \n\nA documentary video by Kane Diep @kanediep\nhttps://www.kanediep.com/\n\n⚡THE TRY GUYS LEGENDS OF THE INTERNET TOUR ⚡: tickets available on Friday, May 10 at 12PM EST / 9AM PST at https://tryguys.com/tour\n\n🎧THE TRYPOD 🎧: watch our new podcast at https://youtube.com/trypod or listen at https://tryguys.com/podcast\n\n📘THE HIDDEN POWER OF F*CKING UP 📘: check out our new book at https://tryguys.com/book\n\nGet your official Try Guys color hoodies and phone cases at https://tryguys.com/collections/color... 💙❤️💚💜\n\nSupport us! http://www.patreon.com/tryguys. Join our Patreon to get videos a day early, plus, live streams, chatrooms, BTS footage, exclusive merchandise, and more!\n\nSUBSCRIBE TO AND FOLLOW THE TRY GUYS \nhttp://www.youtube.com/c/tryguys\nhttp://www.facebook.com/tryguys  \nhttp://www.twitter.com/tryguys\nhttps://www.instagram.com/tryguys\n\nFOLLOW THE GUYS\nhttp://www.Instagram.com/keithhabs\nhttp://www.Instagram.com/nedfulmer\nhttp://www.Instagram.com/korndiddy\nhttp://www.instagram.com/eugeneleeyang\n \nhttp://www.twitter.com/keithhabs\nhttp://www.twitter.com/nedfulmer \nhttp://www.twitter.com/korndiddy\nhttp://www.twitter.com/eugeneleeyang \n\nTHE TRY GUYS\nThe #TryGuys is the flagship channel of 2ND TRY, LLC.  Tune in twice a week for shows from Keith, Ned, Zach and Eugene, the creators and stars of The Try Guys.\n\nMUSIC\nLicensed from AudioNetwork\n\nSFX\nLicensed from Audioblocks\n\nVIDEO\nLicensed from Videoblocks\n\nOfficial Try Guys Photos \nBy Mandee Johnson Photography | @mandeephoto\n\n2nd Try, LLC STAFF\nExecutive Producer - Keith Habersberger\nExecutive Producer - Ned Fulmer\nExecutive Producer - Zach Kornfeld\nExecutive Producer - Eugene Lee Yang\nProducer - Rachel Ann Cole\nProducer - Nick Rufca\nProduction Manager - Alexandria Herring\nEditor - Devlin McCluskey\nEditor - YB Chang\nEditor - Elliot Dickerhoof\nAssistant Editor - Will Witwer\nCamera Operator - Miles Bonsignore\nProduction Assistant - Sam Johnson\nSocial Media Manager - Izzy Francke",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/QruHsyt8paY/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/QruHsyt8paY/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/QruHsyt8paY/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/QruHsyt8paY/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/QruHsyt8paY/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "The Try Guys",
    "tags": [
     "try guys",
     "keith",
     "ned",
     "zach",
     "eugene",
     "habersberger",
     "fulmer",
     "kornfeld",
     "yang",
     "buzzfeedvideo",
     "buzzfeed",
     "ariel",
     "ned & ariel",
     "comedy",
     "education",
     "funny",
     "try",
     "learn",
     "fail",
     "experiment",
     "test",
     "tryceratops",
     "eugene lee yang",
     "gay",
     "I'm",
     "I'm Gay",
     "LGBT",
     "LGBTQ",
     "queer",
     "identity",
     "coming out",
     "secret",
     "family",
     "dance",
     "behind the scenes",
     "bts",
     "happiness",
     "sad",
     "viral",
     "The Trevor Project",
     "art",
     "music video",
     "process",
     "documentary",
     "why",
     "struggle",
     "styling",
     "cinematography",
     "choreography",
     "producer",
     "artist",
     "stress"
    ],
    "categoryId": "23",
    "liveBroadcastContent": "none",
    "defaultLanguage": "en",
    "localized": {
     "title": "Why I'm Coming Out As Gay",
     "description": "Behind the scenes of Eugene Lee Yang's coming out video, \"I'm Gay:\" a never-before-seen look inside the mind of a young artist and his complex, dark struggles with his identity, his work, his public persona, and his private life. \n\nA documentary video by Kane Diep @kanediep\nhttps://www.kanediep.com/\n\n⚡THE TRY GUYS LEGENDS OF THE INTERNET TOUR ⚡: tickets available on Friday, May 10 at 12PM EST / 9AM PST at https://tryguys.com/tour\n\n🎧THE TRYPOD 🎧: watch our new podcast at https://youtube.com/trypod or listen at https://tryguys.com/podcast\n\n📘THE HIDDEN POWER OF F*CKING UP 📘: check out our new book at https://tryguys.com/book\n\nGet your official Try Guys color hoodies and phone cases at https://tryguys.com/collections/color... 💙❤️💚💜\n\nSupport us! http://www.patreon.com/tryguys. Join our Patreon to get videos a day early, plus, live streams, chatrooms, BTS footage, exclusive merchandise, and more!\n\nSUBSCRIBE TO AND FOLLOW THE TRY GUYS \nhttp://www.youtube.com/c/tryguys\nhttp://www.facebook.com/tryguys  \nhttp://www.twitter.com/tryguys\nhttps://www.instagram.com/tryguys\n\nFOLLOW THE GUYS\nhttp://www.Instagram.com/keithhabs\nhttp://www.Instagram.com/nedfulmer\nhttp://www.Instagram.com/korndiddy\nhttp://www.instagram.com/eugeneleeyang\n \nhttp://www.twitter.com/keithhabs\nhttp://www.twitter.com/nedfulmer \nhttp://www.twitter.com/korndiddy\nhttp://www.twitter.com/eugeneleeyang \n\nTHE TRY GUYS\nThe #TryGuys is the flagship channel of 2ND TRY, LLC.  Tune in twice a week for shows from Keith, Ned, Zach and Eugene, the creators and stars of The Try Guys.\n\nMUSIC\nLicensed from AudioNetwork\n\nSFX\nLicensed from Audioblocks\n\nVIDEO\nLicensed from Videoblocks\n\nOfficial Try Guys Photos \nBy Mandee Johnson Photography | @mandeephoto\n\n2nd Try, LLC STAFF\nExecutive Producer - Keith Habersberger\nExecutive Producer - Ned Fulmer\nExecutive Producer - Zach Kornfeld\nExecutive Producer - Eugene Lee Yang\nProducer - Rachel Ann Cole\nProducer - Nick Rufca\nProduction Manager - Alexandria Herring\nEditor - Devlin McCluskey\nEditor - YB Chang\nEditor - Elliot Dickerhoof\nAssistant Editor - Will Witwer\nCamera Operator - Miles Bonsignore\nProduction Assistant - Sam Johnson\nSocial Media Manager - Izzy Francke"
    },
    "defaultAudioLanguage": "en"
   }
  }
 ]
}`
)
