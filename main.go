package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const prefix = "!magma"

func main() {
	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	ses, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	AuthorEmbed := discordgo.MessageEmbedAuthor{
		Name:    "AUTOMATICKÁ ZPRÁVA",
		IconURL: "https://cdn.discordapp.com/banners/1391686459903311953/9cecc6de5497b4f8de7647474ef6b7cd?size=512&quot",
	}

	ses.AddHandler(func(s *discordgo.Session, r *discordgo.MessageCreate) {
		if r.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(r.Content, " ")

		if args[0] != prefix {
			return
		}

		if len(args) >= 2 {
			switch args[1] {
			case "rules":
				embed := &discordgo.MessageEmbed{
					Title: "🌋 MagmaRealms - Rules",
					Description: "" +
						"- No racism, homophobia, toxic behavior or bullying\n" +
						"- Maintain decency and respect for others\n" +
						"- Do not violate the laws of the Czech Republic/Slovakia or the rules of the platform (Discord / Mojang)\n" +
						"- Respect admins and moderators - their word is final\n" +
						"- Advertising on other servers is prohibited (without permission)",
					Color:  0xFF0000,
					Author: &AuthorEmbed,
				}
				s.ChannelMessageSendEmbed(os.Getenv("CHANEL_RULES_ID"), embed)
				break
			case "info":
				embed := &discordgo.MessageEmbed{
					Title: "🌋 MagmaRealms - Info",
					Description: "" +
						"👑 Welcome to MagmaRealms! We are a Minecraft & Discord server focused on community, fun and a unique gaming experience. " +
						"Our team is constantly working on news, events and updates that will keep you playing!\n\n" +
						"🧱 What will you find here?\n\n" +
						"🛡️ BoxFight\n\n" +
						"🎁 Regular giveaways and competitions\n\n" +
						"🏆 Challenges and seasonal rewards\n\n" +
						"👥 Active and friendly community\n\n" +
						"💬 Discord for communication, support and chill\n\n" +
						"🌐 WEB: www.magmarealms.fun\n\n" +
						"🗺️ IP: mc.magmarealms.fun\n\n" +
						"⛏️ Version: 1.18.2 - 1.21.5",
					Color:  0xFF0000,
					Author: &AuthorEmbed,
				}
				s.ChannelMessageSendEmbed(os.Getenv("CHANEL_INFO_ID"), embed)
				break
			case "oznameni":
				message := strings.Join(args[2:], " ")
				message = strings.ReplaceAll(message, `\n`, "\n")

				if len(r.Attachments) > 0 {
					for _, attachment := range r.Attachments {
						if strings.HasPrefix(attachment.ContentType, "image/") {
							message += "\n" + attachment.URL
						}
					}
				}

				_, err := s.ChannelMessageSend(os.Getenv("CHANEL_OZNAMENI_ID"), "# 📣 Oznámení\n\n"+message)
				if err != nil {
					panic(err)
				}
				break
			case "oznameni-T":
				message := strings.Join(args[2:], " ")
				message = strings.ReplaceAll(message, `\n`, "\n")

				if len(r.Attachments) > 0 {
					for _, attachment := range r.Attachments {
						if strings.HasPrefix(attachment.ContentType, "image/") {
							message += "\n" + attachment.URL
						}
					}
				}

				_, err := s.ChannelMessageSend(r.ChannelID, "# 📣 Oznámení (T)\n\n"+message)
				if err != nil {
					panic(err)
				}
				break
			case "changelog":
				message := strings.Join(args[2:], " ")
				message = strings.ReplaceAll(message, `\g`, "\n- ")

				if len(r.Attachments) > 0 {
					for _, attachment := range r.Attachments {
						if strings.HasPrefix(attachment.ContentType, "image/") {
							message += "\n" + attachment.URL
						}
					}
				}

				s.ChannelMessageSend(os.Getenv("CHANEL_CHANGELOG_ID"), "# 🔨 Changelog\n\n"+message)
				break
			case "changelog-T":
				message := strings.Join(args[2:], " ")
				message = strings.ReplaceAll(message, `\g`, "\n- ")

				if len(r.Attachments) > 0 {
					for _, attachment := range r.Attachments {
						if strings.HasPrefix(attachment.ContentType, "image/") {
							message += "\n" + attachment.URL
						}
					}
				}

				s.ChannelMessageSend(r.ChannelID, "# 🔨 Changelog (T)\n\n"+message)
				break
			default:
				s.ChannelMessageSend(r.ChannelID, "příkazy: \n"+
					"```\n!magma rules (odešle pravidla do přednastaveného chanelu) [CHANEL_RULES_ID]"+
					"\n!magma info (odešle informace do přednastaveného chanelu) [CHANEL_INFO_ID]"+
					"\n!magma oznameni (odešle oznameni do přednastaveného chanelu) [CHANEL_OZNAMENI_ID]"+
					"\n!magma oznameni-T (odešle oznameni do lokálního chanelu) [/]"+
					"\n!magma changelog (odešle změny do přednastaveného chanelu) [CHANEL_CHANGELOG_ID]"+
					"\n!magma changelog-T (odešle změny do lokálního chanelu) [/]"+
					"\n!magma help (tento výpis) [/]```",
				)
			}
		} else {
			return
		}
	})

	ses.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = ses.Open()
	if err != nil {
		panic(err)
	}

	err = ses.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "idle",
		Activities: []*discordgo.Activity{
			{
				Name: os.Getenv("ACTIVITY"),
				Type: discordgo.ActivityTypeGame,
			},
		},
	})
	if err != nil {
		fmt.Println("Chyba při nastavení statusu:", err)
	}

	defer ses.Close()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
