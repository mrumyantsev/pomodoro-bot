package main

import "github.com/mrumyantsev/pomodoro-bot/internal/app/pomd"

func main() {
	app := pomd.New()

	app.Run()
}
