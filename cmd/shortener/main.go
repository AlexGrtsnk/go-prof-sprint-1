package main

import (
	fun "go-prof-sprint-1/cmd/functions"
)

func main() {
	if err := fun.Run(); err != nil {
		panic(err)
	}

}
