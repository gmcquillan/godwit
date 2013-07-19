package main

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"net/http"
	"strings"
)

type Command struct {
	Name   string   // Name for the Command.
	Args   []string // Any arguments to pass to function we call.
	Source string   // Who requested it?
}

// Pulls out the commands and arguments out of a message directed toward the bot
// the first parameter in the message is the bots name, so it's ignored here.
func ExtractCommand(msg string, source string) (command Command) {
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

func summarizeWebsite(htmlBody string) string {
	strippedData := sanitize.HTML(htmlBody)
	return fmt.Sprintf("%s.", strings.Join(strings.Split(strippedData, " ")[40:50], " "))
}

// Process command, return message for sending to channel.
func RunCommand(command Command) string {
	response := ""
	switch command.Name {

	case "Unknown":
		response = fmt.Sprintf("Sorry I don't understand you, %s", command.Source)

	case "Read":
		if len(command.Args) < 1 {
			response = fmt.Sprintf("Eh, that's not a URL, %s", command.Source)
			break
		}
		url := strings.ToLower(command.Args[0])
		if !strings.HasPrefix(url, "http://") {
			url = fmt.Sprintf("http://%s", url)
		}
		resp, _ := http.Get(url)
		htmlBody, _ := ioutil.ReadAll(resp.Body)
        if htmlBody !=nil{
			response = summarizeWebsite(string(htmlBody))
		} else {
			response = "Didn't quite get that; try again?"
		}

	case "Thanks":
		target := ""
		if len(command.Args) > 0 {
			target = command.Args[0]
		}

		response = fmt.Sprintf("%s, you're welcome.", target)
	}

	return response
}
