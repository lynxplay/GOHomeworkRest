package main

import (
	"database/sql"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/mattn/go-sqlite3"
)

//DatabaseConnection has been outsourced to the sql file as the type is handled in here
type DatabaseConnection struct {
	connection *sql.DB
}

//Open has been outsourced as it handles DatabaseConnection
func (d *DatabaseConnection) Open() error {
	return d.connection.Ping()
}

//Close has been outsourced as it handles DatabaseConnection
func (d *DatabaseConnection) Close() error {
	return d.connection.Close()
}

//GetConnection has been outsourced as it handles DatabaseConnection
func (d *DatabaseConnection) GetConnection() *sql.DB {
	return d.connection
}

//ExecuteSQLScript was outsourced
func (d *DatabaseConnection) ExecuteSQLScript(fileName string) error {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := d.connection.Exec(request)
		if err != nil {
			return err
		}
	}

	return nil
}

func openMysqlConnection(host string, port int, username string, password string, database string) (*DatabaseConnection, error) {
	connection, e := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+strconv.Itoa(port)+")/"+database)
	if e != nil {
		return nil, e
	}

	wrapper := DatabaseConnection{connection: connection}
	connectionError := wrapper.Open()

	return &wrapper, connectionError
}

func openSQLLiteConnection(databaseName string) (*DatabaseConnection, error) {
	connection, e := sql.Open("sqlite3", databaseName)
	if e != nil {
		return nil, e
	}

	wrapper := DatabaseConnection{connection: connection}
	connectionError := wrapper.Open()

	return &wrapper, connectionError
}

func loadAccount(connetion *DatabaseConnection, username string, password string) *AccountData {
	userResultSet, e := connection.GetConnection().Query("SELECT * FROM users WHERE username=? AND password=?", username, password)
	check(e)
	defer userResultSet.Close()

	if !userResultSet.Next() {
		return nil
	}
	var id int
	userResultSet.Scan(&username, &password, &id)

	classes := loadClasses(connection, id)

	return &AccountData{
		Username:      username,
		ID:            id,
		LoginTime:     time.Now(),
		ClassArray:    classes,
		HomeworkArray: loadHomework(connection, id, classes),
	}
}

func loadClasses(connection *DatabaseConnection, playerID int) []*Class {
	classResultSet, err := connection.GetConnection().Query("SELECT * FROM classes WHERE player_id = ?", playerID)
	check(err)
	defer classResultSet.Close()

	var (
		PlayerID int
		ID       int
		Title    string
		Icon     string
		classes  []*Class
	)

	for classResultSet.Next() {
		classResultSet.Scan(&PlayerID, &ID, &Title, &Icon)
		classes = append(classes, &Class{
			ID:    ID,
			Title: Title,
			Icon:  Icon,
		})
	}

	return classes
}

func loadHomework(connection *DatabaseConnection, playerId int, classes []*Class) []*Homework {
	homeworkResultSet, err := connection.GetConnection().Query("SELECT * FROM homework WHERE player_id = ?", playerId)
	check(err)
	defer homeworkResultSet.Close()

	var (
		PlayerID    int
		ClassID     int
		HomeworkID  int
		Description string
		DueDay      mysql.NullTime
		homework    []*Homework
	)

	for homeworkResultSet.Next() {
		homeworkResultSet.Scan(&PlayerID, &ClassID, &HomeworkID, &Description, &DueDay)

		if len(classes) < ClassID {
			continue
		}

		homework = append(homework, &Homework{
			Class:       classes[ClassID],
			Description: Description,
			DueDay:      &DueDay.Time,
		})
	}

	return homework
}
