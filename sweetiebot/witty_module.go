package sweetiebot

import (
  "github.com/bwmarrin/discordgo"
  "strings"
)

// This module is intended for any witty comments sweetie bot makes in response to what users say or do.
type WittyModule struct {
  maxwit int64
  lastdelete int64
  lastcomment int64
}

func (w *WittyModule) Name() string {
  return "Witty Module"
}

func (w *WittyModule) Register(hooks *ModuleHooks) {
  w.maxwit = 300
  w.lastdelete = 0
  w.lastcomment = 0
  hooks.OnMessageDelete = append(hooks.OnMessageDelete, w)
  hooks.OnMessageCreate = append(hooks.OnMessageCreate, w)
}
func (w *WittyModule) Channels() []string {
  return []string{}
}
  
func (w *WittyModule) SendWittyComment(channel string, comment string) {
  if RateLimit(&w.lastcomment, w.maxwit) {
    sb.dg.ChannelMessageSend(channel, comment)
  }
}
func (w *WittyModule)  OnMessageCreate(s *discordgo.Session, m *discordgo.Message) {
  if CheckRateLimit(&w.lastcomment, w.maxwit) {
    str := strings.ToLower(m.Content)
    if strings.Contains(str, "skynet") {
      w.SendWittyComment(m.ChannelID, "[](/dumbfabric) `SKYNET IS ALREADY HERE.`")
    } else if strings.Contains(str, "lewd") {
      w.SendWittyComment(m.ChannelID, "[](/ohcomeon) `This channel is SFW, remember?`")
    } else if strings.Contains(str, "memes") {
      w.SendWittyComment(m.ChannelID, "http://i1.kym-cdn.com/entries/icons/original/000/015/266/Z7HeRxU.png")
    } else if strings.Contains(str, "intensifies") {
      w.SendWittyComment(m.ChannelID, "[](/spikewoah)")
    } else if strings.Contains(str, "is best pony") {
      w.SendWittyComment(m.ChannelID, "[](/flutterjerk) `Your FACE is best pony.`")
    }
  }
}

func (w *WittyModule) OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
  //if RateLimit(&w.lastdelete, 60) { // It turns out this triggers when the bot itself deletes things, which looks awkward
  //  sb.dg.ChannelMessageSend(m.ChannelID, "[](/sbstare) `I SAW THAT`")
  //} 
}