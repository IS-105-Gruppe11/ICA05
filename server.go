package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"

)


type liga struct {
	Name string `json:"name"`
	Clubs []struct {
		Key string `json:"key"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"clubs"`
}

type steamPlayers struct {
	Response struct {
		PlayerCount int `json:"player_count"`
		Result int `json:"result"`
	} `json:"response"`
}

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Get("/", func(r render.Render, x *http.Request) {
		r.HTML(200, "index", getLiga())
	})
	m.RunOnAddr(":8001")
	m.Run()
}



// Get clubs in liga
func getLiga() *liga {
	readApi, err := http.Get("https://raw.githubusercontent.com/openfootball/" +
		"football.json/master/2016-17/es.1.clubs.json")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(readApi.Body)
	if err != nil {
		log.Fatal(err)
	}
	clubs := &liga{}
	if err := json.Unmarshal(bytes, &clubs); err != nil {
		log.Fatal(err)

	}
	return clubs

}
