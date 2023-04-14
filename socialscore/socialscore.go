package socialscore

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Data struct {
	DiscordID string `json:"discordID"`
	Score     int    `json:"score"`
	Name      string `json:"name"`
}

var bigdata = make(map[string]Data)

func ReadUserFile() {
	mySlice := make([]Data, 0)

	// Read the file contents
	fileData, err := os.ReadFile("socialscore\\db.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a slice of Person structs
	err = json.Unmarshal(fileData, &mySlice)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range mySlice {
		bigdata[obj.DiscordID] = obj
	}
}

func WriteJSONFile() error {

	var userSlice []Data
	for _, user := range bigdata {
		userSlice = append(userSlice, user)
	}

	// Marshal struct to JSON
	jsonData, err := json.MarshalIndent(userSlice, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to file
	err = os.WriteFile("socialscore\\db.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddUser(userID string, name string) {
	user := &Data{
		DiscordID: userID,
		Score:     100,
		Name:      name,
	}

	bigdata[user.DiscordID] = *user
}

func BulkAddUser(members []*discordgo.Member) {

	for _, member := range members {

		//check if the user already exists
		doesExists := AlreadyExists(member.User.ID)

		if doesExists {
			break
		} else {
			AddUser(member.User.ID, member.User.Username)
		}
	}

	WriteJSONFile()
}

func AlreadyExists(userID string) bool {
	_, exists := bigdata[userID]

	return exists
}

func updatePoints(userID string, amount int) {
	user := &Data{
		DiscordID: userID,
		Score:     bigdata[userID].Score + amount,
		Name:      bigdata[userID].Name,
	}

	bigdata[userID] = *user

	WriteJSONFile()
}

func GiveStatsForUser(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userID := ""
	if i.Member != nil {
		userID = i.Member.User.ID
	}

	if i.User != nil {
		userID = i.User.ID
	}

	responseString := "<@" + bigdata[userID].DiscordID + ">'s score is " + strconv.Itoa(bigdata[userID].Score)

	embed := &discordgo.MessageEmbed{
		Title:       "Social Score",
		Description: responseString,
		Color:       808000,
	}

	responseBuilder(s, i, embed)
}

// man i love building beans
func responseBuilder(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println("Error sending response:", err)
	}
}

func CheckPointAdjustment(m *discordgo.Message) {
	message := m.Content

	switch true {
	case strings.Contains(message, "jackbox"):
		{
			updatePoints(m.Author.ID, -1)
		}
	case m.Author.ID == "168209080115134464":
		{
			updatePoints(m.Author.ID, -1)
		}
	}
}
