package main

import (
	"KoKo/forum_API/route"
	API "KoKo/forum_API/rsc"
	"fmt"
)

func main() {
	route.RUN()
	err := API.InitDB()
	if err != nil {
		fmt.Println(err)
	}
}
