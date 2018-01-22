package sweetiebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/blackhole12/discordgo"
)

type ModuleID string
type CommandID string

// BotConfig lists all bot configuration options, grouped into structs
type BotConfig struct {
	Version     int  `json:"version"`
	LastVersion int  `json:"lastversion"`
	SetupDone   bool `json:"setupdone"`
	Basic       struct {
		IgnoreInvalidCommands bool                    `json:"ignoreinvalidcommands"`
		Importable            bool                    `json:"importable"`
		ModRole               DiscordRole             `json:"modrole"`
		ModChannel            DiscordChannel          `json:"modchannel"`
		FreeChannels          map[DiscordChannel]bool `json:"freechannels"`
		BotChannel            DiscordChannel          `json:"botchannel"`
		Aliases               map[string]string       `json:"aliases"`
		ListenToBots          bool                    `json:"listentobots"`
		CommandPrefix         string                  `json:"commandprefix"`
		SilenceRole           DiscordRole             `json:"silencerole"`
	} `json:"basic"`
	Modules struct {
		Channels           map[ModuleID]map[DiscordChannel]bool  `json:"modulechannels"`
		Disabled           map[ModuleID]bool                     `json:"moduledisabled"`
		CommandRoles       map[CommandID]map[DiscordRole]bool    `json:"commandroles"`
		CommandChannels    map[CommandID]map[DiscordChannel]bool `json:"commandchannels"`
		CommandLimits      map[CommandID]int64                   `json:"Commandlimits"`
		CommandDisabled    map[CommandID]bool                    `json:"commanddisabled"`
		CommandPerDuration int                                   `json:"commandperduration"`
		CommandMaxDuration int64                                 `json:"commandmaxduration"`
	} `json:"modules"`
	Spam struct {
		ImagePressure      float32                    `json:"imagepressure"`
		PingPressure       float32                    `json:"pingpressure"`
		LengthPressure     float32                    `json:"lengthpressure"`
		RepeatPressure     float32                    `json:"repeatpressure"`
		LinePressure       float32                    `json:"linepressure"`
		BasePressure       float32                    `json:"basepressure"`
		PressureDecay      float32                    `json:"pressuredecay"`
		MaxPressure        float32                    `json:"maxpressure"`
		MaxChannelPressure map[DiscordChannel]float32 `json:"maxchannelpressure"`
		MaxRemoveLookback  int                        `json:"MaxSpamRemoveLookback"`
		IgnoreRole         DiscordRole                `json:"ignorerole"`
		RaidTime           int64                      `json:"maxraidtime"`
		RaidSize           int                        `json:"raidsize"`
		AutoSilence        int                        `json:"autosilence"`
		LockdownDuration   int                        `json:"lockdownduration"`
	} `json:"spam"`
	Users struct {
		TimezoneLocation string               `json:"timezonelocation"`
		WelcomeChannel   DiscordChannel       `json:"welcomechannel"`
		WelcomeMessage   string               `json:"welcomemessage"`
		SilenceMessage   string               `json:"silencemessage"`
		Roles            map[DiscordRole]bool `json:"userroles"`
		NotifyChannel    DiscordChannel       `json:"joinchannel"`
		TrackUserLeft    bool                 `json:"trackuserleft"`
	} `json:"users"`
	Bucket struct {
		MaxItems       int             `json:"maxbucket"`
		MaxItemLength  int             `json:"maxbucketlength"`
		MaxFightHP     int             `json:"maxfighthp"`
		MaxFightDamage int             `json:"maxfightdamage"`
		Items          map[string]bool `json:"items"`
	} `json:"bucket"`
	Markov struct {
		MaxPMlines     int  `json:"maxpmlines"`
		MaxLines       int  `json:"maxquotelines"`
		DefaultLines   int  `json:"defaultmarkovlines"`
		UseMemberNames bool `json:"usemembernames"`
	} `json:"markov"`
	Filter struct {
		Filters   map[string]map[string]bool         `json:"filters"`
		Channels  map[string]map[DiscordChannel]bool `json:"channels"`
		Responses map[string]string                  `json:"responses"`
		Templates map[string]string                  `json:"templates"`
	} `json:"filter"`
	Bored struct {
		Cooldown int64           `json:"maxbored"`
		Commands map[string]bool `json:"boredcommands"`
	}
	Information struct {
		Rules             map[int]string `json:"rules"`
		HideNegativeRules bool           `json:"hidenegativerules"`
	} `json:"help"`
	Log struct {
		Cooldown int64          `json:"maxerror"`
		Channel  DiscordChannel `json:"logchannel"`
	} `json:"log"`
	Witty struct {
		Responses map[string]string `json:"witty"`
		Cooldown  int64             `json:"maxwit"`
	} `json:"Wit"`
	Scheduler struct {
		BirthdayRole DiscordRole `json:"birthdayrole"`
	} `json:"scheduler"`
	Miscellaneous struct {
		MaxSearchResults int `json:"maxsearchresults"`
	} `json:"misc"`
	Status struct {
		Cooldown int             `json:"statusdelaytime"`
		Lines    map[string]bool `json:"lines"`
	} `json:"status"`
	Quote struct {
		Quotes map[DiscordUser][]string `json:"quotes"`
	} `json:"quote"`
}

