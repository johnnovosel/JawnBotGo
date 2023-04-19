package socialscore

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

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

	writeResponse(s, m, points)
}

func checkPointEvent(message string) (bool, int) {
	value, exists := bigtext[strings.ToLower(message)]

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

func writeResponse(s *discordgo.Session, m *discordgo.Message, points int) {

	message := "<@" + m.Author.ID + ">"

	if points > 0 {
		message = message + " gained " + strconv.Itoa(points) + " social point for this sentence:\n\n" + m.Content
	} else {
		message = message + " lost " + strconv.Itoa(int(math.Abs(float64(points)))) + " social point for this sentence:\n\n" + m.Content
	}

	embed := &discordgo.MessageEmbed{
		Title:       "POINT EVENT",
		Description: message,
		Color:       0x00ff00, // Green
	}

	var err error

	_, err = s.ChannelMessageSendEmbed("168486914356281344", embed)
	if err != nil {
		log.Fatal(err)
	}
}
