package socialscore

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Data struct {
	DiscordID string `json:"discordID"`
	Score     int    `json:"score"`
}

var bigdata = make([]Data, 0)

func ReadUserFile() {
	// Read the file contents
	fileData, err := os.ReadFile("socialscore\\db.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a slice of Person structs
	err = json.Unmarshal(fileData, &bigdata)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteJSONFile() error {

	// Marshal struct to JSON
	jsonData, err := json.MarshalIndent(bigdata, "", "    ")
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

func AddUser(s *discordgo.Session, m *discordgo.Message) {
	//check if the user already exists
	doesExists := alreadyExists(m.Author.ID)

	if doesExists {
		s.ChannelMessageSend(m.ChannelID, "You is already added")
		return
	}

	user := &Data{
		DiscordID: m.Author.ID,
		Score:     100,
	}

	bigdata = append(bigdata, *user)

	WriteJSONFile()
}

func alreadyExists(userID string) bool {
	for i := 0; i < len(bigdata); i++ {
		if bigdata[i].DiscordID == userID {
			return true
		}
	}

	return false
}
