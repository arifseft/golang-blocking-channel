package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"os"
	"log"
	"time"
)

type Answer struct {
	OptionChosen  		string
	QuestionNumber 		int
	Student				string
}
type Student struct {
	Id		string
	Name 	string
	Grade	int
}

type Template struct {
	QuestionNumber int
	RightChoice    string
	Value		   int
}


func main() {
	startTime := time.Now()
	client, err := datastore.NewClient(context.Background(),os.Getenv("PROJECT_ID"))
	checkError(err)
	answers := getAnswers(err, client)
	grades := calcGrades(answers, client)
	updateGrades(grades, client)
	endTime := time.Now()
	log.Println("duration:", endTime.Sub(startTime).Seconds())
}



func updateGrades(grades map[string]int, client *datastore.Client) {
	var students []*Student
	var keys []*datastore.Key
	for student, grade := range grades {
		var studentEntity Student
		err := client.Get(context.Background(), datastore.NameKey("students", student, nil), &studentEntity)
		checkError(err)
		studentEntity.Grade = grade
		keys = append(keys, datastore.NameKey("students", student, nil))
		students = append(students, &studentEntity)
		log.Println(studentEntity)
	}
	if len(students) > 0 {
		_,err := client.PutMulti(context.Background(),keys,students)
		checkError(err)
	}
}

func calcGrades(answers []Answer, client *datastore.Client) map[string]int {
	grades := make(map[string]int)
	for _, answer := range answers {
		var template Template
		err := client.Get(context.Background(), datastore.IDKey("templates", int64(answer.QuestionNumber), nil), &template)
		checkError(err)
		if answer.OptionChosen == template.RightChoice {
			grades[answer.Student] = grades[answer.Student] + template.Value
		}
	}
	return grades
}

func getAnswers(err error, client *datastore.Client) []Answer {
	var answers []Answer
	_, err = client.GetAll(context.Background(), datastore.NewQuery("answers"), &answers)
	checkError(err)
	return answers
}


func checkError(err error){
	if err != nil {
		log.Fatal(err)
	}
}