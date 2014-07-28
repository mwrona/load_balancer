package utils

import (
	//"time"
	//"fmt"
	"log"
)

func Log(who, message string) {
	log.Printf("%s : %s", who, message)
	//fmt.Println("[", time.Now(), "]\t", who, ":", message)
}
