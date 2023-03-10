package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/NeverovDS/playlist/internal/api"
	"github.com/NeverovDS/playlist/internal/server"
	"github.com/NeverovDS/playlist/internal/service"
	"github.com/NeverovDS/playlist/logger"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	player := service.NewPlayer()

	go player.Playlist.RunPlayer()

	go func() {
		handlers := api.NewHandler(player)
		srv := server.NewServer(handlers.Init())
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	player.Playlist.LoadSongs()

	for {
		printMenu()
		scan.Scan()

		switch scan.Text() {
		case "1":
			go func() {
				player.Playlist.Play()
			}()
		case "2":
			player.Playlist.Pause()
		case "3":
			player.Playlist.AddSong(scan)
		case "4":
			player.Playlist.Next()
		case "5":
			player.Playlist.Prev()
		default:
			fmt.Println("Пожалуйста, выберите одно из значений предложенных в меню")
		}
	}

}

func printMenu() {
	fmt.Println("1. Play")
	fmt.Println("2. Pause")
	fmt.Println("3. Add Song")
	fmt.Println("4. Next")
	fmt.Println("5. Prev")
}