// ConfigHelp is a map of help strings for the configuration options above
var ConfigHelp = map[string]map[string]string{
	"basic": map[string]string{
		"ignoreinvalidcommands": "If true, the bot won't display an error if a nonsensical command is used. This helps reduce confusion with other bots that also use the `!` prefix.",
		"importable":            "If true, the collections on this server will be importable into another server.",
		"modrole":               "This is intended to point at a moderator role shared by all admins and moderators of the server for notification purposes.",
		"modchannel":            "This should point at the hidden moderator channel, or whatever channel moderates want to be notified on.",
		"freechannels":          "This is a list of all channels that are exempt from rate limiting. Usually set to the dedicated `#botabuse` channel in a server.",
		"botchannel":            "This allows you to designate a particular channel to point users if they are trying to run too many commands at once. Usually this channel will also be included in `basic.freechannels`",
		"aliases":               "Can be used to redirect commands, such as making `!listgroup` call the `!listgroups` command. Useful for making shortcuts.\n\nExample: `!setconfig basic.aliases kawaii \"pick cute\"` sets an alias mapping `!kawaii arg1...` to `!pick cute arg1...`, preserving all arguments that are passed to the alias.",
		"listentobots":          "If true, processes messages from other bots and allows them to run commands. Bots can never trigger anti-spam. Defaults to false.",
		"commandprefix":         "Determines the SINGLE ASCII CHARACTER prefix used to denote bot commands. You can't set it to an emoji or any weird foreign character. The default is `!`. If this is set to an invalid value, it defaults to `!`.",
		"silencerole":           "This should be a role with no permissions, so the bot can quarantine potential spammers without banning them.",
	},
	"modules": map[string]string{
		"commandroles":       "A map of which roles are allowed to run which command. If no mapping exists, everyone can run the command.",
		"commandchannels":    "A map of which channels commands are allowed to run on. No entry means a command can be run anywhere. If \"!\" is included as a channel, it switches from a whitelist to a blacklist, enabling you to exclude certain channels instead of allow certain channels.",
		"commandlimits":      "A map of timeouts for commands. A value of 30 means the command can't be used more than once every 30 seconds.",
		"commanddisabled":    "A list of disabled commands.",
		"commandperduration": "Maximum number of commands that can be run within `commandmaxduration` seconds. Default: 3",
		"commandmaxduration": "Default: 20. This means that by default, at most 3 commands can be run every 20 seconds.",
		"disabled":           "A list of disabled modules.",
		"channels":           "A mapping of what channels a given module can operate on. If no mapping is given, a module operates on all channels. If \"!\" is included as a channel, it switches from a whitelist to a blacklist, enabling you to exclude certain channels instead of allow certain channels.",
	},
	"spam": map[string]string{
		"imagepressure":      "Additional pressure generated by each image, link or attachment in a message. Defaults to (MaxPressure - BasePressure) / 6, instantly silencing anyone posting 6 or more links at once.",
		"repeatpressure":     "Additional pressure generated by a message that is identical to the previous message sent (ignores case). Defaults to BasePressure, effectively doubling the pressure penalty for repeated messages.",
		"pingpressure":       "Additional pressure generated by each unique ping in a message. Defaults to (MaxPressure - BasePressure) / 20, instantly silencing anyone pinging 20 or more people at once.",
		"lengthpressure":     "Additional pressure generated by each individual character in the message. Discord allows messages up to 2000 characters in length. Defaults to (MaxPressure - BasePressure) / 8000, silencing anyone posting 3 huge messages at the same time.",
		"linepressure":       "Additional pressure generated by each newline in the message. Defaults to (MaxPressure - BasePressure) / 70, silencing anyone posting more than 70 newlines in a single message",
		"basepressure":       "The base pressure generated by sending a message, regardless of length or content. Defaults to 10",
		"maxpressure":        "The maximum pressure allowed. If a user's pressure exceeds this amount, they will be silenced. Defaults to 60, which is intended to ban after a maximum of 6 short messages sent in rapid succession.",
		"maxchannelpressure": "Per-channel pressure override. If a channel's pressure is specified in this map, it will override the global maxpressure setting.",
		"pressuredecay":      "The number of seconds it takes for a user to lose Spam.BasePressure from their pressure amount. Defaults to 2.5, so after sending 3 messages, it will take 7.5 seconds for their pressure to return to 0.",
		"maxremovelookback":  "Number of seconds back the bot should delete messages of a silenced user on the channel they spammed on. If set to 0, the bot will only delete the message that caused the user to be silenced. If less than 0, the bot won't delete any messages.",
		"ignorerole":         "If set, the bot will exclude anyone with this role from spam detection. Use with caution.",
		"raidtime":           "In order to trigger a raid alarm, at least `spam.raidsize` people must join the chat within this many seconds of each other.",
		"raidsize":           "Specifies how many people must have joined the server within the `spam.raidtime` period to qualify as a raid.",
		"autosilence":        "Gets the current autosilence state. Use the `!autosilence` command to set this.",
		"lockdownduration":   "Determines how long the server's verification mode will temporarily be increased to tableflip levels after a raid is detected. If set to 0, disables lockdown entirely.",
	},
	"bucket": map[string]string{
		"maxitems":       "Determines the maximum number of items that can be carried in the bucket. If set to 0, the bucket is disabled.",
		"maxitemlength":  "Determines the maximum length of a string that can be added to the bucket.",
		"maxfighthp":     "Maximum HP of the randomly generated enemy for the `!fight` command.",
		"maxfightdamage": "Maximum amount of damage a randomly generated weapon can deal for the `!fight` command.",
		"items":          "List of items in the bucket.",
	},
	"markov": map[string]string{
		"maxpmlines":     "This is the maximum number of lines a response can be before its automatically sent as a PM to avoid cluttering the chat. Default: 5",
		"maxlines":       "Maximum number of lines the `!episodequote` command can be given.",
		"defaultlines":   "Number of lines for the markov chain to spawn when not given a line count.",
		"usemembernames": "Use member names instead of random pony names.",
	},
	"users": map[string]string{
		"timezonelocation": "Sets the timezone location of the server itself. When no user timezone is available, the bot will use this.",
		"welcomechannel":   "If set to a channel ID, the bot will treat this channel as a \"quarantine zone\" for silenced members. If autosilence is enabled, new users will be sent to this channel.",
		"welcomemessage":   "If autosilence is enabled, this message will be sent to a new user upon joining.",
		"silencemessage":   "This message will be sent to users that have been silenced by the `!silence` command.",
		"roles":            "A list of all user-assignable roles. Manage it via !addrole and !removerole",
		"notifychannel":    "If set to a channel ID other than zero, sends a message to that channel whenever a new user joins the server.",
		"trackuserleft":    "If true, tracks users that leave the server if notifychannel is set.",
	},
	"filter": map[string]string{
		"filters":   "A collection of word lists for each filter. These are combined into a single regex of the form `(word1|word2|etc...)`, depending on the filter template.",
		"channels":  "A collection of channel exclusions for each filter.",
		"responses": "The response message sent by each filter when triggered.",
		"templates": "The template used to construct the regex. `%%` is replaced with `(word1|word2|etc...)` using the filter's word list. Example: `\\[\\]\\(\\/r?%%[-) \"]` is transformed into `\\[\\]\\(\\/r?(word1|word2)[-) \"]`",
	},
	"bored": map[string]string{
		"cooldown": "The bored cooldown timer, in seconds. This is the length of time a channel must be inactive before a bored message is posted.",
		"commands": "This determines what commands will be run when nothing has been said in a channel for a while. One command will be chosen from this list at random.\n\nExample: `!setconfig bored.commands !drop \"!pick bored\"`",
	},
	"information": map[string]string{
		"rules":             "Contains a list of numbered rules. The numbers do not need to be contiguous, and can be negative.",
		"hidenegativerules": "If true, `!rules -1` will display a rule at index -1, but `!rules` will not. This is useful for joke rules or additional rules that newcomers don't need to know about.",
	},
	"log": map[string]string{
		"channel":  "This is the channel where log output is sent.",
		"cooldown": "The cooldown time to display an error message, in seconds, intended to prevent the bot from spamming itself. Default: 4",
	},
	"witty": map[string]string{
		"responses": "Stores the replies used by the Witty module and must be configured using `!addwit` or `!removewit`",
		"cooldown":  "The cooldown time for the witty module. At least this many seconds must have passed before the bot will make another witty reply.",
	},
	"scheduler": map[string]string{
		"birthdayrole": " This is the role given to members on their birthday.",
	},
	"miscellaneous": map[string]string{
		"maxsearchresults": "Maximum number of search results that can be requested at once.",
	},
	"spoiler": map[string]string{
		"channels": "A list of channels that are exempt from the spoiler rules.",
	},
	"status": map[string]string{
		"cooldown": "Number of seconds the bot waits before changing its status to a string picked randomly from the `status` collection.",
		"lines":    "List of possible status messages that the bot can have.",
	},
	"quote": map[string]string{
		"quotes": "This is a map of quotes, which should be managed via `!addquote` and `!removequote`.",
	},
}

