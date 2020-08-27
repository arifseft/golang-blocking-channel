package main

import (
	"google.golang.org/api/iterator"
	"cloud.google.com/go/datastore"
	"context"
	"log"
)



func getAnswers(client *datastore.Client) <- chan Answer {
	answers := make(chan Answer,300)
	go func() {
		query := datastore.NewQuery("answers")
		for {
			var encontrados int
			t := client.Run(context.Background(),query)
			for answer := range getInCursor(t) {
				encontrados ++
				answers <- answer
			}
			if encontrados == 0 {
				break
			}
			cursor,err := t.Cursor()
			checkError(err)
			query = query.Start(cursor)
		}
		close(answers)
	}()
	return answers
}

func getInCursor(t *datastore.Iterator) <- chan Answer {
	answers := make(chan Answer)
	go func() {
		for {
			var answer Answer
			_, err := t.Next(&answer)
			if err == iterator.Done {
				break
			}
			if err != nil {
				checkError(err)
				break
			}
			answers <- answer
		}
		close(answers)
	}()
	return answers
}

func updateGrades(grades <- chan Grade, client *datastore.Client) {
	for grade := range grades {
		var studentEntity Student
		err := client.Get(context.Background(), datastore.NameKey("students", grade.Student, nil), &studentEntity)
		checkError(err)
		studentEntity.Grade = grade.Value
		_,err = client.Put(context.Background(),datastore.NameKey("students", grade.Student, nil),&studentEntity)
		checkError(err)
		log.Println(studentEntity)
	}
}
