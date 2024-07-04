package main

import (
	"fmt"
	fun "go-prof-sprint-1/internal/functions"
)

// глобальные переменные флагов
var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	fmt.Printf("version=%s, date=%s, commit=%ss\n", BuildVersion, BuildDate, BuildCommit)
	if err := fun.Run(); err != nil {
		panic(err)
	}

}