func getConfigHelp(module string, option string) (string, bool) {
	x, ok := ConfigHelp[strings.ToLower(module)]
	if !ok {
		return "", false
	}
	s, b := x[strings.ToLower(option)]
	return s, b
}

// ConfigVersion is the latest version of the config file
var ConfigVersion = 21

// DefaultConfig returns a default BotConfig struct. We can't define this as a variable because you can't initialize nested structs in a sane way in Go
func DefaultConfig() *BotConfig {
	config := &BotConfig{
		Version:     ConfigVersion,
		LastVersion: BotVersion.Integer(),
		SetupDone:   false,
	}
	config.Basic.IgnoreInvalidCommands = false
	config.Basic.Importable = false
	config.Basic.CommandPrefix = "!"
	config.Modules.CommandPerDuration = 3
	config.Modules.CommandMaxDuration = 15
	config.Spam.MaxPressure = 60
	config.Spam.BasePressure = 10
	config.Spam.ImagePressure = (config.Spam.MaxPressure - config.Spam.BasePressure) / 6
	config.Spam.PingPressure = (config.Spam.MaxPressure - config.Spam.BasePressure) / 20
	config.Spam.LengthPressure = (config.Spam.MaxPressure - config.Spam.BasePressure) / 8000
	config.Spam.RepeatPressure = config.Spam.BasePressure
	config.Spam.LinePressure = (config.Spam.MaxPressure - config.Spam.BasePressure) / 70
	config.Spam.PressureDecay = 2.5
	config.Spam.MaxRemoveLookback = 4
	config.Spam.RaidTime = 240
	config.Spam.RaidSize = 4
	config.Spam.AutoSilence = 1 // Default to raid mode
	config.Spam.LockdownDuration = 120
	config.Bucket.MaxItems = 10
	config.Bucket.MaxItemLength = 100
	config.Bucket.MaxFightHP = 300
	config.Bucket.MaxFightDamage = 60
	config.Markov.MaxPMlines = 5
	config.Markov.MaxLines = 30
	config.Markov.DefaultLines = 5
	config.Markov.UseMemberNames = true
	config.Bored.Cooldown = 500
	config.Bored.Commands = map[string]bool{"!quote": true, "!drop": true}
	config.Log.Cooldown = 4
	config.Witty.Cooldown = 180
	config.Miscellaneous.MaxSearchResults = 10
	config.Status.Cooldown = 3600

	return config
}

// FixRequest takes a request that is not fully qualified and attempts to find a fully qualified version
func FixRequest(arg string, t reflect.Value) (string, error) {
	args := strings.SplitN(strings.ToLower(arg), ".", 3)
	list := []string{}
	n := t.NumField()

	for i := 0; i < n; i++ {
		if strings.ToLower(t.Type().Field(i).Name) == args[0] {
			return arg, nil
		}
	}

	for i := 0; i < n; i++ {
		switch t.Field(i).Kind() {
		case reflect.Struct:
			f := t.Field(i)
			for j := 0; j < f.NumField(); j++ {
				if strings.ToLower(f.Type().Field(j).Name) == args[0] {
					list = append(list, t.Type().Field(i).Name)
				}
			}
		}
	}
	if len(list) < 1 {
		return arg, nil
	}
	if len(list) == 1 {
		return strings.ToLower(list[0]) + "." + arg, nil
	}
	for k := range list {
		list[k] += "." + args[0]
	}
	return "", errors.New("```\nCould be any of the following:\n" + strings.Join(list, "\n") + "```")
}

func setConfigValue(f reflect.Value, value string, info *GuildInfo) error {
	switch f.Interface().(type) {
	case string:
		f.SetString(value)
	case DiscordRole:
		g, _ := info.GetGuild()
		s, err := ParseRole(value, g)
		if err != nil {
			return err
		}
		f.SetString(s.String())
	case DiscordChannel:
		g, _ := info.GetGuild()
		s, err := ParseChannel(value, g)
		if err != nil {
			return err
		}
		f.SetString(s.String())
	case DiscordUser:
		s, err := ParseUser(value, info)
		if err != nil {
			return err
		}
		f.SetString(s.String())
	case ModuleID:
		value = strings.ToLower(value)
		for _, v := range info.Modules {
			if value == strings.ToLower(v.Name()) {
				f.SetString(value)
				return nil
			}
		}
		return fmt.Errorf("%s is not a module name!", value)
	case CommandID:
		value = strings.ToLower(value)
		if _, ok := info.commands[CommandID(value)]; !ok {
			return fmt.Errorf("%s is not a command name!", value)
		}
		f.SetString(value)
	case int, int8, int16, int32, int64:
		k, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(k)
	case uint, uint8, uint16, uint32, uint64:
		k, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		f.SetUint(k)
	case float32, float64:
		k, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		f.SetFloat(k)
	}
	return nil
}
func setConfigKeyValue(f reflect.Value, key string, value []string, info *GuildInfo) (string, bool) {
	if len(value) == 0 {
		return "No value parameter given", false
	}
	k := reflect.New(f.Type().Key()).Elem()
	if err := setConfigValue(k, key, info); err != nil {
		return "Key error: " + err.Error(), false
	}
	if f.IsNil() {
		f.Set(reflect.MakeMap(f.Type()))
	}
	if len(value[0]) == 0 {
		f.SetMapIndex(k, reflect.Value{})
		return "Deleted " + value[0], false
	}
	v := reflect.New(f.Type().Elem()).Elem()
	if err := setConfigValue(v, value[0], info); err != nil {
		return "Value error: " + err.Error(), false
	}

	f.SetMapIndex(k, v)
	return fmt.Sprintf("%v: %v", k.Interface(), v.Interface()), true
}

