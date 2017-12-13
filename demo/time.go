package main

import (
	"time"
	"log"
	"fmt"
)

func main() {
	timeStr := "2017-10-12T15:43:26"
	t ,err := time.Parse("2006-01-02T15:04:05",timeStr)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("origin : ", t.String())
	fmt.Println("mm-dd-yyyy : ", t.Format("01-02-2006"))
	fmt.Println("yyyy-mm-dd : ", t.Format("2006-01-02"))
	fmt.Println("yyyy.mm.dd : ", t.Format("2006.01.02"))
	fmt.Println("yyyy-mm-dd HH:mm:ss : ", t.Format("2006-01-02 15:04:05"))
	fmt.Println("yyyy-mm-dd HH:mm:ss : ", t.Format("2006-01-02 15:04:05.000000"))
	fmt.Println("yyyy-mm-dd HH:mm:ss : ", t.Format("2006-01-02 15:04:05.000000000"))
}

