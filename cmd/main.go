package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	discordyetanotherremoterunner "github.com/lajosbencz/discord-yet-another-remote-runner"
)

var configFile string = "config.yaml"
var config discordyetanotherremoterunner.Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %s\n", err)
	}
	if val := os.Getenv("YARR_CONFIG"); val != "" {
		configFile = val
	}
	if config, err = discordyetanotherremoterunner.ReadConfig(configFile); err != nil {
		log.Fatalf("error loading config file: %s\n", err)
	}
}

func main() {
	log.Println("starting")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(configFile)
	if err != nil {
		log.Fatalf("error adding file to watcher: %v", err)
	}

	token := os.Getenv("DISCORD_TOKEN")
	bot, err := discordyetanotherremoterunner.NewBot(token, config)
	if err != nil {
		log.Panicln(err)
	}
	bot.AddCommands(&discordyetanotherremoterunner.CommandServer{})
	if err := bot.Open(); err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := bot.Close(); err != nil {
			log.Printf("error while closing Discord connection: %s\n", err)
		}
	}()

	go func() {
		defer log.Println("watcher exited")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) && event.Name == configFile {
					config, err := discordyetanotherremoterunner.ReadConfig(configFile)
					if err != nil {
						log.Panicf("error reloading config: %v\n", err)
					} else {
						if config.Guild != "" && len(config.Servers) > 0 {
							log.Println("reloaded configuration")
							if err := bot.Close(); err != nil {
								log.Printf("failed to close bot: %s\n", err)
							}
							bot.SetConfig(config)
							if err := bot.Open(); err != nil {
								log.Panicln(err)
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watcher error: %v\n", err)
			}
		}
	}()

	// wait for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("shutting down")
}