func setConfigList(f reflect.Value, values []string, info *GuildInfo) (string, bool) {
	switch f.Kind() {
	case reflect.Slice:
		f.Set(reflect.MakeSlice(f.Type(), 0, len(values)))
		if len(values[0]) > 0 {
			for _, value := range values {
				v := reflect.New(f.Type().Elem()).Elem()
				if err := setConfigValue(v, value, info); err != nil {
					return "Value error: " + err.Error(), false
				}
				f.Set(reflect.Append(f, v))
			}
		}
		return fmt.Sprint(f.Interface()), true
	case reflect.Map:
		if f.Type().Elem() != reflect.TypeOf(true) {
			return "Map sent into list function!", false
		}
		f.Set(reflect.MakeMap(f.Type()))
		stripped := []string{}
		if len(values[0]) > 0 {
			for _, value := range values {
				v := reflect.New(f.Type().Key()).Elem()
				if err := setConfigValue(v, value, info); err != nil {
					return "Value error: " + err.Error(), false
				}
				f.SetMapIndex(v, reflect.ValueOf(true))
				stripped = append(stripped, fmt.Sprint(v.Interface()))
			}
		}
		return "[" + strings.Join(stripped, ", ") + "]", true
	}
	return "Unknown list type!", false
}

func setConfigMapList(f reflect.Value, key string, values []string, info *GuildInfo) (string, bool) {
	if f.IsNil() {
		f.Set(reflect.MakeMap(f.Type()))
	}
	if len(key) == 0 {
		return "No key specified", false
	}
	k := reflect.New(f.Type().Key()).Elem()
	if err := setConfigValue(k, key, info); err != nil {
		return "Key error: " + err.Error(), false
	}
	if len(values) == 0 {
		return deleteFromMapReflect(f, k), false
	}

	v := reflect.New(f.Type().Elem()).Elem()
	s, ok := setConfigList(v, values, info)
	if !ok {
		return s, false
	}
	f.SetMapIndex(k, v)
	return fmt.Sprintf("%v: %s", k, s), true
}

// SetConfig sets the given config option with the given value along with any extra parameters
func (config *BotConfig) SetConfig(info *GuildInfo, name string, value string, extra ...string) (string, bool) {
	names := strings.SplitN(strings.ToLower(name), ".", 3)
	t := reflect.ValueOf(config).Elem()
	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Type().Field(i).Name) == names[0] {
			if len(names) < 2 {
				return "Can't set a configuration category! Use \"Category.Option\" to set a specific option.", false
			}
			switch t.Field(i).Kind() {
			case reflect.Struct:
				for j := 0; j < t.Field(i).NumField(); j++ {
					if strings.ToLower(t.Field(i).Type().Field(j).Name) == names[1] {
						f := t.Field(i).Field(j)
						switch f.Interface().(type) {
						case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, uint64, DiscordChannel, DiscordRole, DiscordUser:
							if err := setConfigValue(f, value, info); err != nil {
								return "Error: " + err.Error(), false
							}
						case map[DiscordChannel]bool, map[string]bool, map[DiscordRole]bool, map[CommandID]bool, map[ModuleID]bool:
							return setConfigList(f, append([]string{value}, extra...), info)
						case bool:
							switch strings.ToLower(value) {
							case "true":
								f.SetBool(true)
							case "false":
								f.SetBool(false)
							default:
								return name + " must be set to either 'true' or 'false'", false
							}
						case map[string]string, map[CommandID]int64, map[DiscordChannel]float32, map[int]string:
							return setConfigKeyValue(f, strings.ToLower(value), extra, info)
						case map[string]map[DiscordChannel]bool, map[CommandID]map[DiscordRole]bool, map[string]map[string]bool, map[DiscordUser][]string, map[CommandID]map[DiscordChannel]bool, map[ModuleID]map[DiscordChannel]bool:
							return setConfigMapList(f, strings.ToLower(value), extra, info)
						default:
							return "That config option has an unknown type!", false
						}
						return fmt.Sprint(f.Interface()), true
					}
				}
			default:
				return "Not a configuration category!", false
			}
		}
	}
	return "Could not find configuration parameter " + name + "!", false
}

func getConfigValue(f reflect.Value, state *discordgo.State, guild string) string {
	switch f.Interface().(type) {
	case DiscordRole:
		if r, err := state.Role(guild, f.String()); err == nil {
			return "@" + r.Name
		}
	case DiscordChannel:
		if ch, err := state.Channel(f.String()); err == nil {
			return "#" + ch.Name
		}
	case DiscordUser:
		if m, err := state.Member(guild, f.String()); err == nil {
			if len(m.Nick) > 0 {
				return m.Nick
			}
			return m.User.Username
		}
		//case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
	}
	return fmt.Sprint(f.Interface())
}
func getConfigList(f reflect.Value, state *discordgo.State, guild string) (s []string) {
	switch f.Kind() {
	case reflect.Slice:
		for i := 0; i < f.Len(); i++ {
			s = append(s, getConfigValue(f.Index(i), state, guild))
		}
	case reflect.Map:
		keys := f.MapKeys()
		if f.Type().Elem() == reflect.TypeOf(true) {
			for _, key := range keys {
				s = append(s, getConfigValue(key, state, guild))
			}
		} else {
			for _, key := range keys {
				s = append(s, "\""+getConfigValue(key, state, guild)+"\": "+getConfigValue(f.MapIndex(key), state, guild))
			}
		}
	}
	return
}

func getConfigMapList(f reflect.Value, state *discordgo.State, guild string) (s []string) {
	keys := f.MapKeys()
	for _, key := range keys {
		v := f.MapIndex(key)
		k := getConfigValue(key, state, guild)

		if v.Len() == 1 {
			s = append(s, fmt.Sprintf("\"%s\": %s", k, strings.Join(getConfigList(v, state, guild), ", ")))
		} else {
			s = append(s, fmt.Sprintf("\"%s\": [%v items]", k, v.Len()))
		}
	}
	return
}

func (config *BotConfig) GetConfig(f reflect.Value, state *discordgo.State, guild string) (s []string) {
	switch f.Interface().(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, uint64, DiscordChannel, DiscordRole, DiscordUser, ModuleID, CommandID, bool:
		s = append(s, getConfigValue(f, state, guild))
	case map[DiscordChannel]bool, map[string]bool, map[DiscordRole]bool, map[string]string, map[CommandID]int64, map[DiscordChannel]float32, map[int]string, map[CommandID]bool, map[ModuleID]bool:
		s = getConfigList(f, state, guild)
	case map[string]map[DiscordChannel]bool, map[CommandID]map[DiscordRole]bool, map[string]map[string]bool, map[DiscordUser][]string, map[CommandID]map[DiscordChannel]bool, map[ModuleID]map[DiscordChannel]bool:
		s = getConfigMapList(f, state, guild)
	default:
		data, err := json.Marshal(f.Interface())
		if err != nil {
			s = append(s, "[JSON Error]")
		} else {
			s = append(s, string(data))
		}
	}
	return
}

// IsModuleDisabled returns a string if a module is disabled
func (config *BotConfig) IsModuleDisabled(module Module) string {
	_, ok := config.Modules.Disabled[ModuleID(strings.ToLower(module.Name()))]
	if ok {
		return " [disabled]"
	}
	return ""
}

