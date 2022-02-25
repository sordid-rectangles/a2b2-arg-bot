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

const Version = "v0.0.2-alpha"

var dg *discordgo.Session
var TOKEN string
var GUILDID string

var timerExp = time.Date(2022, 3, 1, 17, 0, 0, 0, time.UTC)

var link = "https://www.youtube.com/watch?v=kpnwMgV9HeI"

var motivations = []string{
	"`You can do it champ!`",
	"`I believe in you my child`",
	"`Tomorrow will be a better day!`",
	"`I love you, you are safe here`",
	"`Keep trying, I'm sure you've almost got it!`",
	"`Keep up the good work kid!`",
}

var cmdStrings map[string]string = map[string]string{
	"observer": "confirm",
}

var msgEmbed = &discordgo.MessageEmbed{
	Author:      &discordgo.MessageEmbedAuthor{},
	Color:       0xfff200,
	Description: "> Serving catchup.json",
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "PHASE ONE",
			Value:  "- Unnassuming a2b2 branded USB Keys were sold to the masses of fans at Night of Fire 2 in LA and New York. \n- The usb keys turned out to have files on them, some of which are unreleased music, some were a bit cryptic.\n- People discovered a fragmented rar, and tracked down and assembled the pieces to decompress it.\n- The rar contained a congratulatory video, and a link to a mysterious website.",
			Inline: false,
		},
		{
			Name:   "PHASE TWO",
			Value:  "- The website was periodically updated with codes hidden in varied media formats, including various forms of steganography, and SSTV encoding.\n- The end goal was to find a set of key-codes to open up \"the portal\"\n- Some of the content hidden on the website, lead to an interactive livestream event in which further codes were hidden\n- At the end of phase 2, the site was wiped and a videogame was uploaded that contained fragments of the final portal key. Participants failed miserably to decipher this masterpiece, and as such, phase 3 CODENAME: EASYMODE was developed.",
			Inline: false,
		},
		{
			Name:   "PHASE THREE",
			Value:  "Welcome to EASYMODE :)",
			Inline: false,
		},
		{
			Name:   "END?",
			Value:  "Easymode contained a link to a cryptic reading of a john dee work, with a codephrase midaway. The codephrase when entered into the bot returned an image file, which had a metadata tag with a link to the 2 deaths grips remixes Andy dropped at NOF2. Angel Room remains unsolved, and any intrepid souls who wish to continue may do so at their own initiative : ).\nThanks for playing along?",
			Inline: false,
		},
	},
	Thumbnail: &discordgo.MessageEmbedThumbnail{
		URL:    "https://cdn.discordapp.com/attachments/920884723138584599/946153967258312764/a2b2.1.gif",
		Width:  200,
		Height: 200,
	},
	Timestamp: time.Now().Format(time.RFC3339),
	Title:     "bot@a2b2.org:~/$ catchup",
}

var videoEmbed = &discordgo.MessageEmbed{
	Author:      &discordgo.MessageEmbedAuthor{},
	Color:       0xfff200,
	Description: "> Serving euclid.exe",
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "bot@a2b2.org:~/$ easymode",
			Value:  "https://youtu.be/7VbMtMsUQzw",
			Inline: false,
		},
	},
	Video: &discordgo.MessageEmbedVideo{
		URL:    "https://youtu.be/7VbMtMsUQzw",
		Width:  300,
		Height: 300,
	},
	Image: &discordgo.MessageEmbedImage{
		URL:    "https://cdn.discordapp.com/attachments/920884723138584599/946153967258312764/a2b2.1.gif",
		Width:  200,
		Height: 200,
	},
	Timestamp: time.Now().Format(time.RFC3339),
	Title:     "bot@a2b2.org:~/$ easymode true",
}

var runEmbed = &discordgo.MessageEmbed{
	Author:      &discordgo.MessageEmbedAuthor{},
	Color:       0xfff200,
	Description: "`> Serving geometry.json`",
	Image: &discordgo.MessageEmbedImage{
		URL:    "https://cdn.discordapp.com/attachments/920884723138584599/946863634229911572/regarding_geometry.png",
		Width:  300,
		Height: 300,
	},
	Timestamp: time.Now().Format(time.RFC3339),
	Title:     "bot@a2b2.org:~/$ run true observer",
}

var timerEmbed = &discordgo.MessageEmbed{
	Author:      &discordgo.MessageEmbedAuthor{},
	Color:       0xfff200,
	Description: "`> Serving timer.json`",
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "`Timer`",
			Value:  fmt.Sprintf("`Timer expires: %s`", timerExp),
			Inline: false,
		},
	},
	Thumbnail: &discordgo.MessageEmbedThumbnail{
		URL:    "https://cdn.discordapp.com/attachments/920884723138584599/946153967258312764/a2b2.1.gif",
		Width:  200,
		Height: 200,
	},
	Timestamp: time.Now().Format(time.RFC3339),
	Title:     "bot@a2b2.org:~/$ timer",
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
		{
			Name:        "timer",
			Description: "bot@a2b2.org:~/$ timer",
		},
		{
			Name:        "nightcore",
			Description: "bot@a2b2.org:~/$ nightcore",
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

			var vidEmbed_array []*discordgo.MessageEmbed
			if activate {
				content = "`bot@a2b2.org:~/$ easymode activate: true`"
				vidEmbed_array = []*discordgo.MessageEmbed{videoEmbed}
			} else {
				content = "`bot@a2b2.org:~/$ easymode activate: false`"
				// vidEmbed_array = []*discordgo.MessageEmbed{videoEmbed}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
					Embeds:  vidEmbed_array,
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

			var msgEmbed_array []*discordgo.MessageEmbed
			if mode {
				if checkCMD(command) == "confirm" {
					msgEmbed_array = []*discordgo.MessageEmbed{runEmbed}
				} else {
					content = "`bot@a2b2.org:~/$ run activate: fail`"
				}

			} else {
				content = "`bot@a2b2.org:~/$ run activate: false`"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
					Embeds:  msgEmbed_array,
				},
			})
		},
		"timer": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			var msgEmbed_array []*discordgo.MessageEmbed = []*discordgo.MessageEmbed{timerEmbed}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
					Embeds:  msgEmbed_array,
				},
			})

		},
		"nightcore": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//var content string

			check, err := comesFromDM(s, i)
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
				return
			}
			err = handleDmCheck(s, i, check)

			//content = fmt.Sprintf("@everyone %s", link)

			//_, err = s.ChannelMessageSend(i.ChannelID, content)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("NIGHTCORE! %s", link),
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

	return "`error: {type:invalid_CMD}`"

}
