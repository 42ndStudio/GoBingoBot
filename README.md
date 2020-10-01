# GoBingoBot
Telegram Bingo Bot written in GoLang

Features | Funciones: 
- Generate Boards
  - Each has a unique per game ID 
  - PNG Image can be rendered for each board
  - Organizers can ask for board images like: All Boards, Unsold Boards, Specific Board

- Multiple Game Modes
  - Horizontal / Vertical Line
  - C
  - O
  - Diagonal / & \

- Can handle Multiple Games with Separate Boards

- Keeps Track of Balots
 - Organizers can be notified of winning boards
 - Organizers can ask if Board has Bingo
 - Can Relay messages to Group Chat

## Run
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=gobingobot
export DB_USER=bingopher
export DB_PASS=secreto
export TG_KEY="secreto"
go run . 
```