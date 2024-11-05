package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
)

func main() {
	// logHandler, err := ycloggingslog.New(ycloggingslog.Options{
	// 	LogGroupId:   os.Getenv("LOGGING_GROUP_ID"),
	// 	ResourceType: "app",
	// 	ResourceId:   "season-3",
	// 	Credentials:  ycsdk.OAuthToken(os.Getenv("LOGGING_TOKEN")),
	// 	Level:        slog.LevelDebug,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set logger globally for convenience
	// slog.SetDefault(slog.New(logHandler))

	client, err := api.New(os.Getenv("SERVER_URL"), os.Getenv("SERVER_TOKEN"))
	if err != nil {
		slog.With(slog.Any("error", err)).Error("fail init API client")
		os.Exit(1)
	}

	game := game.New(client)
	_ = game

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// TODO

	<-ctx.Done()
	slog.Info("got stop signal...")
	stopNotify()
}
