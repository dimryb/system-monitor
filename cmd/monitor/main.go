package main

import (
	"fmt"
	"github.com/dimryb/system-monitor/internal/app"
)

func main() {
	application := app.NewApp()

	fmt.Println(application)

	fmt.Println("Starting system-monitor...")
}
