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

//More information on the Joke API here: http://api.apekool.nl/services/jokes/getjoke.html
const jokeAPI = "http://api.apekool.nl/services/jokes/getjoke.php?type="

var jokeType = []string{"alg", "be", "nl", "xxx"}

//GetJoke will return a random Dutch joke
func GetJoke() string {
	var noJoke = "Geen grapjes vandaag :("

	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	//build API url with a certain type of joke
	resp, err := http.Get(jokeAPI + jokeType[rand.Intn(len(jokeType))])
	if err != nil {
		log.Fatalln(err)
		return noJoke
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return noJoke
	}

	var joke map[string]interface{}
	if err := json.Unmarshal(body, &joke); err != nil {
		log.Fatalln(err)
		return noJoke
	}

	return joke["joke"].(string)
}