// IsCommandDisabled returns a string if a command is disabled
func (config *BotConfig) IsCommandDisabled(command Command) (str string) {
	_, disabled := config.Modules.CommandDisabled[CommandID(strings.ToLower(command.Info().Name))]
	if disabled {
		str = " [disabled]"
	}
	return
}

// FillConfig ensures root maps are not nil
func (config *BotConfig) FillConfig() {
	if len(config.Basic.FreeChannels) == 0 {
		config.Basic.FreeChannels = make(map[DiscordChannel]bool)
	}
	if len(config.Basic.Aliases) == 0 {
		config.Basic.Aliases = make(map[string]string)
	}
	if len(config.Modules.Channels) == 0 {
		config.Modules.Channels = make(map[ModuleID]map[DiscordChannel]bool)
	}
	if len(config.Modules.Disabled) == 0 {
		config.Modules.Disabled = make(map[ModuleID]bool)
	}
	if len(config.Modules.CommandRoles) == 0 {
		config.Modules.CommandRoles = make(map[CommandID]map[DiscordRole]bool)
	}
	if len(config.Modules.CommandChannels) == 0 {
		config.Modules.CommandChannels = make(map[CommandID]map[DiscordChannel]bool)
	}
	if len(config.Modules.CommandLimits) == 0 {
		config.Modules.CommandLimits = make(map[CommandID]int64)
	}
	if len(config.Modules.CommandDisabled) == 0 {
		config.Modules.CommandDisabled = make(map[CommandID]bool)
	}
	if len(config.Spam.MaxChannelPressure) == 0 {
		config.Spam.MaxChannelPressure = make(map[DiscordChannel]float32)
	}
	if len(config.Users.Roles) == 0 {
		config.Users.Roles = make(map[DiscordRole]bool)
	}
	if len(config.Bucket.Items) == 0 {
		config.Bucket.Items = make(map[string]bool)
	}
	if len(config.Filter.Filters) == 0 {
		config.Filter.Filters = make(map[string]map[string]bool)
	}
	if len(config.Filter.Channels) == 0 {
		config.Filter.Channels = make(map[string]map[DiscordChannel]bool)
	}
	if len(config.Filter.Responses) == 0 {
		config.Filter.Responses = make(map[string]string)
	}
	if len(config.Filter.Templates) == 0 {
		config.Filter.Templates = make(map[string]string)
	}
	if len(config.Bored.Commands) == 0 {
		config.Bored.Commands = make(map[string]bool)
	}
	if len(config.Information.Rules) == 0 {
		config.Information.Rules = make(map[int]string)
	}
	if len(config.Witty.Responses) == 0 {
		config.Witty.Responses = make(map[string]string)
	}
	if len(config.Status.Lines) == 0 {
		config.Status.Lines = make(map[string]bool)
	}
	if len(config.Quote.Quotes) == 0 {
		config.Quote.Quotes = make(map[DiscordUser][]string)
	}
}

type legacyBotConfig struct {
	Version               int                        `json:"version"`
	LastVersion           int                        `json:"lastversion"`
	Maxerror              int64                      `json:"maxerror"`
	Maxwit                int64                      `json:"maxwit"`
	Maxbored              int64                      `json:"maxbored"`
	BoredCommands         map[string]bool            `json:"boredcommands"`
	MaxPMlines            int                        `json:"maxpmlines"`
	Maxquotelines         int                        `json:"maxquotelines"`
	Maxsearchresults      int                        `json:"maxsearchresults"`
	Defaultmarkovlines    int                        `json:"defaultmarkovlines"`
	Commandperduration    int                        `json:"commandperduration"`
	Commandmaxduration    int64                      `json:"commandmaxduration"`
	StatusDelayTime       int                        `json:"statusdelaytime"`
	MaxRaidTime           int64                      `json:"maxraidtime"`
	RaidSize              int                        `json:"raidsize"`
	Witty                 map[string]string          `json:"witty"`
	Aliases               map[string]string          `json:"aliases"`
	MaxBucket             int                        `json:"maxbucket"`
	MaxBucketLength       int                        `json:"maxbucketlength"`
	MaxFightHP            int                        `json:"maxfighthp"`
	MaxFightDamage        int                        `json:"maxfightdamage"`
	MaxImageSpam          int                        `json:"maximagespam"`
	MaxAttachSpam         int                        `json:"maxattachspam"`
	MaxPingSpam           int                        `json:"maxpingspam"`
	MaxMessageSpam        map[int64]int              `json:"maxmessagespam"`
	MaxSpamRemoveLookback int                        `json:maxspamremovelookback`
	IgnoreInvalidCommands bool                       `json:"ignoreinvalidcommands"`
	UseMemberNames        bool                       `json:"usemembernames"`
	Importable            bool                       `json:"importable"`
	HideNegativeRules     bool                       `json:"hidenegativerules"`
	Timezone              int                        `json:"timezone"`
	TimezoneLocation      string                     `json:"timezonelocation"`
	AutoSilence           int                        `json:"autosilence"`
	AlertRole             uint64                     `json:"alertrole"`
	SilentRole            uint64                     `json:"silentrole"`
	LogChannel            uint64                     `json:"logchannel"`
	ModChannel            uint64                     `json:"modchannel"`
	WelcomeChannel        uint64                     `json:"welcomechannel"`
	WelcomeMessage        string                     `json:"welcomemessage"`
	SilenceMessage        string                     `json:"silencemessage"`
	BirthdayRole          uint64                     `json:"birthdayrole"`
	SpoilChannels         []uint64                   `json:"spoilchannels"`
	FreeChannels          map[string]bool            `json:"freechannels"`
	Command_roles         map[string]map[string]bool `json:"command_roles"`
	Command_channels      map[string]map[string]bool `json:"command_channels"`
	Command_limits        map[string]int64           `json:command_limits`
	Command_disabled      map[string]bool            `json:command_disabled`
	Module_disabled       map[string]bool            `json:module_disabled`
	Module_channels       map[string]map[string]bool `json:module_channels`
	Collections           map[string]map[string]bool `json:"collections"`
	Groups                map[string]map[string]bool `json:"groups"`
	Quotes                map[uint64][]string        `json:"quotes"`
	Rules                 map[int]string             `json:"rules"`
}

type legacyBotConfigV10 struct {
	Basic struct {
		Commandperduration int   `json:"commandperduration"`
		Commandmaxduration int64 `json:"commandmaxduration"`
	} `json:"basic"`
}

type legacyBotConfigV12 struct {
	Spam struct {
		MaxImages int `json:"maximagespam"`
		MaxPings  int `json:"maxpingspam"`
	} `json:"spam"`
}

type legacyBotConfigV13 struct {
	Basic struct {
		Groups map[string]map[string]bool `json:"groups"`
	} `json:"basic"`
}

