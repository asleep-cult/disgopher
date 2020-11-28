package disgopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const baseURL = "https://discord.com/api/v7"

//HTTPError ...
type HTTPError struct {
	Message        string `json:"message"`
	HTTPStatusCode int
	Code           int `json:"code"`
}

//Error ...
func (err *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError [HTTP Status Code: %s, Code: %s]: %s", fmt.Sprint(err.HTTPStatusCode), fmt.Sprint(err.Code), err.Message)
}

//HTTPSession ...
type HTTPSession struct {
	state                *clientState
	httpClient           http.Client
	Token                string
	ratelimitBuckets     map[string]*ratelimitBucket
	globallyRatelimited  bool
	globalRatelimitMutex *sync.Mutex
}

type ratelimitBucket struct {
	path        string
	mutex       *sync.Mutex
	ratelimited bool
	maxRetries  int
}

type ratelimitResponse struct {
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
}

func (httpSession *HTTPSession) newRatelimitBucket(path string, maxRetries int) *ratelimitBucket {
	bucket := &ratelimitBucket{path: path, maxRetries: maxRetries, mutex: new(sync.Mutex)}
	httpSession.ratelimitBuckets[path] = bucket
	return bucket
}

func (httpSession *HTTPSession) request(req *http.Request, bucketPath string) (*http.Response, []byte, error) {
	bucket := httpSession.ratelimitBuckets[bucketPath]
	if bucket == nil {
		bucket = httpSession.newRatelimitBucket(bucketPath, 5)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", httpSession.Token))
	req.Header.Add("X-Ratelimit-Precision", "millisecond")
	if httpSession.globallyRatelimited {
		fmt.Println("Global Ratelimit!")
		httpSession.globalRatelimitMutex.Lock()
		httpSession.globalRatelimitMutex.Unlock()
	}
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	for try := 0; try < bucket.maxRetries; try++ {
		resp, err := httpSession.httpClient.Do(req)
		if err != nil {
			return resp, make([]byte, 1), err
		}
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == 429 {
			ratelimit := new(ratelimitResponse)
			json.Unmarshal(data, ratelimit)
			bucket.ratelimited = true
			if ratelimit.Global {
				httpSession.globallyRatelimited = true
				httpSession.globalRatelimitMutex.Lock()
			}
			time.Sleep(time.Duration(ratelimit.RetryAfter/1000) * time.Second)
			bucket.ratelimited = false
			if ratelimit.Global {
				httpSession.globallyRatelimited = false
				httpSession.globalRatelimitMutex.Unlock()
			}
		} else {
			remaining := resp.Header.Values("X-Ratelimit-Remaining")
			if len(remaining) > 0 {
				if remaining[0] == "0" {
					bucket.ratelimited = true
					var duration time.Duration
					resetAfter := resp.Header.Values("X-Ratelimit-Reset-After")
					if len(resetAfter) > 0 {
						parsed, _ := strconv.ParseFloat(resetAfter[0], 64)
						duration = time.Duration(parsed*1000) * time.Millisecond
					} else {
						reset := resp.Header.Values("X-Ratelimit-Reset-After")
						if len(reset) > 0 {
							currentTime := fmt.Sprint(time.Now().UTC().Unix())
							parsedCurrTime, _ := strconv.ParseFloat(currentTime, 64)
							parsed, _ := strconv.ParseFloat(reset[0], 64)
							duration = time.Duration(parsedCurrTime*1000-parsed*1000) * time.Millisecond
						}
					}
					time.Sleep(duration)
					bucket.ratelimited = false
				}
			}
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, data, nil
		}
		actualerr := &HTTPError{HTTPStatusCode: resp.StatusCode}
		json.Unmarshal(data, actualerr)
		return resp, data, actualerr
	}
	return nil, make([]byte, 1), nil
}

func (httpSession *HTTPSession) getChannel(channelID string) ([]byte, error) {
	path := fmt.Sprintf("%s/cahnnels/%s", baseURL, channelID)
	bucketPath := fmt.Sprintf("GET-channels/%s", channelID)
	httpreq, _ := http.NewRequest("GET", path, nil)
	_, data, err := httpSession.request(httpreq, bucketPath)
	return data, err
}

func (httpSession *HTTPSession) modifyChannel(channelID string, req interface{}) ([]byte, error) {
	path := fmt.Sprintf("%s/channels/%s", baseURL, channelID)
	bucketPath := fmt.Sprintf("PATCH-channels/%s", channelID)
	data, _ := json.Marshal(req)
	httpreq, _ := http.NewRequest(
		"PATCH",
		path,
		bytes.NewBuffer(data))
	_, respdata, err := httpSession.request(httpreq, bucketPath)
	return respdata, err
}

func (httpSession *HTTPSession) deleteChannel(channelID string) ([]byte, error) {
	path := fmt.Sprintf("%s/channels/%s", baseURL, channelID)
	bucketPath := "DELETE-channels" //I'm 99% sure this has no major params??
	httpreq, _ := http.NewRequest(
		"DELETE",
		path,
		nil)
	_, respdata, err := httpSession.request(httpreq, bucketPath)
	return respdata, err
}

//getMessage

//MessageCreateRequest ...
type MessageCreateRequest struct {
	Content string      `json:"content"`
	Nonce   interface{} `json:"nonce"`
	TTS     bool        `json:"tts"`
	//File
	Embed *Embed `json:"embed,omitempty"`
	//PayloadJSON
	//AllowedMentions
	//MessageReference
}

func (httpSession *HTTPSession) messageCreate(channelID string, req *MessageCreateRequest) ([]byte, error) {
	path := fmt.Sprintf("%s/channels/%s/messages", baseURL, channelID)
	bucketPath := fmt.Sprintf("POST-%s/messages", channelID)
	data, _ := json.Marshal(req)
	httpreq, _ := http.NewRequest(
		"POST",
		path,
		bytes.NewBuffer(data))
	_, respdata, err := httpSession.request(httpreq, bucketPath)
	return respdata, err
}

//crossPostMessage

//createReaction

//deleteReaction

//getReactions

//deleteAllReactions

//deleteAllReactionsForEmoji

//editMessage (Needs to be treated like modifyChannel)

//more endpoints i dont feel like dealing with
