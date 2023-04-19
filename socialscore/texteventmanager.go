package socialscore

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type TriggerText struct {
	Text   string `json:"text"`
	Points int    `json:"points"`
}

var bigtext = make(map[string]TriggerText)

func MessageEvaluatinator(s *discordgo.Session, m *discordgo.Message) {
	isEvent, points := checkPointEvent(m.Content)

	if !isEvent {
		return
	}

	updatePoints(m.Author.ID, points)
}

func checkPointEvent(message string) (bool, int) {
	value, exists := bigtext[message]

	if exists {
		return true, value.Points
	} else {
		return false, 0
	}
}

func ReadMatchData() {
	mySlice := make([]TriggerText, 0)

	// Read the file contents
	fileData, err := os.ReadFile("socialscore\\wordstocheck.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a slice of Person structs
	err = json.Unmarshal(fileData, &mySlice)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range mySlice {
		bigtext[obj.Text] = obj
	}
}
