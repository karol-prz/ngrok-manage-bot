# How to Setup
- Open settings.json
- Put your bot token for the BotID string
- Set NgrokCmd to the ngrok exec file. So if you copy that into this folder NgrokCmd: "./ngrok"
- If you leave AdminID list empty, anyone can manage ngrok
- But you can use !id command to get your id and paste it in there as a string

# How to Start
- cd into this repository
- run ./start-ngrok-bot-background

# How to exit
- run ps -eo pid,comm | grep ngrok
- The pid of processes is the first column
- run kill pid
- pid can be a list of pids seperated by a space
