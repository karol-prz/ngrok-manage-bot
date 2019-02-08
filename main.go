package main

import (
	"strings"
	"log"
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Settings struct{
	NgrokCmd string
	BotID string
	AdminID []string
}

var settings Settings
var startCmd string

func main(){
	readSettings()
	runDiscordBot()
}

func runDiscordBot(){
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + settings.BotID)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(dealWithMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Discord Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func dealWithMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if !strings.HasPrefix(m.Content, "!"){
		return
	}

	if m.Content == "!id"{
		s.ChannelMessageSend(m.ChannelID, m.Author.ID)
	}

	// Check id admin is sending messages
	// If admin id is empty allow
	valid := false
	for _, id := range(settings.AdminID){
		if id == m.Author.ID {
			valid = true
		}
	}
	if !valid && len(settings.AdminID) != 0{
		return
	}

	if m.Content == "!status"{
		running, _ := checkIfNgrokRunning()
		r := "off"
		if running {
			r = "running"
		}
		s.ChannelMessageSend(m.ChannelID, "Ngrok is " + r)
	}

	if m.Content == "!start" {
		startNgrok(m.ChannelID, s)
	}

	if m.Content == "!stop" {
		stopNgrok(m.ChannelID, s)
	}

	if m.Content == "!port" {
		getNgrokPort(m.ChannelID, s)
	}

	if m.Content == "!help" || m.Content == "!ngrok"{
		sendHelp(m.ChannelID, s)
	}

	if m.Content == "!stopbot"{
		s.Close()
	}
}

func sendHelp(id string, s *discordgo.Session){
	msg := "Ngrok allows ssh without port forwarding\nCommands:"
	msg += "\n\t!help or !ngrok to show help"
	msg += "\n\t!status checks if ngrok is running"
	msg += "\n\t!start starts ngrok and sends port number"
	msg += "\n\t!stop stops ngrok"
	msg += "\n\t!port sends port number"
	msg += "\n\t!id sends your user id"
	msg += "\n\t!stopbot Stops ME"
	s.ChannelMessageSend(id, msg)
}


func readSettings(){
	file, err := os.Open("settings.json")
	endOnError(err)
	byteValue, err := ioutil.ReadAll(file)
	endOnError(err)
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &settings)
	// fmt.Println("NgrokCmd", settings.NgrokCmd)
	// fmt.Println("BotID", settings.BotID)
	// fmt.Println("AdminID's", settings.AdminID)
	startCmd = settings.NgrokCmd + " tcp -log=stdout --region=eu 22"
}

func endOnError(err error){
	if err != nil{
	log.Fatal(err.Error())
	}
}

func printOnError(err error){
	if err != nil{
		fmt.Println(err.Error())
	}
}