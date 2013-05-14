package main

import (
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
        if (v == value) {
            return true
        }
    }
    return false
}

type Command struct {
    Name string         // Name for the Command
    Args []string       // Any arguments to pass to function we call.
    Source string       // Who requested it?
}

// Pulls out the commands and arguments out of a message directed toward the bot
// the first parameter in the message is the bots name, so it's ignored here.
func extractCommand(msg string, source string) (command Command) {
    subStrings := strings.Split(msg, " ")
    commandName := "Unknown"
    if len(subStrings) > 1 {
        commandName = strings.Title(subStrings[1])
    }
    args := []string{}
    if len(subStrings) > 2 {
        args = subStrings[2:]
    }
    return Command{commandName, args, source}
}

func saveUserInput(nick string, msg string, userInput map[string][]string) {
    userInput[nick] = append(userInput[nick], msg)
    fmt.Printf("Looking at User values: \n%#v\n", userInput[nick])
}


func main() {

    myName := "godwit"
    myHome := "#gobots"
    myChannels := []string{}
    myServer := "irc.freenode.net:6667"
    userInput := map[string][]string{}

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
        myChannels = append(myChannels, curChan)
        ircobj.Privmsg(curChan, "Hey! What'd I miss?")
    })

    // It turns out that all messages are PRIVMSG
    // Figure out if there are commands involved and figure out the scope of the
    // message.
    ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
        if strings.HasPrefix(event.Message, myName) {
            myChan := event.Arguments[0]
            command := extractCommand(event.Message, event.Nick)

            switch command.Name {

                case "Unknown":
                    ircobj.Privmsg(myChan, fmt.Sprintf("Sorry, I don't understand you, %s", event.Nick))

                case "Impersonate":
                    target := event.Nick
                    // Use the initiator if no other name is specified.
                    if len(command.Args) > 0 {
                        target = command.Args[0]
                    }
                    c := NewChain(2)
                    c.Build(userInput, command.Args[0])
                    ircobj.Privmsg(myChan, c.Generate(10))

                case "Read":
                    ircobj.Privmsg(myChan, fmt.Sprintf("%s: Make me", event.Nick))

                case "Thanks":
                    target := ""
                    if len(command.Args) > 0 {
                        target = command.Args[0]
                    }

                    ircobj.Privmsg(myChan, fmt.Sprintf("!m %s", target))
            }
        }

        // Save this user's data so we can impersonate them later!
        // But only if this isn't a command for this bot.
        saveUserInput(event.Nick, event.Message, userInput)


        fmt.Printf("%#v\n", event)
    })

    ircobj.AddCallback("INVITE", func(event *irc.Event) {
        //fmt.Printf("%#v\n", event)
        ircobj.Join(event.Message)
    })

    ircobj.Loop()

}

