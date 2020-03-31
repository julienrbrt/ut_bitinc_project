//Package util contains useful functions such as a joke generator
package util

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

//More information on http://api.apekool.nl/services/jokes/getjoke.html
const dutchJokeAPI = "http://api.apekool.nl/services/jokes/getjoke.php?type="

//More information on https://sv443.net/jokeapi/v2
const englishJokeAPI = "https://sv443.net/jokeapi/v2/joke/Miscellaneous,Dark?blacklistFlags=religious&type=single"

//apekool joke types
var dutchJokeType = []string{"alg", "be", "nl", "xxx"}

//GetJoke returns a random joke (Dutch and English)
func GetJoke(lang string) string {
	var noJoke string
	var jokeAPI string

	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	switch lang {
	case "NL":
		noJoke = "Geen grapjes voor deze week :("
		jokeAPI = dutchJokeAPI + dutchJokeType[rand.Intn(len(dutchJokeType))]
	default:
		noJoke = "No jokes for this week"
		jokeAPI = englishJokeAPI
	}

	//build API url with a certain type of joke
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
