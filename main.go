package main

import (
	"bytes"
	"fmt"
	"net/http"
	"crypto/tls"
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
	Assigned_to SnowUser
	State string
}

type SnowUser struct {
	Link string
	Value string
}

type User struct {
	Value string
	Name string
}

type Config struct {
	Url string
	Login string
	Password string
	Period int
	InfluxURL string
	Redisurl string
	Redispassword string
	Redisdb int
	Influxdbname string
	Influxmeasurementname string
	Influxloginpasswordtoken string
	Influxorgname string
	Telegramboturl string
	Chatid string
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

/* will be decomposed later i swear a god
func init(params) int {
	
}
*/

func run(config Config) {

	usersMap := make(map[string]string)
	
	jsonUsers, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonUsers.Close()
	jsonstr,_ := ioutil.ReadAll(jsonUsers)
	json.Unmarshal(jsonstr, &usersMap)

	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
        Addr: config.Redisurl,
        Password: config.Redispassword,
        DB:       config.Redisdb,
    })

	client := &http.Client{
		    Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	oldjson, err := rdb.Get(ctx, "cache").Result()
    if err != nil { fmt.Println(err.Error()) }

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

	counts := make(map[string]int)
	counts["unassigned"] = 0
	
	for e := range usersMap {
		counts[usersMap[e]] = 0
	}

	for i := 0; i < len(newincs); i++ {
		if newincs[i].Assigned_to == (SnowUser{}){
			counts["unassigned"]++
		} else {
			counts[usersMap[newincs[i].Assigned_to.Value]]++
		}
	}

	// TODO: find out how to replace this "mydb" to variable from config
	influx := influxdb2.NewClient(config.InfluxURL, config.Influxloginpasswordtoken)
	writeAPI := influx.WriteAPIBlocking(config.Influxorgname, "mydb" )
	p := influxdb2.NewPointWithMeasurement("tickets")

	for e := range counts {
		p = p.AddField(e, counts[e])
	}
	p = p.AddField("number", len(newincs))
	p = p.SetTime(time.Now())

	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Printf("Write error: %s\n", err.Error());
	}

	var appeared []Incident
	var disappeared []Incident

	appeared, disappeared = compareincs(oldincs, newincs)

	appFlag := false;
	disappFlag := false;

	for i := 0; i < len(appeared); i++ {
		if appeared[i] != (Incident{}) {
			appFlag = true
			break
		}
	}

	for i := 0; i < len(disappeared); i++ {
		if disappeared[i]!=(Incident{}) {
			disappFlag = true;
			break
		}
	}

	var resultMsg string

	if appFlag {
		for i := 0; i < len(appeared); i++ {
			if appeared[i]!=(Incident{}) {
				resultMsg = resultMsg + "+ " + appeared[i].Number + "\n"
			}	
		}
	}

	if disappFlag {
		for i := 0; i < len(disappeared); i++ {
			if disappeared[i] != (Incident{}) {
				resultMsg = resultMsg + "- " + disappeared[i].Number + "\n"
			}
		}	
	}

	if resultMsg != "" {
		
		var jsonMsg = []byte(`{"text":"` + resultMsg + `","chat_id":"` + config.Chatid + `"}`)
		req, err := http.NewRequest("POST", config.Telegramboturl, bytes.NewBuffer(jsonMsg))
	    req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)

	}

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

	for {
		run(config)
		time.Sleep( time.Duration(config.Period) * time.Second )
	}
	
}
