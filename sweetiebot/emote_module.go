package sweetiebot

import (
  "github.com/bwmarrin/discordgo"
  "regexp"
)

// The emote module detects banned emotes and deletes them
type EmoteModule struct {
  ModuleEnabled
  emoteban *regexp.Regexp
  lastmsg int64
}

func (w *EmoteModule) Name() string {
  return "Emote"
}

func (w *EmoteModule) Register(hooks *ModuleHooks) {
  w.lastmsg = 0
  w.emoteban = regexp.MustCompile("\\[\\]\\(\\/r?(canada|BlockJuice|octybelleintensifies|angstybloom|alltheclops|bob|darklelicious|flutterbutts|juice|doitfor24|allthetables|ave|sbrapestare|gak|beforetacoswerecool|bigenough)[-) \"]")
  hooks.OnMessageCreate = append(hooks.OnMessageCreate, w)
  hooks.OnMessageUpdate = append(hooks.OnMessageUpdate, w)
  hooks.OnCommand = append(hooks.OnCommand, w)
}
func (w *EmoteModule) Channels() []string {
  return []string{}
}

func (w *EmoteModule) HasBigEmote(s *discordgo.Session, m *discordgo.Message) bool {
  if w.emoteban.Match([]byte(m.Content)) {
    s.ChannelMessageDelete(m.ChannelID, m.ID)
    if RateLimit(&w.lastmsg, 5) {
      s.ChannelMessageSend(m.ChannelID, "`That emote was way too big! Try to avoid using large emotes, as they can clutter up the chatroom.`")
    }
    return true
  }
  return false
}

func (w *EmoteModule) OnMessageCreate(s *discordgo.Session, m *discordgo.Message) {
  w.HasBigEmote(s, m)
}
  
func (w *EmoteModule) OnMessageUpdate(s *discordgo.Session, m *discordgo.Message) {
  w.HasBigEmote(s, m)
}

func (w *EmoteModule) OnCommand(s *discordgo.Session, m *discordgo.Message) bool {
  return w.HasBigEmote(s, m)
}