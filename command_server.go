package discordyetanotherremoterunner

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandServer struct {
	bot *Bot
}

func (c *CommandServer) Name() string {
	return "server"
}

func (c *CommandServer) Description() string {
	return "Server commands"
}

func (c *CommandServer) Build(bot *Bot) *discordgo.ApplicationCommand {
	c.bot = bot
	servers := []*discordgo.ApplicationCommandOptionChoice{}
	for k, v := range bot.config.Servers {
		servers = append(servers, &discordgo.ApplicationCommandOptionChoice{
			Name:  v.Name,
			Value: k,
		})
	}
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Select action",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Start",
						Value: "start",
					},
					{
						Name:  "Stop",
						Value: "stop",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "server",
				Description: "Select server",
				Required:    true,
				Choices:     servers,
			},
		},
	}
}

func (c *CommandServer) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := commandAcknowledge(s, i); err != nil {
		fmt.Println(err)
		return
	}

	options := commandOptionMap(i.ApplicationCommandData().Options)
	var action, server string
	if opt, ok := options["action"]; ok {
		action = opt.StringValue()
	}
	if opt, ok := options["server"]; ok {
		server = opt.StringValue()
	}
	cfgServer, ok := c.bot.config.Servers[server]
	if !ok {
		commandResponse(s, i, "no such server", nil)
		return
	}
	msg := fmt.Sprintf("%s server: %s", action, server)
	log.Println(msg)
	var err error = nil
	switch action {
	case "start":
		_, err = execute(cfgServer.Start)
	case "stop":
		_, err = execute(cfgServer.Stop)
	default:
		log.Printf("unknown server action: %s\n", action)
		commandResponse(s, i, fmt.Sprintf("unknown action: %s", action), nil)
		return
	}
	if err != nil {
		commandResponse(s, i, err.Error(), nil)
		return
	}
	commandResponse(s, i, ":white_check_mark: "+msg, nil)
}
