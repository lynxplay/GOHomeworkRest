package main

import (
	"time"
)

//AccountData was
type AccountData struct {
	Username  string
	ID        int
	LoginTime time.Time

	ClassArray    []*Class
	HomeworkArray []*Homework
}

//Class was outsourced
type Class struct {
	ID    int
	Title string
	Icon  string
}

//Homework was outsourced
type Homework struct {
	Class       *Class
	Description string
	DueDay      *time.Time
}

func (a *AccountData) isSessionValid() bool {
	return time.Now().After(a.LoginTime)
}
