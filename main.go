package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/* Type definitions */

// To parse and unparse this JSON data, add this code to your project and do:
//
//    topLevel, err := UnmarshalTopLevel(bytes)
//    bytes, err = topLevel.Marshal()

func UnmarshalTopLevel(data []byte) (TopLevel, error) {
	var r TopLevel
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TopLevel) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TopLevel struct {
	Recenttracks Recenttracks `json:"recenttracks"`
}

type Recenttracks struct {
	Attr  RecenttracksAttr `json:"@attr"`
	Track []Track          `json:"track"`
}

type RecenttracksAttr struct {
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	User       string `json:"user"`
	Total      string `json:"total"`
	TotalPages string `json:"totalPages"`
}

type Track struct {
	Artist     Album      `json:"artist"`
	Attr       *TrackAttr `json:"@attr,omitempty"`
	Mbid       string     `json:"mbid"`
	Album      Album      `json:"album"`
	Streamable string     `json:"streamable"`
	URL        string     `json:"url"`
	Name       string     `json:"name"`
	Image      []Image    `json:"image"`
	Date       *Date      `json:"date,omitempty"`
}

type Album struct {
	Mbid string `json:"mbid"`
	Text string `json:"#text"`
}

type TrackAttr struct {
	Nowplaying string `json:"nowplaying"`
}

type Date struct {
	Uts  string `json:"uts"`
	Text string `json:"#text"`
}

type Image struct {
	Size Size   `json:"size"`
	Text string `json:"#text"`
}

type Size string

const (
	Extralarge Size = "extralarge"
	Large      Size = "large"
	Medium     Size = "medium"
	Small      Size = "small"
)

/* Type Scrobbles */

type Scrobbles []Scrobble

func UnmarshalScrobbles(data []byte) (Scrobbles, error) {
	var r Scrobbles
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Scrobbles) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Scrobble struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
	Album  string `json:"album"`
	When   When   `json:"when"`
}

type When struct {
	POSIX int64  `json:"posix"`
	Human string `json:"human"`
}

/* End Types*/

func fetchScrobbles(apiKey string, username string) []byte {
	var apiEndpoint = "http://ws.audioscrobbler.com/2.0/"

	httpClient := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Golang; +https://github.com/ricardobalk/go-fetch-scrobbles)")
	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", apiKey)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func parsedRawResponse(input []byte) string {
	scrobbles := TopLevel{}
	jsonErr := json.Unmarshal(input, &scrobbles)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	bytes, err := json.Marshal(scrobbles.Recenttracks)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes[:])
}

func buildList(body []byte) string {
	scrobbles := TopLevel{}
	var output string

	jsonErr := json.Unmarshal(body, &scrobbles)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	tracks := scrobbles.Recenttracks.Track
	for i := range tracks {
		artist := tracks[i].Artist.Text
		song := tracks[i].Name
		album := tracks[i].Album.Text
		var date int64

		if tracks[i].Date != nil {
			val, err := strconv.ParseInt(tracks[i].Date.Uts, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			date = val
		} else {
			if tracks[i].Attr != nil {
				_, err := strconv.ParseBool(tracks[i].Attr.Nowplaying)
				if err != nil {
					log.Fatal(err)
				}
				date = time.Now().Unix()
			}
		}
		parsedDate := time.Unix(date, 0).UTC()

		output = output + fmt.Sprintf("%02d: %s - %s [%s] at %s\n", i+1, artist, song, album, parsedDate.String())
	}
	return output
}

func buildJSONResponse(input []byte) string {
	scrobbles := TopLevel{}
	output := Scrobbles{}

	jsonErr := json.Unmarshal(input, &scrobbles)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	tracks := scrobbles.Recenttracks.Track
	for i := range tracks {
		artist := tracks[i].Artist.Text
		song := tracks[i].Name
		album := tracks[i].Album.Text
		var date int64

		if tracks[i].Date != nil {
			val, err := strconv.ParseInt(tracks[i].Date.Uts, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			date = val
		} else {
			if tracks[i].Attr != nil {
				_, err := strconv.ParseBool(tracks[i].Attr.Nowplaying)
				if err != nil {
					log.Fatal(err)
				}
				date = time.Now().Unix()
			}
		}
		parsedDate := time.Unix(date, 0).UTC()
		data := Scrobble{artist, song, album, When{date, parsedDate.String()}}
		output = append(output, data)
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes[:])
}

func formatScrobbles(input []byte, format string) string {
	switch format {
	case "raw":
		return parsedRawResponse(input)
	case "list":
		return buildList(input)
	case "json":
		return buildJSONResponse(input)
	default:
		flag.Usage()
		log.Fatal("Invalid output format.")
	}
	return ""
}

func main() {
	apiTokenPtr := flag.String("api-token", "", "Last.fm API Token")
	usernamePtr := flag.String("username", "Batmaniosaurus", "Username")
	formatPtr := flag.String("format", "list", "Format of output, currently allowed values are default, list or raw.\nDefault output means processed JSON, list shows a list of songs and raw shows response data from Last.fm.")
	flag.Parse()

	if *apiTokenPtr == "" {
		flag.Usage()
		log.Fatal("Missing Last.fm API token.")
	}

	fetchScrobbles(*apiTokenPtr, *usernamePtr)
}