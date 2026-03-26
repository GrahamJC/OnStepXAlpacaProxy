package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"onstepx-alpaca-proxy/alpaca"
	"onstepx-alpaca-proxy/config"
	"onstepx-alpaca-proxy/logger"
	"onstepx-alpaca-proxy/onstepx"
	"onstepx-alpaca-proxy/server"
)

//go:embed webui/dist
var embeddedFS embed.FS

var (
	uiFS fs.FS
)

func main() {

	var err error

	// Setup UI file system
	uiFS, err = fs.Sub(embeddedFS, "webui/dist")
	if err != nil {
		// This is a critical error at startup. A message box is appropriate.
		fmt.Printf("FATAL: failed to load embedded UI files - %v\n", err)
		os.Exit(1)
	}

	// Initialise config
	if err := config.Load(); err != nil {
		fmt.Printf("FATAL: failed to load configuration - %v\n", err)
		os.Exit(1)
	}
	cfg := config.Get()

	// Create console logger and set as default
	//logger := slog.New(logger.NewHandler(&slog.HandlerOptions{Level: cfg.LogLevel}))
	logger := slog.New(logger.NewHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// If there is no COM port in the config attempt to find an OnStepX device
	comPort := cfg.ComPort
	if comPort == "" {
		comPort, err = onstepx.FindPort(cfg.BaudRate)
		//if err != nil {
		//	fmt.Printf("FATAL: error finding COM port - %v\n", err)
		//	os.Exit(1)
		//}
	}
	onstepx := onstepx.NewDevice(comPort, cfg.BaudRate)

	// Enable Alpaca discovery
	go alpaca.RespondToDiscovery()

	// Start HTTP server
	server.Start(uiFS, "0.0.0", onstepx)
}
