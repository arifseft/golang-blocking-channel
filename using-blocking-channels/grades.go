package main

import ( "sync"
	"cloud.google.com/go/datastore"
	"context"
)

func checkGrades(answers <- chan Answer, client *datastore.Client) <- chan Grade {
	grades := make(chan Grade)
	go func() {
		for answer := range answers {
			var template Template
			err := client.Get(context.Background(), datastore.IDKey("templates", int64(answer.QuestionNumber), nil), &template)
			checkError(err)
			grade := Grade{Student:answer.Student}
			if answer.OptionChosen == template.RightChoice {
				grade.Value=  template.Value
			}
			grades <- grade
		}
		close(grades)
	}()
	return grades
}


func mergeGrades(gradeChannels ...<- chan Grade) <- chan Grade {
	grades := make(chan Grade)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(gradeChannels))

		for _,gradeChannel := range gradeChannels {
			go func(channel <- chan Grade) {
				for grade := range channel {
					grades <- grade
				}
				wg.Done()
			}(gradeChannel)
		}

		wg.Wait()
		close(grades)
	}()


	return grades
}



func accumulateGrades(grades <- chan Grade) <- chan Grade {
	gradesAcumulated := make(chan Grade)
	go func() {
		studentProcessedTimes := make(map[string]int)
		sumStudentGrades := make(map[string]int)
		for grade := range grades {
			studentProcessedTimes[grade.Student] = studentProcessedTimes[grade.Student] + 1
			sumStudentGrades[grade.Student] = sumStudentGrades[grade.Student] + grade.Value
			if studentProcessedTimes[grade.Student] == 10 {
				gradesAcumulated <- Grade{grade.Student, sumStudentGrades[grade.Student]}
				delete(studentProcessedTimes,grade.Student)
				delete(sumStudentGrades,grade.Student)
			}
		}
		for student,grade := range sumStudentGrades {
			gradesAcumulated <- Grade{student, grade}
		}
		close(gradesAcumulated)
	}()

	return gradesAcumulated
}
