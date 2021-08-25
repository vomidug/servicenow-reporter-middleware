package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
    "context"
    "github.com/go-redis/redis/v8"
	"os"
	"time"
    "github.com/influxdata/influxdb-client-go/v2"
)

type Incident struct {
	Number string
	Short_description string
	Priority int
	Assigned_to string
	State string
}

type Config struct {
	Url string
	Login string
	Password string
	Period int
}

func compareSingle(first Incident, second Incident) bool {
	if (first.Number == second.Number) {
		return true
	} else {
		return false
	}
}

func chooseGreater(first int, second int) int {
	if (first > second){
		return first
	} else {
		return second
	}
}

func compareincs( oldincs []Incident, newincs []Incident ) ([]Incident, []Incident) {

	appeared := make([]Incident, chooseGreater( len(oldincs), len(newincs) ))
	disappeared := make([]Incident, chooseGreater( len(oldincs), len(newincs) ))
	for i := 0; i < len(oldincs); i++ {
		var flag bool = false;
		for j := 0; j < len(newincs); j++ {
			if ( compareSingle( oldincs[i], newincs[j] ) ) {
				flag = true
			}
		}
		if (!flag){
			disappeared = append(disappeared, oldincs[i])
			fmt.Println(disappeared)
		}
	}

	for i := 0; i < len(newincs); i++ {
		var flag bool = false;
		for j := 0; j < len(oldincs); j++ {
			if ( compareSingle( newincs[i], oldincs[j] ) ) {
				flag = true
			}
		}
		if (!flag){
			appeared = append(appeared, newincs[i])
		}
	}

	return appeared, disappeared
}

func main() {

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonstr,_ := ioutil.ReadAll(jsonFile)

	config := Config{}

	json.Unmarshal(jsonstr, &config)

	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	client := &http.Client{}

	oldjson, err := rdb.Get(ctx, "cache").Result()
    if err != nil { panic(err) }

	req, err := http.NewRequest("GET", config.Url, nil)
	req.Header.Add("Accept", `application/json`)
	req.SetBasicAuth(config.Login, config.Password)
	resp, err := client.Do(req)
	if err != nil{ log.Fatal(err) }
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	newjson := s[10:len(s)-1]

	err = rdb.Set(ctx, "cache", newjson, 0).Err()
    if err != nil { panic(err) }

	var oldincs []Incident	
	var newincs []Incident
	json.Unmarshal([]byte(oldjson), &oldincs)
	json.Unmarshal([]byte(newjson), &newincs)

	var appeared []Incident
	var disappeared []Incident

	appeared, disappeared = compareincs(oldincs, newincs)

	fmt.Println("Appeared:")
	fmt.Println(appeared)
	fmt.Println("Disppeared")
	fmt.Println(disappeared)

}
