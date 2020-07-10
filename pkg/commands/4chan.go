package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/valyala/fasthttp"

	"sync-bot/pkg/storages"
)

type Thread struct {
	No            int    `json:"no"`
	Sticky        int    `json:"sticky,omitempty"`
	Closed        int    `json:"closed,omitempty"`
	Now           string `json:"now"`
	Name          string `json:"name"`
	Sub           string `json:"sub,omitempty"`
	Com           string `json:"com"`
	Filename      string `json:"filename"`
	Ext           string `json:"ext"`
	W             int    `json:"w"`
	H             int    `json:"h"`
	TnW           int    `json:"tn_w"`
	TnH           int    `json:"tn_h"`
	Tim           int64  `json:"tim"`
	Time          int    `json:"time"`
	Md5           string `json:"md5"`
	Fsize         int    `json:"fsize"`
	Resto         int    `json:"resto"`
	Capcode       string `json:"capcode,omitempty"`
	SemanticURL   string `json:"semantic_url"`
	Replies       int    `json:"replies"`
	Images        int    `json:"images"`
	OmittedPosts  int    `json:"omitted_posts,omitempty"`
	OmittedImages int    `json:"omitted_images,omitempty"`
	LastReplies   []struct {
		No       int    `json:"no"`
		Now      string `json:"now"`
		Name     string `json:"name"`
		Com      string `json:"com"`
		Filename string `json:"filename"`
		Ext      string `json:"ext"`
		W        int    `json:"w"`
		H        int    `json:"h"`
		TnW      int    `json:"tn_w"`
		TnH      int    `json:"tn_h"`
		Tim      int64  `json:"tim"`
		Time     int    `json:"time"`
		Md5      string `json:"md5"`
		Fsize    int    `json:"fsize"`
		Resto    int    `json:"resto"`
		Capcode  string `json:"capcode"`
	} `json:"last_replies"`
	LastModified int `json:"last_modified"`
	Bumplimit    int `json:"bumplimit,omitempty"`
	Imagelimit   int `json:"imagelimit,omitempty"`
}

type Page struct {
	Page    int      `json:"page"`
	Threads []Thread `json:"threads"`
}

type Post struct {
	No          int    `json:"no"`
	Now         string `json:"now"`
	Name        string `json:"name"`
	Sub         string `json:"sub,omitempty"`
	Com         string `json:"com,omitempty"`
	Filename    string `json:"filename"`
	Ext         string `json:"ext"`
	W           int    `json:"w,omitempty"`
	H           int    `json:"h,omitempty"`
	TnW         int    `json:"tn_w,omitempty"`
	TnH         int    `json:"tn_h,omitempty"`
	Tim         int64  `json:"tim,omitempty"`
	Time        int    `json:"time"`
	Md5         string `json:"md5,omitempty"`
	Fsize       int    `json:"fsize,omitempty"`
	Resto       int    `json:"resto"`
	Bumplimit   int    `json:"bumplimit,omitempty"`
	Imagelimit  int    `json:"imagelimit,omitempty"`
	SemanticURL string `json:"semantic_url,omitempty"`
	Replies     int    `json:"replies,omitempty"`
	Images      int    `json:"images,omitempty"`
	UniqueIps   int    `json:"unique_ips,omitempty"`
}

type ThreadResponse struct {
	Posts []Post `json:"posts"`
}

type PageResponse []Page

type fourchanCommand struct {
	cache *storages.CacheStorage
}

func NewFourchanCommand() CommandExecutor {
	return &fourchanCommand{
		cache: storages.NewCacheStorage(),
	}
}

func (c *fourchanCommand) GetMinRequiredRank() float64 {
	return 1
}

func (c *fourchanCommand) Validate(args []string) error {
	if len(args) == 0 || args[0] == "" {
		return errors.New("Укажите название доски.")
	}
	return nil
}

func (c *fourchanCommand) Exec(args []string) (*CommandResult, error) {
	err := c.Validate(args)
	if err != nil {
		return nil, err
	}
	boardName := args[0]
	var board *PageResponse
	cachedBoard := c.cache.Get(fmt.Sprintf("4chan_%s", boardName))
	if cachedBoard != nil {
		board = cachedBoard.(*PageResponse)
	} else {
		fourchan, err := fetchCatalog(boardName)
		if err != nil {
			return nil, err
		}
		board = fourchan
		c.cache.Set(fmt.Sprintf("4chan_%s", boardName), fourchan, time.Minute*5)
	}

	page := (*board)[rand.Intn(len(*board))]
	threadWithImages := make([]*Thread, 0)
	for _, thread := range page.Threads {
		if thread.Images > 0 {
			threadWithImages = append(threadWithImages, &thread)
		}
	}

	randomThread := threadWithImages[rand.Intn(len(threadWithImages))]
	thread, err := fetchThreadInfo(boardName, randomThread.No)
	if err != nil {
		return nil, err
	}
	images := make([]string, 0)
	for _, post := range thread.Posts {
		if post.Ext != "" {
			images = append(images, fmt.Sprintf("http://i.4cdn.org/%s/%d%s", boardName, post.Tim, post.Ext))
		}
	}
	randomPost := images[rand.Intn(len(images))]
	payloads := []*Event{
		{
			Method: "chatMsg",
			Message: EventPayload{
				Message: randomPost,
				Meta:    struct{}{},
			},
		},
	}
	return NewCommandResult(payloads), nil
}

func fetchThreadInfo(board string, number int) (*ThreadResponse, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("https://a.4cdn.org/%s/thread/%d.json", board, number))
	fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	thread := &ThreadResponse{}
	err := json.Unmarshal(bodyBytes, thread)
	return thread, err
}

func fetchCatalog(board string) (*PageResponse, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("https://a.4cdn.org/%s/catalog.json", board))

	fasthttp.Do(req, resp)

	bodyBytes := resp.Body()
	page := &PageResponse{}
	err := json.Unmarshal(bodyBytes, page)
	return page, err
}
