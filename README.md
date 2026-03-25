# Learn Go by TDD: Poker App

This application is built by following the tutorials from [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests).

## Features
* Two applications: a command-line application and a web server.
* With the CLI, users can spin up the app locally.
* With the web server, navigate to `/game` to start playing.
* After either application is started, input the number of players and start the game.
* Every 5 + `[number of players]` minutes, the blind will increase by 100, triggering an alert.
* With the CLI, users can end the game by recording the winner with the format: `[Player name] wins`.
* With the web server, simply enter the player's name into the winner form.
* All data is stored locally in `game.db.json`.
