package discordyetanotherremoterunner

import (
	"bytes"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandBasic struct {
	name        string
	description string
	command     ConfigCommand
}

func (c *CommandBasic) Name() string {
	return c.name
}

func (c *CommandBasic) Description() string {
	return c.description
}

func (c *CommandBasic) Build(bot *Bot) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
	}
}

func (c *CommandBasic) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := commandAcknowledge(s, i); err != nil {
		fmt.Println(err)
		return
	}

	out, err := execute(c.command)
	if err != nil {
		log.Println(err)
		commandResponse(s, i, fmt.Sprintf("error running %s command: %s", c.Name(), err), nil)
		return
	}
	commandResponse(s, i, "Info", &discordgo.File{
		Name:        "info.txt",
		ContentType: "text/plain",
		Reader:      bytes.NewBufferString(out),
	})
}

func NewCommandBasic(name, description string, command ConfigCommand) *CommandBasic {
	return &CommandBasic{
		name:        name,
		description: description,
		command:     command,
	}
}
