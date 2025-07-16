package main

import (
	"drillCore/internal/config"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
