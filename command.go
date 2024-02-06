package discordyetanotherremoterunner

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Name() string
	Description() string
	Build(*Bot) *discordgo.ApplicationCommand
	Handle(*discordgo.Session, *discordgo.InteractionCreate)
}

func commandAcknowledge(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "..."},
	})
	if err != nil {
		return fmt.Errorf("failed to acknowledge command: %s", err)
	}
	return nil
}

func commandResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string, file *discordgo.File) {
	data := &discordgo.WebhookEdit{
		Content: &content,
	}
	if file != nil {
		data.Files = []*discordgo.File{file}
	}
	_, err := s.InteractionResponseEdit(i.Interaction, data)
	if err != nil {
		log.Printf("error while responding: %s\n", err)
	}
}

func commandOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
