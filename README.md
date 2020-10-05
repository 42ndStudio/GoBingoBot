# GoBingoBot
Telegram Bingo Bot written in GoLang

Features | Funciones: 
- [x] Every Game can have multiple Organizers (Organizers change game mode, and register drawn ballots)
  - [x] Keep Track of sold boards

- [x] Generate Boards 
  - [x] Each has a unique per game ID 
  - [ ] PNG Image can be rendered for each board
  - [ ] Organizers can ask for boards like: All Boards, Unsold Boards, Specific Board

- [ ] Multiple Game Modes (Organizers can change it)
  - [x] Horizontal / Vertical Line
  - [x] C
  - [x] O
  - [x] Diagonal / & \
  - [x] All

- [x] Can handle Multiple Games with Separate Boards

- [x] Keeps Track of Balots
 - [x] Organizers are notified on number of winners as soon as they draw a ballot
 - [x] Organizers can ask if Board has Bingo
 - [ ] Organizers can ask for drawn balots
 - [ ] Can Relay messages to Group Chat

## Run
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=gobingobot
export DB_USER=bingopher
export DB_PASS=secreto
export TG_KEY=secreto
export MASTER_ID=master_telegram_user_id  #MASTER_ID is the manager's TelegramID
go run . 
```
