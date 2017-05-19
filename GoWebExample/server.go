package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"net/http"
	"sync"
)

// RedditResponse is for Unmarshalling the response from subreddit
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// RedditPost is to represent the data we need from each post
type RedditPost struct {
	Subreddit   string  `json:"subreddit"`
	Created     float64 `json:"created"`
	CreatedUTC  time.Time
	Title       string `json:"title"`
	NumComments int    `json:"num_comments"`
}

// Slice of subreddits to get posts from
var subReddits = []string{"all", "golang", "videos", "pics", "gifs"}

func main() {

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // so we can read it..
	}))

	m.Get("/", func(r render.Render, x *http.Request) {
		r.HTML(200, "index", getRedditPosts())
	})

	m.RunOnAddr(":5050")
	m.Run()
}

// Get posts from the subreddit provided in argument, send it to channel ch
func getSubredditPosts(sr string, ch chan <- []RedditPost, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: 15 * time.Second}          // Please read: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	url := fmt.Sprintf("https://www.reddit.com/r/%s.json", sr) // Formatting URL with subreddit name

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Changing user agent, otherwise reddit denies the request
	req.Header.Set("User-Agent", "GoWebExample")

	// Make GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal the JSON to struct
	redditResponse := &RedditResponse{}
	if err := json.Unmarshal(body, &redditResponse); err != nil {
		log.Fatal(err)
	}
	// Get posts from the RedditResponse, we only need the RedditPost struct
	posts := []RedditPost{}
	for _, p := range redditResponse.Data.Children {
		p.Data.CreatedUTC = time.Unix(int64(p.Data.Created), 0) // Change UNIX timestamp to user readable form
		posts = append(posts, p.Data)                           // append post to be returned
	}
	ch <- posts
	log.Println(sr + " done")
}

// Function to get posts from the subReddits
func getRedditPosts() []RedditPost {
	var wg sync.WaitGroup
	wg.Add(len(subReddits)) // Add number of subreddits to WaitGroup
	data := []RedditPost{} // It will be returned as response, contains all the posts
	ch := make(chan []RedditPost) // Channel that will receive the posts from subreddit
	//  Looping over subReddits slice
	for _, subreddit := range subReddits {
		go getSubredditPosts(subreddit, ch, &wg) // Start a goroutine
	}
	go func(){
		// Wait till all the goroutines finish
		wg.Wait()
		// Close the channel or it will block
		 close(ch)
	}()
	// Loop over the channel and get data
	for posts := range ch{
		data = append(data, posts...) // Keep appending the posts to create final response
	}
	return data
}
