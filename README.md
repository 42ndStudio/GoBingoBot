# GoBingoBot
Telegram Bingo Bot written in GoLang

Features | Funciones: 
- [x] Every Game can have multiple Organizers (Organizers change game mode, and register drawn ballots)
  - [x] Keep Track of sold boards

- [x] Generate Boards 
  - [x] Each has a unique per game ID 
  - [x] Unique boards mode
  - [x] PNG Image can be rendered for each board
  - [ ] Organizers can ask for boards like: All Boards, Unsold Boards, Specific Board

- [x] Multiple Game Modes (Organizers can change it)
  - [x] Horizontal / Vertical Line
  - [x] C
  - [x] O
  - [ ] U
  - [x] N
  - [x] Diagonal / & \
  - [x] All

- [x] Can handle Multiple Games with Separate Boards

- [x] Keeps Track of Balots
 - [x] Organizers are notified on number of winners as soon as they draw a ballot
 - [x] Organizers can ask if Board has Bingo
 - [ ] Organizers can ask for drawn balots
 - [ ] Can Relay messages to Group Chat

## Font
[MoonGet](https://www.dafont.com/moon-get.font)

Important: because we use a CGO enabled package, we are required to set the environment variable CGO_ENABLED=1 and have a gcc compile present within your path.
## Run
```bash
# If using a MySQL DB (TODO (for now using only sqlite))
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=gobingobot
export DB_USER=bingopher
export DB_PASS=secreto

# If you want to run a Telegram Bingo Bot
export TG_KEY="1295700987:AAG2Mr6Z77enrTCUMBZ7oNdPV-QBpqE2DJw"
export MASTER_ID=596025632  #MASTER_ID is the manager's TelegramID

export CGO_ENABLED=1
go run . 
```

