package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dungtc/kafka-playground/youtube-stream/setting"
	// "google.golang.org/api/googleapi/transport"
	// youtube "google.golang.org/api/youtube/v3"
)

// const developerKey = "AIzaSyCzpC2mqJ-g5423Aas7xZDrwhg7inmkhYY"

type Snippet struct {
	PublishedAt string `json:"publishedAt"`
	ChannelId   string `json:"channelId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  string `json:"categoryId"`
}

type Video struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type VideoResponse struct {
	Kind  string  `json:"kind"`
	Etag  string  `json:"etag"`
	Items []Video `json:"items"`
}

type YoutubeRepository interface {
	GetVideos(categoryID string) (*VideoResponse, error)
}

type youtubeRepository struct {
	conf setting.GoogleConnection
}

func NewYoutubeRepository(conf setting.GoogleConnection) YoutubeRepository {
	return &youtubeRepository{
		conf: conf,
	}
}

func (y *youtubeRepository) GetVideos(categoryID string) (*VideoResponse, error) {
	resp, err := http.Get("https://www.googleapis.com/youtube/v3/videos?part=snippet&chart=mostPopular&videoCategoryId=" + categoryID + "&key=" + y.conf.ApiKey)
	if err != nil {
		fmt.Errorf("err %v", err)
	}
	defer resp.Body.Close()

	data := &VideoResponse{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	fmt.Printf("data %+v", data)

	return data, nil
}
