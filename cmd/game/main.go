package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/characters"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

func main() {
	setup()

	client, err := api.New(os.Getenv("SERVER_URL"), os.Getenv("SERVER_TOKEN"))
	if err != nil {
		slog.With(slog.Any("error", err)).Error("fail init API client")
		os.Exit(1)
	}

	game := game.New(client)
	ishtar := characters.NewIshtar(client, game)

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go ishtar.Live(ctx)

	<-ctx.Done()
	slog.Info("got stop signal...")
	stopNotify()
}

func setup() {
	logHandler, err := ycloggingslog.New(ycloggingslog.Options{
		LogGroupId:   os.Getenv("LOGGING_GROUP_ID"),
		ResourceType: "app",
		ResourceId:   "season-3",
		Credentials:  ycsdk.OAuthToken(os.Getenv("LOGGING_TOKEN")),
		Level:        slog.LevelDebug,
	})
	if err != nil {
		panic(err)
	}

	// set logger globally for convenience
	slog.SetDefault(slog.New(logHandler))

	monitoringClient := ycmonitoringgo.NewClient(
		os.Getenv("MONITORING_FOLDER"),
		os.Getenv("MONITORING_TOKEN"),
		ycmonitoringgo.WithLogger(slog.Default().With(ycloggingslog.Stream, "monitoring")),
	)

	go monitoringClient.Run(context.Background(), ycmonitoringgo.DefaultRegistry, 30*time.Second)
}
