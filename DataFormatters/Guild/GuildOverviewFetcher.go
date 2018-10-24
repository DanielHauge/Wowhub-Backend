package Guild

import (
	"../../Integrations/BlizzardOpenAPI"
	"../../Integrations/Raider.io"
	"../../Integrations/WarcraftLogs"
	"../../Redis"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
	"time"
)

func FetchFullGuildOverview(id int, result *interface{}) error {

	// TODO: Error handling??? Logging?
	var FullGuildOverview Overview
	guildstring := Redis.Get("GUILD:" + strconv.Itoa(id))
	split := strings.Split(guildstring, ":")
	Guild := struct {
		Name   string
		Realm  string
		Region string
	}{Name: split[0], Realm: split[1], Region: split[2]}

	FullGuildOverview.Name = Guild.Name
	FullGuildOverview.Realm = Guild.Realm
	Progression, e := Raider_io.GetRaiderIOGuild(Guild.Region, Guild.Realm, Guild.Name)

	FullGuildOverview.Progress = Progression.RaidProgression
	t1 := time.Now().Add(-time.Hour * 24 * 7)
	t2 := time.Now().Add(time.Hour * 24 * 7)

	// TODO: Create struct that are more helpfull. ie. Link maybe or something.
	WarcraftlogsReports, e := WarcraftLogs.GetWarcraftGuildReports(Guild.Name, Guild.Realm, Guild.Region, t1.Unix(), t2.Unix())
	FullGuildOverview.WarcraftlogReports = WarcraftlogsReports

	Roster, e := BlizzardOpenAPI.GetBlizzardGuildMembers(Guild.Name, Guild.Region, Guild.Realm)

	var guildmembers []Member
	for _, v := range Roster.Members {
		if v.Rank < 5 {
			guildmembers = append(guildmembers, Member{Name: v.Character.Name, Rank: v.Rank, Role: v.Character.Spec.Role, Class: v.Character.Class})
		}
	}

	FullGuildOverview.Roster = guildmembers

	//TODO: Create postgress and have real data here:
	FullGuildOverview.LastRaid = strconv.FormatInt(time.Now().Unix(), 10)
	FullGuildOverview.NextRaid = strconv.FormatInt(time.Now().Add(time.Hour*24*3).Unix(), 10)
	FullGuildOverview.RaidDays = []RaidNight{{DayOfTheWeek: 2, StartTime: "19:00", EndTime: "23:00"}, {DayOfTheWeek: 6, StartTime: "20:00", EndTime: "22:00"}}
	copier.Copy(result, FullGuildOverview)
	return e

}
