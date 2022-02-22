package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const Version = "v0.0.1-alpha"

var dg *discordgo.Session
var TOKEN string
var GUILDID string

var motivations = []string{
	"you can do it champ!",
	"I believe in you my child",
	"Tomorrow will be a better day!",
	"I love you, you are safe here",
	"Keep trying, I'm sure you've almost got it!",
	"Keep up the good work kid!",
}

var cmdStrings map[string]string = map[string]string{
	"testcmd": "Return String",
}

var msgEmbed = &discordgo.MessageEmbed{
	Author:      &discordgo.MessageEmbedAuthor{},
	Color:       0xfff200,
	Description: "> Serving catchup.json",
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "PHASE ONE",
			Value:  "Unnassuming USB Keys were sold to the masses of fans at Night of Fire 2 in LA and New York",
			Inline: false,
		},
		{
			Name:   "PHASE TWO",
			Value:  "usb keys turn out to have files on them, some of which are unreleased music, some are a bit cryptic",
			Inline: false,
		},
		{
			Name:   "PHASE THREE",
			Value:  "people discover fragmented rar",
			Inline: false,
		},
		{
			Name:   "PHASE FOUR",
			Value:  "people discover steganography",
			Inline: false,
		},
	},
	Image: &discordgo.MessageEmbedImage{
		URL: "https://cdn.discordapp.com/attachments/945144841908662332/945510431072612413/a2b2_green_logo_a2b2green2.png",
	},
	Thumbnail: &discordgo.MessageEmbedThumbnail{
		URL: "https://cdn.discordapp.com/attachments/945144841908662332/945510431072612413/a2b2_green_logo_a2b2green2.png",
	},
	Timestamp: time.Now().Format(time.RFC3339),
	Title:     "bot@a2b2.org:~/$ catchup",
}

func init() {
	// Print out a fancy logo!
	fmt.Printf(`arg-bot! %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	//Load Token from env (simulated with godotenv)
	TOKEN = os.Getenv("BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading token from env")
		os.Exit(1)
	}

	GUILDID = os.Getenv("GUILD_ID")
	if GUILDID == "" {
		log.Println("No GuildID specified in env")
		GUILDID = "" //this effectively specifies command registration as global
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "easymode",
			Description: "bot@a2b2.org:~/$ easymode",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "activate",
					Description: "TRUE/FALSE",
					Required:    true,
				},
			},
		},
		{
			Name:        "catchup",
			Description: "bot@a2b2.org:~/$ catchup",
		},
		{
			Name:        "motivation-simulator",
			Description: "bot@a2b2.org:~/$ motivate",
		},
		{
			Name:        "run",
			Description: "bot@a2b2.org:~/$ run <flags>",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "activate",
					Description: "TRUE/FALSE",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "flags",
					Description: "<cmdstring>",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"easymode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			activate := i.ApplicationCommandData().Options[0].BoolValue()

			if activate {
				content = "bot@a2b2.org:~/$ easymode activate: true"
			} else {
				content = "bot@a2b2.org:~/$ easymode activate: false"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"catchup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			var msgEmbed_array []*discordgo.MessageEmbed = []*discordgo.MessageEmbed{msgEmbed}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
					Embeds:  msgEmbed_array,
				},
			})

		},
		"motivation-simulator": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			content = motivate()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"run": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			mode := i.ApplicationCommandData().Options[0].BoolValue()
			command := i.ApplicationCommandData().Options[1].StringValue()
			if mode {
				content = checkCMD(command)
			} else {
				content = "bot@a2b2.org:~/$ run activate: false"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
	}
)

func init() {
	var err error
	dg, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating discordgo session!")
		os.Exit(1)
	}
}

func main() {
	var err error
	//Configure discordgo session bot
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//Register Bot Intents with Discord
	//worth noting MakeIntent is a no-op, but I want it there for doing something with pointers later
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GUILDID, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)

	defer dg.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")

	// Exit Normally.
	//exit
}

//Helper funcs

func comesFromDM(s *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		if channel, err = s.Channel(i.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func handleDmCheck(s *discordgo.Session, i *discordgo.InteractionCreate, dm bool) error {
	var err error
	if dm {
		log.Println("Message in dm")
		var content = "I can only be used in servers ;-("

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(content),
			},
		})
		return err
	}
	return nil
}

func motivate() string {
	rand.Seed(time.Now().UnixNano())
	l := len(motivations)
	i := rand.Intn(l)
	return motivations[i]

}

func checkCMD(cmd string) string {
	msg, ok := cmdStrings[cmd]
	if ok {
		return msg
	}

	return `error: {"type":"invalid_CMD"}`

}
