package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/datastore"
)

type QuestionOption string

const (
	A = "A"
	B = "B"
	C = "C"
	D = "D"
)

type Student struct {
	Id    string
	Name  string
	Grade int
}

type Template struct {
	QuestionNumber int
	RightChoice    QuestionOption
	Value          int
}

type Answer struct {
	OptionChosen   QuestionOption
	QuestionNumber int
	Student        string
}

func readStudents() <-chan Student {
	students := make(chan Student)
	go func() {
		file, err := os.Open("students.csv")
		if err != nil {
			close(students)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			studentData := strings.Split(scanner.Text(), ";")
			students <- Student{Id: studentData[0], Name: studentData[1]}
		}
		close(students)
	}()
	return students
}

func readTemplate() <-chan Template {
	templateAnwsers := make(chan Template)
	go func() {
		file, err := os.Open("template.csv")
		if err != nil {
			close(templateAnwsers)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := strings.Split(scanner.Text(), ";")
			questionNumber, _ := strconv.Atoi(data[0])
			questionValue, _ := strconv.Atoi(data[2])
			templateAnwsers <- Template{
				QuestionNumber: questionNumber,
				RightChoice:    getChoice(data[1]),
				Value:          questionValue}
		}
		close(templateAnwsers)
	}()
	return templateAnwsers
}

func readAnswsers() <-chan Answer {
	answers := make(chan Answer)
	go func() {
		file, err := os.Open("answers.csv")
		if err != nil {
			close(answers)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := strings.Split(scanner.Text(), ";")
			questionNumber, _ := strconv.Atoi(data[1])
			answers <- Answer{
				OptionChosen:   getChoice(data[0]),
				QuestionNumber: questionNumber,
				Student:        data[2]}
		}
		close(answers)
	}()
	return answers
}

func persistStudent(students <-chan Student) {
	client, err := datastore.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	checkError(err)
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 20; i++ {
		persistStudents(students, client, &wg)
	}
	wg.Wait()

}
func persistStudents(students <-chan Student, client *datastore.Client, wg *sync.WaitGroup) {
	for student := range students {
		log.Print(student)
		_, err := client.Put(context.Background(), datastore.NameKey("students", student.Id, nil), &student)
		checkError(err)
	}
	wg.Done()
}

func persistTemplates(templates <-chan Template) {
	client, err := datastore.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	checkError(err)
	for template := range templates {
		_, err := client.Put(context.Background(), datastore.IDKey("templates", int64(template.QuestionNumber), nil), &template)
		checkError(err)
	}
}

func persistAnswers(answers <-chan Answer) {
	client, err := datastore.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	checkError(err)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		insertAnswer(answers, client, &wg)
	}
	wg.Wait()
}
func insertAnswer(answers <-chan Answer, client *datastore.Client, wg *sync.WaitGroup) {
	for answer := range answers {
		_, err := client.Put(context.Background(), datastore.IDKey("answers", 0, nil), &answer)
		checkError(err)
		log.Print(answer)
	}
	wg.Done()
}

func main() {
	students := readStudents()
	templates := readTemplate()
	answers := readAnswsers()
	persistStudent(students)
	persistTemplates(templates)
	persistAnswers(answers)
}

func getChoice(choice string) QuestionOption {
	switch choice {
	case "A":
		return A
	case "B":
		return B
	case "C":
		return C
	default:
		return D
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
