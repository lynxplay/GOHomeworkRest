package main

import (
	"time"
)

type AccountData struct {
	Username  string
	ID        int
	LoginTime time.Time

	ClassArray    []*Class
	HomeworkArray []*Homework
}

type Class struct {
	ID    int
	Title string
	Icon  string
}

type Homework struct {
	Class       *Class
	Description string
	DueDay      *time.Time
}

func (a *AccountData) isSessionValid() bool {
	return time.Now().After(a.LoginTime)
}
