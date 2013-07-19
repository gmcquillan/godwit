package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent"
	"math/rand"
	"strings"
	"time"
)

/*
* Written for GO 1.1 +
* Godwit - a small, clever bird in Northeastern America.
*
* Super helpful for deciphering IRC codes:
* https://www.alien.net.au/irc/irc2numerics.html
 */

// containsStr takes a slice and a string and returns true if the string
// exists in the slice.
func containsStr(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

var nick = flag.String("nick", "godwit", "nick of the bot")
var channel = flag.String("chan", "#gobots", "channel to join")
var server = flag.String("server", "irc.freenode.net:6667", "server and port to join")

func main() {
	flag.Parse()
	myName := *nick
	myHome := *channel
	myChannels := []string{}
	myServer := *server

	ircobj := irc.IRC(myName, myName)

	//ircobj.UseTLS = true //default is false
	//ircobj.TLSOptions //set ssl options
	//ircobj.Password = "[server password]"

	err := ircobj.Connect(myServer)
	if err != nil {
		fmt.Println("Connection Error")
		return
	}

	rand.Seed(time.Now().UnixNano())
	// Welcome message -- we've successfully connected!
	ircobj.AddCallback("001", func(e *irc.Event) { ircobj.Join(myHome) })

	ircobj.AddCallback("JOIN", func(event *irc.Event) {
		//fmt.Printf("%#v\n", event)
		curChan := event.Arguments[0]
		if event.Nick == myName {
			myChannels = append(myChannels, curChan)
			ircobj.Privmsg(curChan, "Hey! What'd I miss?")
		}
	})

	// It turns out that all messages are PRIVMSG
	// Figure out if there are commands involved and figure out the scope of the
	// message.
	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		if strings.HasPrefix(event.Message, myName) {
			myChan := event.Arguments[0]
			command := ExtractCommand(event.Message, event.Nick)
			response := RunCommand(command)
			if len(response) > 0 {
				ircobj.Privmsg(myChan, response)
			}
		}

		// Save this user's data so we can impersonate them later!
		// But only if this isn't a command for this bot.

		fmt.Printf("%#v\n", event)
	})

	ircobj.AddCallback("INVITE", func(event *irc.Event) {
		//fmt.Printf("%#v\n", event)
		ircobj.Join(event.Message)
	})

	ircobj.Loop()

}
