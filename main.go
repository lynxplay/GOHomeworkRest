package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.LoadHTMLGlob("resources/*.html")
	router.Static("/assets", "resources/assets")

	router.GET("/", func(context *gin.Context) {

		context.HTML(200, "index.html", gin.H{
			"payload": mapHomeworkToTimelineElements(&[]Homework{
				{Title: "Math Homework", Description: "What is 1+1"},
				{Title: "German Homework", Description: "Read some random books"},
				{Title: "DataStructure Homework", Description: "Work on that SQL thingy"},
			}),
		})
	})
	router.POST("/upload", func(context *gin.Context) {
		bytes, _ := context.GetRawData()
		fmt.Println(string(bytes[:]))
	})

	router.Run()
}

type Homework struct {
	Title       string
	Description string
}

type HomeworkTimelineElement struct {
	IsReverse   string
	HasNextNode string
	Title       string
	Description string
}

func mapHomeworkToTimelineElements(homework *[]Homework) []HomeworkTimelineElement {
	var timelineElements []HomeworkTimelineElement
	for index, homework := range *homework {
		timelineElements = append(timelineElements, createHomeworkTimelineObject(homework.Title, homework.Description, index%2 == 0))
	}
	timelineElements[len(timelineElements)-1].HasNextNode = ""
	return timelineElements
}

func createHomeworkTimelineObject(title string, description string, reverse bool) HomeworkTimelineElement {
	element := HomeworkTimelineElement{HasNextNode: "separline"}
	element.Title = title
	element.Description = description
	if reverse {
		element.IsReverse = "reverse"
	} else {
		element.IsReverse = ""
	}

	return element
}
