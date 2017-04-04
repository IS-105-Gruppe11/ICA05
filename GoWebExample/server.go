package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/google/gops/agent"

	"net/http"
)

// RedditResponse er for å "Unnmarshale" responsen fra subreddit apiene
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// RedditPost er for å representere dataene vi får fra hver post
type RedditPost struct {
	Subreddit   string  `json:"subreddit"`
	Created     float64 `json:"created"`
	CreatedUTC  time.Time
	Title       string `json:"title"`
	NumComments int    `json:"num_comments"`
}

// Slice med subreddits vi får post fra
var subReddits = []string{"RustPlay", "golang", "videos", "pics", "gifs"}

func main() {

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // so we can read it..
	}))

	m.Get("/", func(r render.Render, x *http.Request) {
		r.HTML(200, "index", getRedditPosts())
	})

	m.RunOnAddr(":8001")
	m.Run()

	if err := agent.Listen(&agent.Options{
		Addr: "158.39.77.203:8001",
	}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Hour)
}

//Få postene fra subreddit i et argument
func getSubredditPosts(sr string) []RedditPost {
	client := &http.Client{Timeout: 15 * time.Second}          // Les: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	url := fmt.Sprintf("https://www.reddit.com/r/%s.json", sr) //Formaterer URL med subreddit navn

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//Endrer user agent, ellers benekter reddit requesten
	req.Header.Set("User-Agent", "GoWebExample")

	// Gjør GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	// Leser response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal JSON til struct
	redditResponse := &RedditResponse{}
	if err := json.Unmarshal(body, &redditResponse); err != nil {
		log.Fatal(err)
	}
	// Får post fra RedditResponse, vi trenger bare RedditPost structen
	posts := []RedditPost{}
	for _, p := range redditResponse.Data.Children {
		p.Data.CreatedUTC = time.Unix(int64(p.Data.Created), 0) // Endrer UNIX timestamp til lesbar form
		posts = append(posts, p.Data)                           // Appender posten som skal returneres
	}
	return posts[0:1] // Får 1 post fra hver subreddit
}

// Funksjon for å få post fra subredditene
func getRedditPosts() []RedditPost {
	data := []RedditPost{}  // Vil bli returnert som en respons, inneholder all postene
	// Løkke som går over over subreddit slicene
	for _, subreddit := range subReddits {
		posts := getSubredditPosts(subreddit)
		data = append(data, posts...) // Slår sammen de returnerte postene med de andre
	}
	return data
}
