package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type DatabaseConnection struct {
	connection *sql.DB
}

func (d *DatabaseConnection) Open() error {
	return d.connection.Ping()
}

func (d *DatabaseConnection) Close() error {
	return d.connection.Close()
}

func (d *DatabaseConnection) GetConnection() *sql.DB {
	return d.connection
}

func openConnection(host string, port int, username string, password string, database string) *DatabaseConnection {
	connection, e := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+strconv.Itoa(port)+")/"+database)
	check(e)

	wrapper := DatabaseConnection{connection: connection}
	wrapper.Open()

	return &wrapper
}

func loadAccount(connetion *DatabaseConnection, username string, password string) (*AccountData) {
	userResultSet, e := connection.GetConnection().Query("SELECT * FROM users WHERE username=? AND password=?", username, password)
	check(e)

	if !userResultSet.Next() {
		return nil
	}
	var id int
	userResultSet.Scan(&username, &password, &id)

	classes := loadClasses(connection, id)

	return &AccountData{
		Username:      username,
		Id:            id,
		LoginTime:     time.Now(),
		ClassArray:    classes,
		HomeworkArray: loadHomework(connection, id, classes),
	}
}

func loadClasses(connection *DatabaseConnection, playerId int) []*Class {
	classResultSet, err := connection.GetConnection().Query("SELECT * FROM classes WHERE playerID = ?", playerId)
	check(err)

	var (
		PlayerID    int
		Id      int
		Title   string
		Icon    string
		classes []*Class
	)

	for classResultSet.Next() {
		classResultSet.Scan(&PlayerID, &Id, &Title, &Icon)
		classes = append(classes, &Class{
			Id:    Id,
			Title: Title,
			Icon:  Icon,
		})
	}

	return classes
}

func loadHomework(connection *DatabaseConnection, playerId int, classes []*Class) []*Homework {
	homeworkResultSet, err := connection.GetConnection().Query("SELECT * FROM homework WHERE playerID = ?", playerId)
	check(err)

	var (
		PlayerID    int
		ClassId     int
		Description string
		DueDay      mysql.NullTime
		homework    []*Homework
	)

	for homeworkResultSet.Next() {
		homeworkResultSet.Scan(&PlayerID, &ClassId, &Description, &DueDay)
		homework = append(homework, &Homework{
			Class:       classes[ClassId - 1],
			Description: Description,
			DueDay:      &DueDay.Time,
		})
	}

	return homework
}
