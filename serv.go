package main

import (
	API "KoKo/site_web/rsc"
	"fmt"
)

func main() {
	API.Init()
	API.RUN()
	err := API.InitDB()
	if err != nil {
		fmt.Println(err)
	}
}
