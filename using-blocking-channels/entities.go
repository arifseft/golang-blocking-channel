package main

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

type Grade struct {
	Student string
	Value int
}

