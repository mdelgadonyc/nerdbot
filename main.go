package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"nerdbot/delphi"
	"nerdbot/translate"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	input := goDotEnvVariable("TOKEN")
	discord, err := discordgo.New("Bot " + input)

	if err != nil {
		panic(nil)
	}

	// The "ready" event is a crucial point in the bot's lifecycle and is where you can start implementing your bot's functionality and interactions with Discord.
	discord.AddHandler(onReady)
	discord.AddHandler(onMessageCreate)

	// Open a WebSocket connection to Discord
	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}

	// Keep the bot running until it's stopped manually or encounters an error
	fmt.Printf("Listening for events\n")
	select {}
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	// Set the bot's status to "online"
	err := s.UpdateListeningStatus("Online")
	if err != nil {
		log.Printf("Error setting bot status: %v", err)
	}

	log.Println(event.User.Username + " is now online.")

}

func onMessageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	// Ensure message sender is not bot itself.
	if event.Author.ID == s.State.User.ID {
		return
	}

	message := event.Message.Content
	author := event.Author.Username

	log.Printf("Message from %s in channel %s: %s", author, event.ChannelID, message)

	// Check if the last character of the message is a question mark (?)
	if message[len(message)-1] == "?"[0] {
		//log.Println("We have encountered a question!")
		chatString := delphi.Delphi(message)
		// fmt.Println(chatString)
		s.ChannelMessageSend(event.ChannelID, chatString)
	}

	words := strings.Split(message, " ")

	if len(words) > 0 {
		firstWord := strings.ToLower(words[0])
		words = words[1:]
		message = strings.Join(words, " ")

		// fmt.Println("The first word is: ", firstWord)
		// fmt.Println("The rest of the sentence is: " + message)

		// if event.Content == "bye" {
		if firstWord == "bye" || firstWord == "bye!" {
			s.ChannelMessageSend(event.ChannelID, "Bye!")

			// Sleep for 5 seconds before closing connection
			time.Sleep(5 * time.Second)

			err := s.Close()
			if err != nil {
				log.Printf("Error closing session: %v", err)
			} else {
				log.Println(s.State.User.Username + " has been logged out and disconnected.")
			}
			os.Exit(0)
		}

		if firstWord == "hello" || firstWord == "hello!" || firstWord == "hi" || firstWord == "hi!" {
			s.ChannelMessageSend(event.ChannelID, "Hi!")

			// translate.Translate()
		}

		if firstWord == "translate" {

			translatedString := translate.Translate(message)

			// Convert the string to a rune slice
			runes := []rune(message)

			runes[0] = unicode.ToUpper(runes[0])
			message = string(runes)

			s.ChannelMessageSend(event.ChannelID, message+" in Mandarin is "+translatedString+".")
		}

		if firstWord == "thanks" || firstWord == "thanks!" {
			s.ChannelMessageSend(event.ChannelID, "No problem!")
		}
	} else {
		fmt.Println("The string is empty.")
		return
	}

}
