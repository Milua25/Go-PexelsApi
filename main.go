package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2X"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

// NewClient Function
func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

// SearchPhotos Method
func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json")
	}
	return &result, err
}

// CuratedPhotos Method
func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?&per_page=%d&page=%d", perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var curated CuratedResult
	err = json.Unmarshal(data, &curated)
	if err != nil {
		return nil, err
	}
	return &curated, nil
}

// GetPhoto Method
func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getPhoto Photo

	err = json.Unmarshal(data, &getPhoto)
	if err != nil {
		return nil, err
	}

	return &getPhoto, nil
}

// GetRandomPhoto Method
func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}

// requestDoWithAuth Method
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)

	response, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	times, err := strconv.Atoi(response.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return nil, err
	} else {
		c.RemainingTimes = int32(times)
	}
	return response, nil
}

type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}

type Video struct {
	Id            int32           `json:"id"`
	Width         int32           `json:"width"`
	Height        int32           `json:"height"`
	Url           string          `json:"url"`
	Image         string          `json:"image"`
	FullRes       interface{}     `json:"full_res"`
	Duration      float64         `json:"duration"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type VideoFiles struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPictures struct {
	Id      int32  `json:"id"`
	Picture string `json:"picture"`
	Nr      int32  `json:"nr"`
}

type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

func (c *Client) SearchVideo(query string, perPage, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result VideoSearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) PopularVideo(perPage, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PopularVideos
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.PopularVideo(1, randNum)
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
}

func (c *Client) GetRemainingRequestsInThisMonth() int32 {
	return c.RemainingTimes
}
func main() {
	err := os.Setenv("PexelsToken", "xxx")

	if err != nil {
		fmt.Errorf("could not set environment variable:%v", err)
	}

	var TOKEN = os.Getenv("PexelsToken")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves", 15, 1)

	if err != nil {
		fmt.Errorf("search error:%v", err)
	}
	if result.Page == 0 {
		fmt.Errorf("search result wrong:%v", err)
	}

	// print out result
	fmt.Println(result)
}
