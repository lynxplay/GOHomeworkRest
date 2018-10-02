package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var connection *DatabaseConnection
var userCache = make(map[int]*AccountData)

func main() {
	router := gin.New()
	restServerConfiguration, e := LoadServerConfiguration("configuration.json", RestServerConfiguration{
		SessionKey:       "superSecretSessionKey",
		Database:         "gin-rest",
		DatabaseHost:     "localhost",
		DatabasePassword: "password",
		DatabaseUsername: "username",
		DatabasePort:     3306,
		DatabaseFileName: "./gin-rest-database.db",
	})
	check(e)

	fmt.Println("[GIN-INFO] Connecting to remote sql server...")
	connection, e = openMysqlConnection(restServerConfiguration.DatabaseHost,
		restServerConfiguration.DatabasePort,
		restServerConfiguration.DatabaseUsername,
		restServerConfiguration.DatabasePassword,
		restServerConfiguration.Database)
	if e != nil {
		fmt.Println("[GIN-ERROR] Could not connect to remote sql server!!")
		fmt.Println("[GIN-INFO] Opening SQLLite database...")

		connection, e = openSQLLiteConnection(restServerConfiguration.DatabaseFileName)
		if e != nil {
			fmt.Println("[GIN-ERROR] Could not open connection to SQLLite file")
			panic(e)
		}
	} else {
		fmt.Println("[GIN-INFO] Connected to sql server!")
	}

	fmt.Println("[GIN-INFO] Executing sql setup script...")
	setupError := connection.ExecuteSQLScript("./sql-setup.sql")
	if setupError != nil {
		fmt.Println("[GIN-ERROR] Could not execute setup script")
		panic(setupError)
	} else {
		fmt.Println("[GIN-INFO] Finished database setup")
	}

	store := memstore.NewStore([]byte(restServerConfiguration.SessionKey))
	fmt.Println("[GIN-INFO] Using memstore key", restServerConfiguration.SessionKey)

	router.LoadHTMLGlob("resources/*.html")
	router.Static("/assets", "resources/assets")
	router.Use(sessions.Sessions("gin-rest-sessions", store))

	router.GET("/", func(context *gin.Context) {
		session := sessions.Default(context)
		account := getAccountData(session)
		if account == nil || !account.isSessionValid() {
			context.Redirect(302, "/login.html")
			return
		}

		homework := mapHomeworkToTimelineElements(account.HomeworkArray)
		sort.Slice(homework, func(h1, h2 int) bool {
			return homework[h1].DueDateRaw.Before(*homework[h2].DueDateRaw)
		})

		context.HTML(200, "index.html", gin.H{
			"payload":  mapHomeworkToTimelineElements(account.HomeworkArray),
			"Username": account.Username,
		})
	})

	router.GET("/login.html", func(context *gin.Context) {
		context.HTML(200, "login.html", gin.H{})
	})

	router.POST("/login.html", func(context *gin.Context) {
		context.Request.ParseForm()

		username := context.Request.FormValue("signin-email")
		password := context.Request.FormValue("signin-password")

		account := loadAccount(connection, username, password)
		if account == nil {
			context.JSON(403, "Could not find the user account")
		} else {
			session := sessions.Default(context)
			session.Set("userID", account.ID)
			userCache[account.ID] = account
			session.Save()

			context.Redirect(302, "/")
		}
	})

	router.Run()
}

type HomeworkTimelineElement struct {
	IsReverse   string
	HasNextNode string
	Title       string
	Description string
	DueDate     string
	DueDateRaw  *time.Time
	Icon        string
}

func mapHomeworkToTimelineElements(homework []*Homework) []*HomeworkTimelineElement {
	var timelineElements []*HomeworkTimelineElement
	for index, homework := range homework {
		timelineElements = append(timelineElements,
			createHomeworkTimelineObject(homework, index%2 == 0))
	}

	length := len(timelineElements)
	if length > 0 {
		timelineElements[length-1].HasNextNode = ""
	}

	return timelineElements
}

func createHomeworkTimelineObject(homework *Homework, reverse bool) *HomeworkTimelineElement {
	element := HomeworkTimelineElement{HasNextNode: "separline"}
	element.Title = homework.Class.Title
	element.Description = homework.Description
	element.Icon = homework.Class.Icon
	element.DueDate = strconv.Itoa(homework.DueDay.Day()) + "." + strconv.Itoa(int(homework.DueDay.Month())) + "." + strconv.Itoa(homework.DueDay.Year())
	element.DueDateRaw = homework.DueDay
	if reverse {
		element.IsReverse = "reverse"
	} else {
		element.IsReverse = ""
	}

	return &element
}

func getAccountData(s sessions.Session) *AccountData {
	if s == nil {
		return nil
	}

	userID := s.Get("userID")
	if userID == nil {
		return nil
	}

	if value, castable := userID.(int); castable {
		return userCache[value]
	}
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
