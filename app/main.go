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

func main() {
	token := "TOKEN_DO_SEU_BOT_DISCORD"

	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar a sessão do Discord:", err)
		return
	}

	dg.AddMessageCreateHandler(messageCreate)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "/play" {
		s.ChannelMessageSend(m.ChannelID, "Tocando música!")
		// Aqui você conectaria o bot a um reprodutor de música, mas para esta etapa, vamos apenas enviar a mensagem.
	} else if m.Content == "/stop" {
		s.ChannelMessageSend(m.ChannelID, "Música parada!")
		// Aqui você pararia a música que está sendo reproduzida, mas para esta etapa, enviamos apenas a mensagem.
	} else if m.Content == "/list" {
		if len(queue) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Não há músicas na lista.")
			return
		}

		message := "Lista de músicas na fila:\n"
		for i, song := range queue {
			message += fmt.Sprintf("%d. %s\n", i+1, song)
		}
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
