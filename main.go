package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/johnnovosel/JawnBotGo/socialscore"
)

func main() {

	socialscore.ReadUserFile()

	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(handleCommand)

	// grab all the users of the server and add them to the db
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		// idk how to do it for multiple guilds rn lmao
		if len(r.Guilds) > 0 {
			guildID := r.Guilds[0].ID

			// Fetch all members of the guild
			members, err := dg.GuildMembers(guildID, "", 1000)
			if err != nil {
				fmt.Println("Error fetching guild members: ", err)
				return
			}

			socialscore.BulkAddUser(members)
		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
		socialscore.AddUser(r.User.ID, r.User.Username)
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	registerCommands(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// oh crap delete the commands so deficient people don't try to spam commands
	deleteCommands(dg)

	// Cleanly close down the Discord session.
	dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Create a new goroutine to handle the message.
	go handleMessage(s, m.Message)
}

func handleMessage(s *discordgo.Session, m *discordgo.Message) {

	// i hate pleb bot
	if m.Author.ID == "168209080115134464" {
		socialscore.CheckPointAdjustment(m)
	}

	//then we send the message to see if they gain or lose points
	socialscore.CheckPointAdjustment(m)
}

func handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "ss" {
		command := i.ApplicationCommandData().Options[0].StringValue()

		if command == "mystats" {
			socialscore.GiveStatsForUser(s, i)
		}
	}
}

func registerCommands(s *discordgo.Session) {

	cmd := &discordgo.ApplicationCommand{
		Name:        "ss",
		Description: "social score commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "what Social Score command do you want",
				Required:    true,
			},
		},
	}
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "931418680074596452", cmd)
	if err != nil {
		panic(err)
	}
}

func deleteCommands(s *discordgo.Session) {
	// get the list of commands
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		fmt.Println("Error getting commands:", err)
		return
	}

	for _, command := range commands {
		err = s.ApplicationCommandDelete(s.State.User.ID, "", command.ID)
		if err != nil {
			fmt.Printf("Error deleting command %s: %s\n", command.Name, err)
		}
	}
}