type legacyBotConfigV19 struct {
	Basic struct {
		Collections map[string]map[string]bool `json:"collections"`
	} `json:"basic"`
}

type legacyBotConfigV20 struct {
	Collections map[string]map[string]bool `json:"collections"`
	Spam        struct {
		SilentRole     DiscordRole `json:"silentrole"`
		SilenceMessage string      `json:"silencemessage"`
	} `json:"spam"`
	Basic struct {
		AlertRole     DiscordRole `json:"alertrole"`
		TrackUserLeft bool        `json:"trackuserleft"`
	} `json:"basic"`
	Search struct {
		MaxResults int `json:"maxsearchresults"`
	} `json:"search"`
	Spoiler struct {
		Channels []DiscordChannel `json:"spoilchannels"`
	} `json:"spoiler"`
	Schedule struct {
		BirthdayRole DiscordRole `json:"birthdayrole"`
	} `json:"schedule"`
}

func restrictCommand(v string, roles map[CommandID]map[DiscordRole]bool, modrole DiscordRole) {
	id := CommandID(v)
	_, ok := roles[id]
	if !ok && modrole != "" {
		roles[id] = make(map[DiscordRole]bool)
		roles[id][modrole] = true
	}
}

// MigrateSettings from earlier config version
func (guild *GuildInfo) MigrateSettings(config []byte) error {
	err := json.Unmarshal(config, &guild.Config)
	if err != nil {
		return err
	}

	if guild.Config.Version < 10 {
		legacy := legacyBotConfig{}
		err := json.Unmarshal(config, &legacy)
		if err != nil {
			return err
		}

		if legacy.Version == 0 {
			if len(legacy.Command_roles) == 0 {
				legacy.Command_roles = make(map[string]map[string]bool)
			}
			legacy.MaxImageSpam = 3
			legacy.MaxAttachSpam = 1
			legacy.MaxPingSpam = 24
			legacy.MaxMessageSpam = make(map[int64]int)
			legacy.MaxMessageSpam[1] = 4
			legacy.MaxMessageSpam[9] = 10
			legacy.MaxMessageSpam[12] = 15
		}

		if legacy.Version <= 1 {
			if len(legacy.Aliases) == 0 {
				legacy.Aliases = make(map[string]string)
			}
			legacy.Aliases["cute"] = "pick cute"
		}

		if legacy.Version <= 3 {
			legacy.BoredCommands = make(map[string]bool)
		}

		if legacy.Version <= 5 {
			legacy.TimezoneLocation = "Etc/GMT"
			if legacy.Timezone < 0 {
				legacy.TimezoneLocation += "+"
			}
			legacy.TimezoneLocation += strconv.Itoa(-legacy.Timezone) // Etc has the sign reversed
		}

		guild.Config.Basic.ModRole = NewDiscordRole(legacy.AlertRole)
		guild.Config.Basic.Aliases = legacy.Aliases
		guild.Config.Filter.Filters = legacy.Collections
		guild.Config.Basic.FreeChannels = make(map[DiscordChannel]bool)
		for k, v := range legacy.FreeChannels {
			if ch, err := ParseChannel(k, nil); err == nil {
				guild.Config.Basic.FreeChannels[ch] = v
			}
		}
		guild.Config.Basic.IgnoreInvalidCommands = legacy.IgnoreInvalidCommands
		guild.Config.Basic.Importable = legacy.Importable
		guild.Config.Basic.ModChannel = NewDiscordChannel(legacy.ModChannel)
		guild.Config.Basic.SilenceRole = NewDiscordRole(legacy.SilentRole)
		guild.Config.Modules.CommandChannels = make(map[CommandID]map[DiscordChannel]bool)
		for key, _ := range legacy.Command_channels {
			guild.Config.Modules.CommandChannels[CommandID(key)] = make(map[DiscordChannel]bool)
			for k, v := range legacy.Command_channels[key] {
				if ch, err := ParseChannel(k, nil); err == nil {
					guild.Config.Modules.CommandChannels[CommandID(key)][ch] = v
				}
			}
		}
		guild.Config.Modules.CommandDisabled = make(map[CommandID]bool)
		for key, _ := range legacy.Command_disabled {
			guild.Config.Modules.CommandDisabled[CommandID(key)] = true
		}
		guild.Config.Modules.CommandLimits = make(map[CommandID]int64)
		for key, v := range legacy.Command_limits {
			guild.Config.Modules.CommandLimits[CommandID(key)] = v
		}
		guild.Config.Modules.CommandRoles = make(map[CommandID]map[DiscordRole]bool)
		for key, _ := range legacy.Command_roles {
			guild.Config.Modules.CommandRoles[CommandID(key)] = make(map[DiscordRole]bool)
			for k, v := range legacy.Command_roles[key] {
				if r, err := ParseRole(k, nil); err == nil {
					guild.Config.Modules.CommandRoles[CommandID(key)][r] = v
				}
			}
		}

		guild.Config.Modules.CommandMaxDuration = legacy.Commandmaxduration
		guild.Config.Modules.CommandPerDuration = legacy.Commandperduration
		guild.Config.Modules.Channels = make(map[ModuleID]map[DiscordChannel]bool)
		for key, _ := range legacy.Module_channels {
			guild.Config.Modules.Channels[ModuleID(key)] = make(map[DiscordChannel]bool)
			for k, v := range legacy.Module_channels[key] {
				if ch, err := ParseChannel(k, nil); err == nil {
					guild.Config.Modules.Channels[ModuleID(key)][ch] = v
				}
			}
		}
		guild.Config.Modules.Disabled = make(map[ModuleID]bool)
		for key, _ := range legacy.Module_disabled {
			guild.Config.Modules.Disabled[ModuleID(key)] = true
		}
		guild.Config.Spam.AutoSilence = legacy.AutoSilence
		//guild.Config.Spam.MaxAttach = legacy.MaxAttachSpam
		//guild.Config.Spam.MaxImages = legacy.MaxImageSpam
		//guild.Config.Spam.MaxMessages = legacy.MaxMessageSpam
		//guild.Config.Spam.MaxPings = legacy.MaxPingSpam
		guild.Config.Spam.RaidTime = legacy.MaxRaidTime
		guild.Config.Spam.MaxRemoveLookback = legacy.MaxSpamRemoveLookback
		guild.Config.Spam.RaidSize = legacy.RaidSize
		guild.Config.Bucket.MaxItems = legacy.MaxBucket
		guild.Config.Bucket.MaxItemLength = legacy.MaxBucketLength
		guild.Config.Bucket.MaxFightDamage = legacy.MaxFightDamage
		guild.Config.Bucket.MaxFightHP = legacy.MaxFightHP
		guild.Config.Markov.DefaultLines = legacy.Defaultmarkovlines
		guild.Config.Markov.MaxPMlines = legacy.MaxPMlines
		guild.Config.Markov.MaxLines = legacy.Maxquotelines
		guild.Config.Markov.UseMemberNames = legacy.UseMemberNames
		guild.Config.Users.TimezoneLocation = legacy.TimezoneLocation
		guild.Config.Users.WelcomeChannel = NewDiscordChannel(legacy.WelcomeChannel)
		guild.Config.Users.WelcomeMessage = legacy.WelcomeMessage
		guild.Config.Users.SilenceMessage = legacy.SilenceMessage
		guild.Config.Bored.Commands = legacy.BoredCommands
		guild.Config.Bored.Cooldown = legacy.Maxbored
		guild.Config.Information.HideNegativeRules = legacy.HideNegativeRules
		guild.Config.Information.Rules = legacy.Rules
		guild.Config.Log.Channel = NewDiscordChannel(legacy.LogChannel)
		guild.Config.Log.Cooldown = legacy.Maxerror
		guild.Config.Witty.Cooldown = legacy.Maxwit
		guild.Config.Witty.Responses = legacy.Witty
		guild.Config.Scheduler.BirthdayRole = NewDiscordRole(legacy.BirthdayRole)
		guild.Config.Miscellaneous.MaxSearchResults = legacy.Maxsearchresults
		guild.Config.Filter.Channels = make(map[string]map[DiscordChannel]bool)
		guild.Config.Filter.Channels["spoiler"] = make(map[DiscordChannel]bool)
		for _, v := range legacy.SpoilChannels {
			guild.Config.Filter.Channels["spoiler"][NewDiscordChannel(v)] = true
		}
		guild.Config.Status.Cooldown = legacy.StatusDelayTime
		guild.Config.Quote.Quotes = make(map[DiscordUser][]string)
		for k, v := range legacy.Quotes {
			guild.Config.Quote.Quotes[NewDiscordUser(k)] = v
		}

		newcommands := []string{"addevent", "addbirthday", "autosilence", "silence", "unsilence", "wipewelcome", "new", "addquote", "removequote", "removealias", "delete", "createpoll", "deletepoll", "addoption"}
		for _, v := range newcommands {
			restrictCommand(v, guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		}
	}

	if guild.Config.Version == 10 {
		legacy := legacyBotConfigV10{}
		err := json.Unmarshal(config, &legacy)
		if err == nil {
			guild.Config.Modules.CommandMaxDuration = legacy.Basic.Commandmaxduration
			guild.Config.Modules.CommandPerDuration = legacy.Basic.Commandperduration
		} else {
			fmt.Println(err.Error())
		}
	}

	if guild.Config.Version <= 11 {
		restrictCommand("getaudit", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
	}

	if guild.Config.Version <= 12 {
		guild.Config.Spam.BasePressure = 10.0
		guild.Config.Spam.MaxPressure = 60.0
		guild.Config.Spam.ImagePressure = ((guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / 6.0)
		guild.Config.Spam.PingPressure = ((guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / 24.0)
		guild.Config.Spam.LengthPressure = ((guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / (2000.0 * 4))
		guild.Config.Spam.RepeatPressure = guild.Config.Spam.BasePressure
		guild.Config.Spam.PressureDecay = 2.5

		legacy := legacyBotConfigV12{}
		err := json.Unmarshal(config, &legacy)
		if err == nil {
			if legacy.Spam.MaxImages > 0 {
				guild.Config.Spam.ImagePressure = ((guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / float32(legacy.Spam.MaxImages+1))
			} else {
				guild.Config.Spam.ImagePressure = 0
			}
			if legacy.Spam.MaxPings > 0 {
				guild.Config.Spam.PingPressure = ((guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / float32(legacy.Spam.MaxPings+1))
			} else {
				guild.Config.Spam.PingPressure = 0
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	if guild.Config.Version <= 13 {
		legacy := legacyBotConfigV13{}
		err := json.Unmarshal(config, &legacy)
		if err == nil {
			guild.Config.Users.Roles = make(map[DiscordRole]bool, len(legacy.Basic.Groups))
			idmap := make(map[string]string, len(legacy.Basic.Groups)) // Map initial group name to new role ID

			for k, v := range legacy.Basic.Groups {
				role := k
				check, err := GetRoleByName(role, guild)
				if check != nil {
					role = "sb-" + role
				}
				r, err := guild.Bot.DG.GuildRoleCreate(guild.ID)
				if err == nil {
					r, err = guild.Bot.DG.GuildRoleEdit(guild.ID, r.ID, role, 0, false, 0, true)
				}
				if err == nil {
					idmap[strings.ToLower(k)] = r.ID
					if id, err := ParseRole(r.ID, nil); err == nil {
						guild.Config.Users.Roles[id] = true
					}

					for u := range v {
						err = guild.Bot.DG.GuildMemberRoleAdd(guild.ID, u, r.ID)
						if err != nil {
							fmt.Println(err)
						}
					}
				} else {
					fmt.Println(err)
				}
			}

			stmt, err := guild.Bot.DB.Prepare("SELECT ID, Data FROM schedule WHERE Guild = ? AND Type = 7")
			stmt2, err := guild.Bot.DB.Prepare("UPDATE schedule SET Data = ? WHERE ID = ?")
			if err != nil {
				fmt.Println(err)
			} else {
				q, err := stmt.Query(SBatoi(guild.ID))
				if err != nil {
					fmt.Println(err)
				} else {
					defer q.Close()
					for q.Next() {
						var id uint64
						var dat string
						if err := q.Scan(&id, &dat); err == nil {
							datas := strings.SplitN(dat, "|", 2)
							groups := strings.Split(datas[0], "+")
							for i := range groups {
								rid, ok := idmap[strings.ToLower(groups[i])]
								if ok {
									groups[i] = "<@&" + rid + ">"
								}
							}
							_, err = stmt2.Exec(strings.Join(groups, " ")+"|"+datas[1], id)
							if err != nil {
								fmt.Println(err)
							}
						}
					}
				}
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	if guild.Config.Version <= 14 {
		restrictCommand("addrole", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("removerole", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("deleterole", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
	}

	if guild.Config.Version <= 15 {
		restrictCommand("bannewcomers", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		guild.Config.Spam.LockdownDuration = 120
	}

	if guild.Config.Version <= 16 {
		guild.Config.Basic.CommandPrefix = "!"
	}

	if guild.Config.Version <= 17 {
		guild.Config.SetupDone = true
	}

	if guild.Config.Version <= 18 {
		restrictCommand("banraid", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("getraid", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("wipe", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("bannewcomers", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("getpressure", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		guild.Config.Spam.LinePressure = (guild.Config.Spam.MaxPressure - guild.Config.Spam.BasePressure) / 70.0
	}

	if guild.Config.Version <= 19 {
		guild.Bot.GuildsLock.Lock()
		if len(guild.Config.Filter.Filters) == 0 {
			guild.Config.Filter.Filters = make(map[string]map[string]bool)
		}
		legacy := legacyBotConfigV19{}
		err := json.Unmarshal(config, &legacy)
		if err == nil {
			guild.Config.Bucket.Items = legacy.Basic.Collections["bucket"]
			guild.Config.Filter.Filters["emote"] = legacy.Basic.Collections["emote"]
			guild.Config.Status.Lines = legacy.Basic.Collections["status"]
			guild.Config.Filter.Filters["spoiler"] = legacy.Basic.Collections["spoiler"]
			delete(legacy.Basic.Collections, "bucket")
			delete(legacy.Basic.Collections, "emote")
			delete(legacy.Basic.Collections, "status")
			delete(legacy.Basic.Collections, "spoiler")

			gID := SBatoi(guild.ID)
			for k, v := range legacy.Basic.Collections {
				if len(v) > 0 {
					fmt.Println("Importing:", k)
					guild.Bot.DB.CreateTag(k, gID)
					tag, err := guild.Bot.DB.GetTag(k, gID)
					if err == nil {
						for item := range v {
							id, err := guild.Bot.DB.AddItem(item)
							if err == nil || err != ErrDuplicateEntry {
								guild.Bot.DB.AddTag(id, tag)
							}
						}
					}
				} else {
					fmt.Println("Skipping empty collection:", k)
				}
			}
		} else {
			fmt.Println(err.Error())
		}
		guild.Bot.GuildsLock.Unlock()
		restrictCommand("addset", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("removeset", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
		restrictCommand("searchset", guild.Config.Modules.CommandRoles, guild.Config.Basic.ModRole)
	}

	if guild.Config.Version <= 20 {
		legacy := legacyBotConfigV20{}
		err := json.Unmarshal(config, &legacy)
		if err == nil {
			guild.Config.Basic.ModRole = legacy.Basic.AlertRole
			guild.Config.Miscellaneous.MaxSearchResults = legacy.Search.MaxResults
			guild.Config.Scheduler.BirthdayRole = legacy.Schedule.BirthdayRole
			guild.Config.Filter.Filters = make(map[string]map[string]bool)
			guild.Config.Filter.Channels = make(map[string]map[DiscordChannel]bool)
			guild.Config.Filter.Responses = make(map[string]string)
			guild.Config.Filter.Templates = make(map[string]string)
			guild.Config.Bucket.Items = make(map[string]bool)
			guild.Config.Status.Lines = make(map[string]bool)
			guild.Config.Users.TrackUserLeft = legacy.Basic.TrackUserLeft
			guild.Config.Users.SilenceMessage = legacy.Spam.SilenceMessage
			guild.Config.Basic.SilenceRole = legacy.Spam.SilentRole

			if bucket, ok := legacy.Collections["bucket"]; ok {
				for k, v := range bucket {
					guild.Config.Bucket.Items[k] = v
				}
			}

			if status, ok := legacy.Collections["status"]; ok {
				for k, v := range status {
					guild.Config.Status.Lines[k] = v
				}
			}

			if guild.Config.Spam.AutoSilence == -2 {
				guild.Config.Users.NotifyChannel = guild.Config.Log.Channel
			} else if guild.Config.Spam.AutoSilence != 0 {
				guild.Config.Users.NotifyChannel = guild.Config.Basic.ModChannel
			}
			if guild.Config.Spam.AutoSilence < 0 {
				guild.Config.Spam.AutoSilence = 0
			}

			if spoilers, ok := legacy.Collections["spoiler"]; (ok && len(spoilers) > 0) || len(legacy.Spoiler.Channels) > 0 {
				guild.Config.Filter.Filters["spoiler"] = make(map[string]bool)
				if ok {
					for k, v := range spoilers {
						guild.Config.Filter.Filters["spoiler"][k] = v
					}
				}
				guild.Config.Filter.Channels["spoiler"] = make(map[DiscordChannel]bool)
				for _, v := range legacy.Spoiler.Channels {
					guild.Config.Filter.Channels["spoiler"][v] = true
				}
				guild.Config.Filter.Responses["spoiler"] = "[](/nospoilers) ```\nNO SPOILERS! Posting spoilers is a bannable offense. All discussion about new and future content MUST be in #mylittlespoilers.```"
			}

			if emotes, ok := legacy.Collections["emote"]; ok && len(emotes) > 0 {
				guild.Config.Filter.Filters["emote"] = make(map[string]bool)
				for k, v := range emotes {
					guild.Config.Filter.Filters["emote"][k] = v
				}
				guild.Config.Filter.Channels["emote"] = make(map[DiscordChannel]bool)
				guild.Config.Filter.Responses["emote"] = "```\nThat emote isn't allowed here! Try to avoid using large or disturbing emotes, as they can be problematic.```"
				guild.Config.Filter.Templates["emote"] = "\\[\\]\\(\\/r?%%[-) \"]"
			}
		}

		if guild.Config.Basic.ModRole == "0" {
			guild.Config.Basic.ModRole = ""
		}
		if guild.Config.Basic.ModChannel == "0" {
			guild.Config.Basic.ModChannel = ""
		}
		if guild.Config.Basic.SilenceRole == "0" {
			guild.Config.Basic.SilenceRole = ""
		}
		if guild.Config.Spam.IgnoreRole == "0" {
			guild.Config.Spam.IgnoreRole = ""
		}
		if guild.Config.Users.WelcomeChannel == "0" {
			guild.Config.Users.WelcomeChannel = ""
		}
		if guild.Config.Users.NotifyChannel == "0" {
			guild.Config.Users.NotifyChannel = ""
		}
		if guild.Config.Log.Channel == "0" {
			guild.Config.Log.Channel = ""
		}
		if guild.Config.Scheduler.BirthdayRole == "0" {
			guild.Config.Scheduler.BirthdayRole = ""
		}

		for k := range guild.Config.Modules.Channels {
			switch k {
			case "schedule":
				guild.Config.Modules.Channels["scheduler"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			case "anti-spam":
				guild.Config.Modules.Channels["spam"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			case "help/about":
				guild.Config.Modules.Channels["information"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			}
		}

		for k := range guild.Config.Modules.Disabled {
			switch k {
			case "schedule":
				guild.Config.Modules.Channels["scheduler"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			case "anti-spam":
				guild.Config.Modules.Channels["spam"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			case "help/about":
				guild.Config.Modules.Channels["information"] = guild.Config.Modules.Channels[k]
				delete(guild.Config.Modules.Channels, k)
			}
		}
	}

	if guild.Config.Version != ConfigVersion {
		guild.Config.Version = ConfigVersion // set version to most recent config version
		guild.SaveConfig()
	}
	return nil
}