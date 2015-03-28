package main

import (
	"github.com/jebjerg/go-bot/bot"
	cfg "github.com/jebjerg/go-bot/bot/config"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	"math/rand"
	"net"
	"strings"
	"time"
)

func OnPrivMsg(client *rpc2.Client, line *irc.Line, reply *bool) error {
	var msg string
	channel, text := line.Args[0], line.Args[1]
	if text == "WHAT!" {
		msg = "Daaanger zooone"
	} else if text == "im coming" {
		msg = "\002\0032ARCHER \003\002 Phrasing!"
	} else {
		return nil
	}
	client.Call("privmsg", &bot.PrivMsgArgs{channel, msg}, nil)
	return nil
}

type boilerplate_conf struct {
	Channels []string `json:"channels"`
	BotHost  string   `json:"bot_host"`
}

func main() {
	conf := &boilerplate_conf{}
	cfg.NewConfig(conf, "archer.json")

	// RPC
	conn, _ := net.Dial("tcp", conf.BotHost)
	c := rpc2.NewClient(conn)
	go c.Run()
	// register privmsg
	c.Handle("privmsg", OnPrivMsg)
	c.Call("register", struct{}{}, nil)
	for _, channel := range conf.Channels {
		c.Call("join", channel, nil)
	}

	Quotes := []string{
		"\002\0036KRIEGER\003\002 Smoke bomb!",
		"\002\0034CHERYL \003\002 It's just like my birthday all over again. Remember?\n\002\0033PAM    \003\002 Nooo!\n\002\0034CHERYL \003\002 Yeah. Because you weren't invited",
		"\002\0032ARCHER \003\002 Yeah, 'child-murderer' sholdn't be hyphenated. That makes it seem like he's a murderer who's also a child",
		"\002\0033PAM    \003\002 Sploosh!",
		"\002\0037LANA   \003\002 Neewp!",
		"\002\0037LANA   \003\002 Nooope!",
		"\002\0032ARCHER \003\002 MULATTO BUTTS",
		"\002\0032ARCHER \003\002 Mmmahp",
		"\002\0032ARCHER \003\002 Maaaaahp",
		"\002\0034CHERYL \003\002 Polo",
		"\002\00311GILETTE\003\002 Get off! Clamydiot.\n\002\0034CHERYL \003\002 Oh I get it. Because of the Chlamidiya. Oh, and I'm an idiot",
		"\002\0034CHERYL \003\002 Oh my God, just like that old, gypsy woman said",
		"\002\0032ARCHER \003\002 You say that",
		"\002\0032ARCHER \003\002 Oh, okay, I guess just pout",
		"\002\0034CHERYL \003\002 Mopeds are fun, but you don't want your buddies to see you riding one",
		"\002\0032ARCHER \003\002 Relax, it's North Korea. The nation-state equivalent of the short bus",
		"\002\0032ARCHER \003\002 Eh, little column A, little column B",
		"\002\0033PAM    \003\002 Aww man, did I miss it?!\n\002\0034CHERYL \003\002 Oh my god, the toilet!?",
		"\002\0032ARCHER \003\002 Just the tip",
	}
	go func() {
		for {
			rand.Seed(time.Now().Unix())
			time.Sleep(time.Duration(rand.Intn(1440-180)+180) * time.Second)
			for _, s := range strings.Split(Quotes[rand.Intn(len(Quotes))], "\n") {
				for _, channel := range conf.Channels {
					c.Call("privmsg", &bot.PrivMsgArgs{channel, s}, nil)
				}
			}
		}
	}()

	// daemon
	forever := make(chan bool)
	<-forever
}
