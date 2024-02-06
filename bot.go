package discordyetanotherremoterunner

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	config   Config
	discord  *discordgo.Session
	commands map[string]Command
}

func (b *Bot) AddCommands(c ...Command) {
	for _, cmd := range c {
		b.commands[cmd.Name()] = cmd
	}
}

func (b *Bot) SetConfig(c Config) {
	b.config = c

	for k, v := range c.Commands {
		b.commands[k] = NewCommandBasic(k, k, v)
	}
}

func (b *Bot) Config() Config {
	return b.config
}

func (b *Bot) Open() error {
	return b.discord.Open()
}

func (b *Bot) Close() error {
	if err := b.discord.Close(); err != nil {
		return err
	}
	log.Println("discord closed")
	return nil
}

func NewBot(t string, c Config) (*Bot, error) {
	d, err := discordgo.New("Bot " + t)
	if err != nil {
		return nil, err
	}
	b := &Bot{
		discord:  d,
		commands: map[string]Command{},
	}

	b.SetConfig(c)

	// login handler
	b.discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as %v#%v for guild %s\n", s.State.User.Username, s.State.User.Discriminator, b.config.Guild)

		// command handlers
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if cmd, ok := b.commands[i.ApplicationCommandData().Name]; ok {
				go cmd.Handle(s, i) // @todo: rate limit
			} else {
				commandResponse(s, i, fmt.Sprintf("no such command: %s", cmd.Name()), nil)
			}
		})

		// command definitions
		for _, cmd := range b.commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, b.config.Guild, cmd.Build(b))
			if err != nil {
				log.Printf("cannot create '%v' command: %v\n", cmd.Name(), err)
			}
		}
	})

	return b, nil
}
