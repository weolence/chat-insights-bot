package main

import (
	"fmt"
	"main/internal/insightsbot"
)

func main() {
	insightsBot, err := insightsbot.NewInsightsBot()
	if err != nil {
		fmt.Println(err)
		return
	}
	insightsBot.Run()
}
