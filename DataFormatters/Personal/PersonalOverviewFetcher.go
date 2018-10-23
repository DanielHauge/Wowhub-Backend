package Personal

import (
	"../../Integrations/BlizzardOauthAPI"
	"../../Integrations/BlizzardOpenAPI"
	"../../Integrations/Raider.io"
	"../../Integrations/WarcraftLogs"
	"../../Integrations/Wowprogress"
	"../../Redis"
	log "../../Utility/Logrus"
	"github.com/avelino/slugify"
	"github.com/jinzhu/copier"
	"strconv"
	"sync"
)

func FetchFullPersonal(id int, profile *interface{}) error {

	var Profile PersonalProfile

	charMap, e := Redis.GetStruct("MAIN:" + strconv.Itoa(id))
	if e != nil {
		log.WithLocation().WithError(e).WithField("User", id).Error("There is no main registered to the user!")
		return e
	}
	char := BlizzardOauthAPI.CharacterMinimalFromMap(charMap)

	var wg sync.WaitGroup
	var blizzwait sync.WaitGroup
	blizzwait.Add(1)
	wg.Add(4)

	var blizzChar BlizzardOpenAPI.FullCharInfo
	go func() {
		blizzChar, e = BlizzardOpenAPI.GetBlizzardChar(char.Realm, char.Name, char.Region)
		go Redis.Set("GUILD:"+strconv.Itoa(id), blizzChar.Guild.Name+":"+blizzChar.Guild.Realm+":"+char.Region)
		Profile.Character = blizzChar
		wg.Done()
		blizzwait.Done()
	}()

	var raiderio Raider_io.CharacterProfile
	go func() {
		raiderio, e = Raider_io.GetRaiderIORank(Raider_io.CharInput{Name: char.Name, Realm: char.Realm, Region: char.Region})
		Profile.RaiderIOProfile = raiderio
		wg.Done()
	}()

	var logs WarcraftLogs.Encounters
	go func() {
		logs, e = WarcraftLogs.GetWarcraftLogsRanks(WarcraftLogs.CharInput{Name: char.Name, Realm: char.Realm, Region: char.Region})
		Profile.WarcraftLogsRanks = logs.Encounters
		wg.Done()
	}()

	var wowprog Wowprogress.GuildRank
	go func() {
		blizzwait.Wait()
		wowprog, e = Wowprogress.GetGuildRank(Wowprogress.Input{Region: char.Region, Realm: slugify.Slugify(char.Realm), Guild: blizzChar.Guild.Name})
		Profile.GuildRank = wowprog
		wg.Done()
	}()

	wg.Wait()

	copier.Copy(profile, Profile)

	return e
}

func FetchRaiderioPersonal(id int, Profile *interface{}) error {

	charMap, e := Redis.GetStruct("MAIN:" + strconv.Itoa(id))
	if e != nil {
		log.WithLocation().WithError(e).WithField("User", id).Error("There is no main registered to the user!")
		return e
	}
	char := BlizzardOauthAPI.CharacterMinimalFromMap(charMap)
	prof, e := Raider_io.GetRaiderIORank(Raider_io.CharInput{Name: char.Name, Realm: char.Realm, Region: char.Region})
	copier.Copy(Profile, &prof)
	return e
}

func FetchWarcraftlogsPersonal(id int, Logs *interface{}) error {
	charMap, e := Redis.GetStruct("MAIN:" + strconv.Itoa(id))
	if e != nil {
		log.WithLocation().WithError(e).WithField("User", id).Error("There is no main registered to the user!")
		return e
	}
	char := BlizzardOauthAPI.CharacterMinimalFromMap(charMap)

	logs, e := WarcraftLogs.GetWarcraftLogsRanks(WarcraftLogs.CharInput{Name: char.Name, Realm: char.Realm, Region: char.Region})
	copier.Copy(Logs, &logs)
	return e
}

func FetchBlizzardPersonal(id int, Profile *interface{}) error {

	charMap, e := Redis.GetStruct("MAIN:" + strconv.Itoa(id))
	if e != nil {
		log.WithLocation().WithError(e).WithField("User", id).Error("There is no main registered to the user!")
		return e
	}
	char := BlizzardOauthAPI.CharacterMinimalFromMap(charMap)
	blizzChar, e := BlizzardOpenAPI.GetBlizzardChar(char.Realm, char.Name, char.Region)
	go Redis.Set("GUILD:"+strconv.Itoa(id), blizzChar.Guild.Name+":"+blizzChar.Guild.Realm+":"+char.Region)
	copier.Copy(Profile, &blizzChar)
	return e
}
