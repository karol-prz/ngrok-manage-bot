package main

import (
	"time"
	"github.com/bwmarrin/discordgo"
	"strings"
	"os/exec"
	"fmt"
)

func checkIfNgrokRunning() (bool, string){
	cmd := "ps -eo pid,comm | grep -E 'ngrok$'"
	out := runCmd(cmd)
	result := strings.Fields(string(out))
	if len(result) < 2 {
		return false, ""
	}
	return result[1] == "ngrok", result[0]
}

func startNgrok(id string, s *discordgo.Session){
	stopNgrok(id, s)
	cmd := "nohup " + startCmd + " > /dev/null & disown"
	startACmd(cmd)
	s.ChannelMessageSend(id, "Started ngrok. Nice!")
	time.Sleep(time.Second * 2)
	getNgrokPort(id, s)
}

func stopNgrok(id string, s *discordgo.Session){
	running, pid := checkIfNgrokRunning()
	if running{
		cmd := "kill " + pid
		runCmd(cmd)
		s.ChannelMessageSend(id, "Stopped ngrok. Nice!")
	} 
}

func getNgrokPort(id string, s *discordgo.Session) {
	if running, _ := checkIfNgrokRunning(); !running{
		startNgrok(id, s)
	}
	cmd := `curl -s http://127.0.0.1:4040/inspect/http | grep -oE 'tcp://0.tcp.eu.ngrok.io:[0-9]{4,6}' | grep -Eo '[0-9]{4,6}'`
	r := runCmd(cmd)
	s.ChannelMessageSend(id, "Ngrok port is " + string(r))
	s.ChannelMessageSend(id, "ssh myuser@0.tcp.eu.ngrok.io -p" + string(r))
}

func runCmd(cmd string) []byte{
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil{
		if err.Error() == "exit status 1"{
			return ([]byte)("")
		}
		fmt.Printf("Error running cmd %s.\n Error %s\n", cmd, err.Error())
	}
	return out
}

func startACmd(cmd string){
	exec.Command("bash", "-c", cmd).Start()
}