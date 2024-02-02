package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	dg    *discordgo.Session
	queue []string
)

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	token := "MTE4NzgxMzA1Njc4MTIzNDMxNw.GnoBFF.Zj8ZyUdXFAGOUTd9VR9pId8smC55xL8zgrz3nI"

	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar a sessão do Discord:", err)
		return
	}

	dg.AddHandler(commandHandler)

	err = dg.Open()
	if err != nil {
		log.Fatal("Erro ao abrir a conexão com o Discord:", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot está funcionando. Pressione CTRL + C para sair.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	switch data.Name {
	case "play":

		voiceChannelID, err := findUserVoiceChannel(s, i.GuildID, i.Interaction.Member.User.ID)
		if err != nil {
			fmt.Println("Erro ao encontrar o canal de voz do usuário:", err)
			return
		}
		if voiceChannelID != "" {
			if !isBotInVoiceChannel(s, i.GuildID, voiceChannelID) {
				_, err := joinVoiceChannel(s, i.GuildID, voiceChannelID)
				if err != nil {
					fmt.Println("Erro ao entrar no canal de voz:", err)
					return
				}

				fmt.Println("Bot entrou no canal de voz do autor da mensagem")
			}
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "https://open.spotify.com/playlist/5z2dTZUjDD90wM4Z9youwS?si=407acac85d6f4b08",
			},
		})
	case "parar":
		// Lógica para o comando /parar
	case "lista":
		// Lógica para o comando /lista
	default:
		// Mensagem para o caso de um comando desconhecido
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID, channelID string) (*discordgo.VoiceConnection, error) {
	voiceConnection, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}
	return voiceConnection, nil
}

func findUserVoiceChannel(s *discordgo.Session, guildID, userID string) (string, error) {

	guild, err := s.State.Guild(guildID)
	if err != nil {
		return "", err
	}

	for _, vs := range guild.VoiceStates {
		fmt.Print("VoiceStateUser: " + vs.UserID)
		fmt.Print("ParamUser: " + userID)
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", fmt.Errorf("Usuário não encontrado no canal de voz")
}

func isBotInVoiceChannel(s *discordgo.Session, guildID, channelID string) bool {
	voiceConnections := s.VoiceConnections

	for _, voiceConnection := range voiceConnections {
		if voiceConnection.GuildID == guildID && voiceConnection.ChannelID == channelID {
			return true
		}
	}

	return false
}
