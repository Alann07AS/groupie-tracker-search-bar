package apimanagement

import (
	"encoding/json"
	"gt-alann/config"
	"io"
	"log"
	"net/http"
)

// init API________
var appConfig *config.Config
var api GroupiApi

func ConfigApi() {
	appConfig = config.ConfigLoad()
	getApi(appConfig.Api, &api)
}

type GroupiApi struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

// get Essential info
type ArtistsSimpleApi struct {
	Id         int    `json:"id"`
	Img        string `json:"image"`
	ArtistName string `json:"name"`
}

// get AllInfo for one artists
// stuct
type AllInfoArtist struct {
	Id           int      `json:"id"`
	Img          string   `json:"image"`
	ArtistName   string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relation     Relation //`json:"relations"`
}
type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func er(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getApi(apiUrl string, savePointer any) {
	apiRep, err := http.Get(apiUrl)
	er(err)
	data, err := io.ReadAll(apiRep.Body)
	er(err)
	err = json.Unmarshal(data, savePointer)
	er(err)
	err = apiRep.Body.Close()
	er(err)
}
