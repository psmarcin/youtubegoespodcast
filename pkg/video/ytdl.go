package video

// import (
// 	"bytes"
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/iawia002/annie/request"
// 	log "github.com/sirupsen/logrus"
// )

// // VideoInfo contains the info a youtube video
// type VideoInfo struct {
// 	// The video ID
// 	ID string `json:"id"`
// 	// The video title
// 	Title string `json:"title"`
// 	// The video description
// 	Description string `json:"description"`
// 	// The date the video was published
// 	DatePublished time.Time `json:"datePublished"`
// 	// Formats the video is available in
// 	Formats FormatList `json:"formats"`
// 	// List of keywords associated with the video
// 	Keywords []string `json:"keywords"`
// 	// Author of the video
// 	Author string `json:"author"`
// 	// Duration of the video
// 	Duration time.Duration

// 	htmlPlayerFile string
// }

// // FormatList is a slice of formats with filtering functionality
// type FormatList []Format

// // Format is a youtube is a static youtube video format
type Format struct {
	Itag          int    `json:"itag"`
	Extension     string `json:"extension"`
	Resolution    string `json:"resolution"`
	VideoEncoding string `json:"videoEncoding"`
	AudioEncoding string `json:"audioEncoding"`
	AudioBitrate  int    `json:"audioBitrate"`
	meta          map[string]interface{}
}

// // FormatKey is a string type containing a key in a video format map
// type FormatKey string

const youtubeBaseURL = "https://www.youtube.com/watch"

// const youtubeEmbededBaseURL = "https://www.youtube.com/embed/"
// const youtubeVideoEURL = "https://youtube.googleapis.com/v/"
// const youtubeVideoInfoURL = "https://www.youtube.com/get_video_info"
// const youtubeDateFormat = "2006-01-02"
// const referer = "https://www.youtube.com"

// // Available format Keys
// const (
// 	FormatExtensionKey     FormatKey = "ext"
// 	FormatResolutionKey    FormatKey = "res"
// 	FormatVideoEncodingKey FormatKey = "videnc"
// 	FormatAudioEncodingKey FormatKey = "audenc"
// 	FormatItagKey          FormatKey = "itag"
// 	FormatAudioBitrateKey  FormatKey = "audbr"
// 	FormatFPSKey           FormatKey = "fps"
// )

// func GetVideoInfoFromID(id string) (*VideoInfo, error) {
// 	u, _ := url.ParseRequestURI(youtubeBaseURL)
// 	values := u.Query()
// 	values.Set("v", id)
// 	u.RawQuery = values.Encode()

// 	resp, err := http.Get(u.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return getVideoInfoFromHTML(id, body)
// }

// func getVideoInfoFromHTML(id string, html []byte) (*VideoInfo, error) {
// 	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
// 	if err != nil {
// 		return nil, err
// 	}

// 	info := &VideoInfo{}

// 	// extract description and title
// 	info.Description = strings.TrimSpace(doc.Find("#eow-description").Text())
// 	info.Title = strings.TrimSpace(doc.Find("#eow-title").Text())
// 	info.ID = id
// 	dateStr, ok := doc.Find("meta[itemprop=\"datePublished\"]").Attr("content")
// 	if !ok {
// 		log.Debug("Unable to extract date published")
// 	} else {
// 		date, err := time.Parse(youtubeDateFormat, dateStr)
// 		if err == nil {
// 			info.DatePublished = date
// 		} else {
// 			log.Debug("Unable to parse date published", err.Error())
// 		}
// 	}

// 	// match json in javascript
// 	re := regexp.MustCompile("ytplayer.config = (.*?);ytplayer.load")
// 	matches := re.FindSubmatch(html)
// 	var jsonConfig map[string]interface{}
// 	if len(matches) > 1 {
// 		err = json.Unmarshal(matches[1], &jsonConfig)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		log.Debug("Unable to extract json from default url, trying embedded url")
// 		var resp *http.Response
// 		resp, err = http.Get(youtubeEmbededBaseURL + id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != 200 {
// 			return nil, fmt.Errorf("Embeded url request returned status code %d	", resp.StatusCode)
// 		}
// 		html, err = ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		//	re = regexp.MustCompile("\"sts\"\\s*:\\s*(\\d+)")
// 		re = regexp.MustCompile("yt.setConfig\\({'PLAYER_CONFIG': (.*?)}\\);")

