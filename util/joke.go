//Package util contains useful functions such as a joke generator
package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//More information on http://api.apekool.nl/services/jokes/getjoke.html
const dutchJokeAPI = "http://api.apekool.nl/services/jokes/getjoke.php?type="

//More information on https://sv443.net/jokeapi/v2
const englishJokeAPI = "https://sv443.net/jokeapi/v2/joke/Any?blacklistFlags=nsfw,religious,political,racist,sexist&type=single"

//More information on
const frenchJokeAPI = "https://blague.xyz/api/vdm/random"
const frenchJokeAPIToken = "_Ni2qnRfhubLAsW27nLjMsWzvJm_GO1yGsloGrim9RgcqmXilbo.wlK7vygNZ7mz"

//apekool joke types
var dutchJokeType = []string{"alg", "be", "nl"}

//GetJoke returns a random joke (Dutch and English)
func GetJoke(lang string) string {
	var noJoke string
	var jokeAPI string

	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	switch lang {
	case "NL":
		noJoke = "Geen grap vandaag :("
		//build API url with a certain type of joke
		jokeAPI = dutchJokeAPI + dutchJokeType[rand.Intn(len(dutchJokeType))]
	case "FR":
		return getJokeFR()
	default:
		noJoke = "No jokes for today :("
		jokeAPI = englishJokeAPI
	}

	resp, _ := http.Get(jokeAPI)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return noJoke
	}

	var joke map[string]interface{}
	if err := json.Unmarshal(body, &joke); err != nil {
		return noJoke
	}

	return joke["joke"].(string)
}

func getJokeFR() string {
	noJoke := "Pas de blague aujourd'hui :("

	req, _ := http.NewRequest("GET", frenchJokeAPI, nil)
	req.Header.Set("Authorization", frenchJokeAPIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return noJoke
	}

	var joke map[string]interface{}
	if err := json.Unmarshal(body, &joke); err != nil {
		return noJoke
	}

	return joke["vdm"].(map[string]interface{})["content"].(string)
}
