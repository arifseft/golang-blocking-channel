package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

func main() {
	startTime := time.Now()
	client, err := datastore.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	checkError(err)
	answers := getAnswers(client)
	grades := checkGrades(answers, client)
	//grades2 := calcGrades(answers, client)
	//grades3 := calcGrades(answers, client)
	//grades4 := calcGrades(answers, client)
	//grades5 := calcGrades(answers, client)
	gradesAcumulated := accumulateGrades(mergeGrades(grades))
	updateGrades(gradesAcumulated, client)
	endTime := time.Now()
	log.Println("duration:", endTime.Sub(startTime).Seconds())
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