// 		matches := re.FindSubmatch(html)
// 		if len(matches) < 2 {
// 			return nil, fmt.Errorf("Error extracting sts from embedded url response")
// 		}
// 		dec := json.NewDecoder(bytes.NewBuffer(matches[1]))
// 		err = dec.Decode(&jsonConfig)
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to extract json from embedded url: %s", err.Error())
// 		}
// 		query := url.Values{
// 			"sts":      []string{strconv.Itoa(int(jsonConfig["sts"].(float64)))},
// 			"video_id": []string{id},
// 			"eurl":     []string{youtubeVideoEURL + id},
// 		}

// 		resp, err = http.Get(youtubeVideoInfoURL + "?" + query.Encode())
// 		if err != nil {
// 			return nil, fmt.Errorf("Error fetching video info: %s", err.Error())
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != 200 {
// 			return nil, fmt.Errorf("Video info response invalid status code")
// 		}
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to read video info response body: %s", err.Error())
// 		}
// 		query, err = url.ParseQuery(string(body))
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to parse video info data: %s", err.Error())
// 		}
// 		args := make(map[string]interface{})
// 		for k, v := range query {
// 			if len(v) > 0 {
// 				args[k] = v[0]
// 			}
// 		}
// 		jsonConfig["args"] = args
// 	}

// 	inf := jsonConfig["args"].(map[string]interface{})
// 	if status, ok := inf["status"].(string); ok && status == "fail" {
// 		return nil, fmt.Errorf("Error %d:%s", inf["errorcode"], inf["reason"])
// 	}
// 	if a, ok := inf["author"].(string); ok {
// 		info.Author = a
// 	} else {
// 		log.Debug("Unable to extract author")
// 	}

// 	if length, ok := inf["length_seconds"].(string); ok {
// 		if duration, err := strconv.ParseInt(length, 10, 64); err == nil {
// 			info.Duration = time.Second * time.Duration(duration)
// 		} else {
// 			log.Debug("Unable to parse duration string: ", length)
// 		}
// 	} else {
// 		log.Debug("Unable to extract duration")
// 	}

// 	// For the future maybe
// 	parseKey := func(key string) []string {
// 		val, ok := inf[key].(string)
// 		if !ok {
// 			return nil
// 		}
// 		vals := []string{}
// 		split := strings.Split(val, ",")
// 		for _, v := range split {
// 			if v != "" {
// 				vals = append(vals, v)
// 			}
// 		}
// 		return vals
// 	}
// 	info.Keywords = parseKey("keywords")
// 	info.htmlPlayerFile = jsonConfig["assets"].(map[string]interface{})["js"].(string)

// 	/*
// 		fmtList := parseKey("fmt_list")
// 		fexp := parseKey("fexp")
// 		watermark := parseKey("watermark")

// 		if len(fmtList) != 0 {
// 			vals := []string{}
// 			for _, v := range fmtList {
// 				vals = append(vals, strings.Split(v, "/")...)
// 		} else {
// 			info["fmt_list"] = []string{}
// 		}

// 		videoVerticals := []string{}
// 		if videoVertsStr, ok := inf["video_verticals"].(string); ok {
// 			videoVertsStr = string([]byte(videoVertsStr)[1 : len(videoVertsStr)-2])
// 			videoVertsSplit := strings.Split(videoVertsStr, ", ")
// 			for _, v := range videoVertsSplit {
// 				if v != "" {
// 					videoVerticals = append(videoVerticals, v)
// 				}
// 			}
// 		}
// 	*/
// 	var formatStrings []string
// 	if fmtStreamMap, ok := inf["url_encoded_fmt_stream_map"].(string); ok {
// 		formatStrings = append(formatStrings, strings.Split(fmtStreamMap, ",")...)
// 	}

// 	if adaptiveFormats, ok := inf["adaptive_fmts"].(string); ok {
// 		formatStrings = append(formatStrings, strings.Split(adaptiveFormats, ",")...)
// 	}
// 	var formats FormatList
// 	for _, v := range formatStrings {
// 		query, err := url.ParseQuery(v)
// 		if err == nil {
// 			itag, _ := strconv.Atoi(query.Get("itag"))
// 			if format, ok := newFormat(itag); ok {
// 				if strings.HasPrefix(query.Get("conn"), "rtmp") {
// 					format.meta["rtmp"] = true
// 				}
// 				for k, v := range query {
// 					if len(v) == 1 {
// 						format.meta[k] = v[0]
// 					} else {
// 						format.meta[k] = v
// 					}
// 				}
// 				formats = append(formats, format)
// 			} else {
// 				log.Debug("No metadata found for itag: ", itag, ", skipping...")
// 			}
// 		} else {
// 			log.Debug("Unable to format string", err.Error())
// 		}
// 	}

// 	if dashManifestURL, ok := inf["dashmpd"].(string); ok {
// 		tokens, err := getSigTokens(info.htmlPlayerFile)
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to extract signature tokens: %s", err.Error())
// 		}
// 		regex := regexp.MustCompile("\\/s\\/([a-fA-F0-9\\.]+)")
// 		regexSub := regexp.MustCompile("([a-fA-F0-9\\.]+)")
// 		dashManifestURL = regex.ReplaceAllStringFunc(dashManifestURL, func(str string) string {
// 			return "/signature/" + decipherTokens(tokens, regexSub.FindString(str))
// 		})
// 		dashFormats, err := getDashManifest(dashManifestURL)
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to extract dash manifest: %s", err.Error())
// 		}

// 		for _, dashFormat := range dashFormats {
// 			added := false
// 			for j, format := range formats {
// 				if dashFormat.Itag == format.Itag {
// 					formats[j] = dashFormat
// 					added = true
// 					break
// 				}
// 			}
// 			if !added {
// 				formats = append(formats, dashFormat)
// 			}
// 		}
// 	}
// 	info.Formats = formats
// 	return info, nil
// }

// func newFormat(itag int) (Format, bool) {
// 	if f, ok := FORMATS[itag]; ok {
// 		f.meta = make(map[string]interface{})
// 		return f, true
// 	}
// 	return Format{}, false
// }

// func getSigTokens(htmlPlayerFile string) ([]string, error) {
// 	u, _ := url.Parse("https://www.youtube.com/watch")
// 	p, err := url.Parse(htmlPlayerFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := request.Get(u.ResolveReference(p).String(), referer, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	objResult := actionsObjRegexp.FindStringSubmatch(body)
// 	funcResult := actionsFuncRegexp.FindStringSubmatch(body)

// 	if len(objResult) < 3 || len(funcResult) < 2 {
// 		return nil, fmt.Errorf("error parsing signature tokens")
// 	}
// 	obj := strings.Replace(objResult[1], "$", "\\$", -1)
// 	objBody := strings.Replace(objResult[2], "$", "\\$", -1)
// 	funcBody := strings.Replace(funcResult[1], "$", "\\$", -1)

// 	var reverseKey, sliceKey, spliceKey, swapKey string
// 	var result []string

// 	if result = reverseRegexp.FindStringSubmatch(objBody); len(result) > 1 {
// 		reverseKey = strings.Replace(result[1], "$", "\\$", -1)
// 	}
// 	if result = sliceRegexp.FindStringSubmatch(objBody); len(result) > 1 {
// 		sliceKey = strings.Replace(result[1], "$", "\\$", -1)
// 	}
// 	if result = spliceRegexp.FindStringSubmatch(objBody); len(result) > 1 {
// 		spliceKey = strings.Replace(result[1], "$", "\\$", -1)
// 	}
// 	if result = swapRegexp.FindStringSubmatch(objBody); len(result) > 1 {
// 		swapKey = strings.Replace(result[1], "$", "\\$", -1)
// 	}

// 	keys := []string{reverseKey, sliceKey, spliceKey, swapKey}
// 	regex, err := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s)\\(a,(\\d+)\\)", obj, strings.Join(keys, "|")))
// 	if err != nil {
// 		return nil, err
// 	}
// 	results := regex.FindAllStringSubmatch(funcBody, -1)
// 	var tokens []string
// 	for _, s := range results {
// 		switch s[1] {
// 		case swapKey:
// 			tokens = append(tokens, "w"+s[2])
// 		case reverseKey:
// 			tokens = append(tokens, "r")
// 		case sliceKey:
// 			tokens = append(tokens, "s"+s[2])
// 		case spliceKey:
// 			tokens = append(tokens, "p"+s[2])
// 		}
// 	}
// 	return tokens, nil
// }

// func decipherTokens(tokens []string, sig string) string {
// 	var pos int
// 	sigSplit := strings.Split(sig, "")
// 	for i, l := 0, len(tokens); i < l; i++ {
// 		tok := tokens[i]
// 		if len(tok) > 1 {
// 			pos, _ = strconv.Atoi(string(tok[1:]))
// 			pos = ^^pos
// 		}
// 		switch string(tok[0]) {
// 		case "r":
// 			reverseStringSlice(sigSplit)
// 		case "w":
// 			s := sigSplit[0]
// 			sigSplit[0] = sigSplit[pos]
// 			sigSplit[pos] = s
// 		case "s":
// 			sigSplit = sigSplit[pos:]
// 		case "p":
// 			sigSplit = sigSplit[pos:]
// 		}
// 	}
// 	return strings.Join(sigSplit, "")
// }

// func getDashManifest(urlString string) (formats []Format, err error) {

// 	resp, err := http.Get(urlString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
// 	}
// 	dec := xml.NewDecoder(resp.Body)
// 	var token xml.Token
// 	for ; err == nil; token, err = dec.Token() {
// 		if el, ok := token.(xml.StartElement); ok && el.Name.Local == "Representation" {
// 			var rep representation
// 			err = dec.DecodeElement(&rep, &el)
// 			if err != nil {
// 				break
// 			}
// 			if format, ok := newFormat(rep.Itag); ok {
// 				format.meta["url"] = rep.URL
// 				if rep.Height != 0 {
// 					format.Resolution = strconv.Itoa(rep.Height) + "p"
// 				} else {
// 					format.Resolution = ""
// 				}
// 				formats = append(formats, format)
// 			} else {
// 				log.Debug("No metadata found for itag: ", rep.Itag, ", skipping...")
// 			}
// 		}
// 	}
// 	if err != io.EOF {
// 		return nil, err
// 	}
// 	return formats, nil
// }

// // FORMATS is a map of all itags and their formats
var FORMATS = map[int]Format{
	5: {
		Extension:     "flv",
		Resolution:    "240p",
		VideoEncoding: "Sorenson H.283",
		AudioEncoding: "mp3",
		Itag:          5,
		AudioBitrate:  64,
	},
	6: {
		Extension:     "flv",
		Resolution:    "270p",
		VideoEncoding: "Sorenson H.263",
		AudioEncoding: "mp3",
		Itag:          6,
		AudioBitrate:  64,
	},
	13: {
		Extension:     "3gp",
		Resolution:    "",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          13,
		AudioBitrate:  0,
	},
	17: {
		Extension:     "3gp",
		Resolution:    "144p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          17,
		AudioBitrate:  24,
	},
	18: {
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          18,
		AudioBitrate:  96,
	},
	22: {
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          22,
		AudioBitrate:  192,
	},
	34: {
		Extension:     "flv",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          34,
		AudioBitrate:  128,
	},
	35: {
		Extension:     "flv",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          35,
		AudioBitrate:  128,
	},
	36: {
		Extension:     "3gp",
		Resolution:    "240p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          36,
		AudioBitrate:  36,
	},
	37: {
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          37,
		AudioBitrate:  192,
	},
	38: {
		Extension:     "mp4",
		Resolution:    "3072p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          38,
		AudioBitrate:  192,
	},
	43: {
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          43,
		AudioBitrate:  128,
	},
	44: {
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          44,
		AudioBitrate:  128,
	},
	45: {
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          45,
		AudioBitrate:  192,
	},
	46: {
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          46,
		AudioBitrate:  192,
	},
	82: {
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		Itag:          82,
		AudioBitrate:  96,
	},
	83: {
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          83,
		AudioBitrate:  96,
	},
	84: {
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          84,
		AudioBitrate:  192,
	},
	85: {
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          85,
		AudioBitrate:  192,
	},
	100: {
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          100,
		AudioBitrate:  128,
	},
	101: {
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          101,
		AudioBitrate:  192,
	},
	102: {
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          102,
		AudioBitrate:  192,
	},
	// DASH (video only)
	133: {
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          133,
		AudioBitrate:  0,
	},
	134: {
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          134,
		AudioBitrate:  0,
	},
	135: {
		Extension:     "mp4",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          135,
		AudioBitrate:  0,
	},
	136: {
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          136,
		AudioBitrate:  0,
	},
	137: {
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          137,
		AudioBitrate:  0,
	},
	138: {
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          138,
		AudioBitrate:  0,
	},
	160: {
		Extension:     "mp4",
		Resolution:    "144p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          160,
		AudioBitrate:  0,
	},
	242: {
		Extension:     "webm",
		Resolution:    "240p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          242,
		AudioBitrate:  0,
	},
	243: {
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          243,
		AudioBitrate:  0,
	},
	244: {
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          244,
		AudioBitrate:  0,
	},
	247: {
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          247,
		AudioBitrate:  0,
	},
	248: {
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          248,
		AudioBitrate:  9,
	},
	264: {
		Extension:     "mp4",
		Resolution:    "1440p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          264,
		AudioBitrate:  0,
	},
	266: {
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          266,
		AudioBitrate:  0,
	},
	271: {
		Extension:     "webm",
		Resolution:    "1440p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          271,
		AudioBitrate:  0,
	},
	272: {
		Extension:     "webm",
		Resolution:    "2160p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          272,
		AudioBitrate:  0,
	},
	278: {
		Extension:     "webm",
		Resolution:    "144p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          278,
		AudioBitrate:  0,
	},
	298: {
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          298,
		AudioBitrate:  0,
	},
	299: {
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          299,
		AudioBitrate:  0,
	},
	302: {
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          302,
		AudioBitrate:  0,
	},
	303: {
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          303,
		AudioBitrate:  0,
	},
	// DASH (audio only)
	139: {
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          139,
		AudioBitrate:  48,
	},
	140: {
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          140,
		AudioBitrate:  128,
	},
	141: {
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          141,
		AudioBitrate:  256,
	},
	171: {
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "vorbis",
		Itag:          171,
		AudioBitrate:  128,
	},
	172: {
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "vorbis",
		Itag:          172,
		AudioBitrate:  192,
	},
	249: {
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          249,
		AudioBitrate:  50,
	},
	250: {
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          250,
		AudioBitrate:  70,
	},
	251: {
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          251,
		AudioBitrate:  160,
	},
	// Live streaming
	92: {
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          92,
		AudioBitrate:  48,
	},
	93: {
		Extension:     "ts",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          93,
		AudioBitrate:  128,
	},
	94: {
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          94,
		AudioBitrate:  128,
	},
	95: {
		Extension:     "ts",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          95,
		AudioBitrate:  256,
	},
	96: {
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          96,
		AudioBitrate:  256,
	},
	120: {
		Extension:     "flv",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          120,
		AudioBitrate:  128,
	},
	127: {
		Extension:     "ts",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          127,
		AudioBitrate:  96,
	},
	128: {
		Extension:     "ts",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          128,
		AudioBitrate:  96,
	},
	132: {
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          132,
		AudioBitrate:  48,
	},
	151: {
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          151,
		AudioBitrate:  24,
	},
}

// var actionsObjRegexp = regexp.MustCompile(fmt.Sprintf(
// 	"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s|%s%s),?\\n?)+)\\};",
// 	jsvarStr, jsvarStr, reverseStr, jsvarStr, sliceStr, jsvarStr, spliceStr, jsvarStr, swapStr,
// ))

// var actionsFuncRegexp = regexp.MustCompile(fmt.Sprintf(
// 	"function(?: %s)?\\(a\\)\\{"+
// 		"a=a\\.split\\(\"\"\\);\\s*"+
// 		"((?:(?:a=)?%s\\.%s\\(a,\\d+\\);)+)"+
// 		"return a\\.join\\(\"\"\\)"+
// 		"\\}", jsvarStr, jsvarStr, jsvarStr,
// ))

// var reverseRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, reverseStr))
// var sliceRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, sliceStr))
// var spliceRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, spliceStr))
// var swapRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, swapStr))

// const (
// 	jsvarStr   = "[a-zA-Z_\\$][a-zA-Z_0-9]*"
// 	reverseStr = ":function\\(a\\)\\{" +
// 		"(?:return )?a\\.reverse\\(\\)" +
// 		"\\}"
// 	sliceStr = ":function\\(a,b\\)\\{" +
// 		"return a\\.slice\\(b\\)" +
// 		"\\}"
// 	spliceStr = ":function\\(a,b\\)\\{" +
// 		"a\\.splice\\(0,b\\)" +
// 		"\\}"
// 	swapStr = ":function\\(a,b\\)\\{" +
// 		"var c=a\\[0\\];a\\[0\\]=a\\[b(?:%a\\.length)?\\];a\\[b(?:%a\\.length)?\\]=c(?:;return a)?" +
// 		"\\}"
// )

// func reverseStringSlice(s []string) {
// 	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
// 		s[i], s[j] = s[j], s[i]
// 	}
// }

// type representation struct {
// 	Itag   int    `xml:"id,attr"`
// 	Height int    `xml:"height,attr"`
// 	URL    string `xml:"BaseURL"`
// }

// // GetDownloadURL gets the download url for a format
// func (info *VideoInfo) GetDownloadURL(format Format) (*url.URL, error) {
// 	return getDownloadURL(format, info.htmlPlayerFile)
// }

// func getDownloadURL(stream url.Values, htmlPlayerFile string) (string, error) {
// 	var signature string
// 	if s := stream.Get("s"); len(s) > 0 {
// 		tokens, err := getSigTokens(htmlPlayerFile)
// 		if err != nil {
// 			return "", err
// 		}
// 		signature = decipherTokens(tokens, s)
// 	} else {
// 		if sig := stream.Get("sig"); len(sig) > 0 {
// 			signature = sig
// 		}
// 	}
// 	var urlString string
// 	if s := stream.Get("url"); len(s) > 0 {
// 		urlString = s
// 	} else if s := stream.Get("stream"); len(s) > 0 {
// 		if c := stream.Get("conn"); len(c) > 0 {
// 			urlString = c
// 			if urlString[len(urlString)-1] != '/' {
// 				urlString += "/"
// 			}
// 		}
// 		urlString += s
// 	} else {
// 		return "", fmt.Errorf("couldn't extract url from format")
// 	}
// 	urlString, err := url.QueryUnescape(urlString)
// 	if err != nil {
// 		return "", err
// 	}
// 	u, err := url.Parse(urlString)
// 	if err != nil {
// 		return "", err
// 	}
// 	query := u.Query()
// 	query.Set("ratebypass", "yes")
// 	if len(signature) > 0 {
// 		if sp := stream.Get("sp"); sp != "" {
// 			query.Set(sp, signature)
// 		} else {
// 			query.Set("signature", signature)
// 		}
// 	}
// 	u.RawQuery = query.Encode()
// 	return u.String(), nil
// }
