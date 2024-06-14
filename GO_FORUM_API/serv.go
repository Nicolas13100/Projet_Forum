package main

import (
	API "KoKo/forum_API/rsc"
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
