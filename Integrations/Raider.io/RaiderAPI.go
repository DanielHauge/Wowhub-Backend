package Raider_io

import (
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var json = jsoniter.ConfigFastest

func GetRaiderIORank(input CharInput) (CharacterProfile, error){
	log.Info("Fetching RaiderIO profile for: ",input)
	url := "https://raider.io/api/v1/characters/profile?region="+input.Region+"&realm="+input.Realm+"&name="+input.Name+"&fields=mythoc_plus_scores%2Cmythic_plus_ranks%2Cmythic_plus_recent_runs%2Cmythic_plus_highest_level_runs%2Cmythic_plus_weekly_highest_level_runs%2C"

	resp, e := http.Get(url)
	if e != nil{
		log.Error(e, " -> Something went wrong in getting data from RaiderIO")
		return CharacterProfile{}, e
	}
	defer resp.Body.Close()

	var rankings CharacterProfile
	e = json.NewDecoder(resp.Body).Decode(&rankings)
	if e != nil { log.Error(e, "Something went wrong in decoding data from RaiderIO") }

	return rankings, e
}
