package disgopher

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ratelimitBucket struct {
	path        string
	mutex       *sync.Mutex
	ratelimited bool
	maxRetries  int
}

//HTTPSession ...
type HTTPSession struct {
	state            *clientState
	httpClient       http.Client
	Token            string
	ratelimitBuckets map[string]*ratelimitBucket
}

func (httpSession *HTTPSession) newRatelimitBucket(path string, maxRetries int) *ratelimitBucket {
	bucket := &ratelimitBucket{path: path, maxRetries: maxRetries, mutex: new(sync.Mutex)}
	httpSession.ratelimitBuckets[path] = bucket
	return bucket
}

func (httpSession *HTTPSession) request(req *http.Request, bucketPath string) (*http.Response, error) {
	bucket := httpSession.ratelimitBuckets[bucketPath]
	if bucket == nil {
		fmt.Println("bucket is nil... we will make a new one :^ )")
		bucket = httpSession.newRatelimitBucket(bucketPath, 5)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", httpSession.Token))
	req.Header.Add("X-Ratelimit-Precision", "millisecond")
	fmt.Println("Waiting on mutex", bucket.ratelimited)
	bucket.mutex.Lock()
	fmt.Println("Got Mutex", bucket.ratelimited)
	for try := 0; try < bucket.maxRetries; try++ {
		res, err := httpSession.httpClient.Do(req)
		if err != nil {
			panic(err)
		}
		remaining := res.Header.Values("X-Ratelimit-Remaining")
		if len(remaining) > 0 {
			fmt.Println(remaining[0])
			if remaining[0] == "0" {
				bucket.ratelimited = true
				var duration time.Duration
				resetAfter := res.Header.Values("X-Ratelimit-Reset-After")
				if len(resetAfter) > 0 {
					parsed, _ := strconv.ParseFloat(resetAfter[0], 64)
					fmt.Println(parsed)
					duration = time.Duration(parsed*1000) * time.Millisecond
				} else {
					reset := res.Header.Values("X-Ratelimit-Reset-After")
					if len(reset) > 0 {
						currentTime := fmt.Sprint(time.Now().UTC().Unix())
						parsedCurrTime, _ := strconv.ParseFloat(currentTime, 64)
						parsed, _ := strconv.ParseFloat(reset[0], 64)
						duration = time.Duration(parsedCurrTime*1000-parsed*1000) * time.Millisecond
					}
				}
				fmt.Printf("Sleeping for %s\n", duration)
				time.Sleep(duration)
				bucket.ratelimited = false
			}
		}
		//if res.StatusCode >= 200 && res.StatusCode < 300 {
		//	/	return res, err
		//	}
		bucket.mutex.Unlock()
		return res, err
	}
	return nil, nil
}

func (httpSession *HTTPSession) MessageCreate(channel string, content string) {
	x := fmt.Sprintf("https://discord.com/api/v7/channels/%s/messages", channel)
	fmt.Println(x)
	req, _ := http.NewRequest(
		"POST",
		x,
		strings.NewReader("{\"content\":\"hello friendo\",\"tts\":false}"))
	httpSession.request(req, "peepeepOOpOO")
}
