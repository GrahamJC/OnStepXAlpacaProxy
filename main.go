package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"onstepx-alpaca-proxy/config"
	"onstepx-alpaca-proxy/logger"
	"onstepx-alpaca-proxy/onstepx"
	"onstepx-alpaca-proxy/server"
)

//go:embed ui-vue/dist
var embeddedFS embed.FS

var (
	uiFS fs.FS
)

func main() {

	var err error

	// Setup UI file system
	uiFS, err = fs.Sub(embeddedFS, "ui-vue/dist")
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
	logger := slog.New(logger.NewHandler(&slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	// If there is no COM port in the config attempt to find an OnStepX device
	comPort := cfg.ComPort
	if comPort == "" {
		comPort, err = onstepx.FindPort(cfg.BaudRate)
		if err != nil {
			fmt.Printf("FATAL: error finding COM port - %v\n", err)
			os.Exit(1)
		}
	}
	onstepx := onstepx.NewDevice(comPort, cfg.BaudRate)

	// Connect OnStepX device
	if !onstepx.Connect() {
		slog.Error("failed to connect to OnStepX device")
	} else {

		// Set date/time
		err = onstepx.SetSiteTime(time.Now())

		dt, _ := onstepx.GetSiteTime()
		slog.Info(fmt.Sprintf("Date/Time: %s", dt))
		onstepx.SetSiteTime(dt)
		lat, _ := onstepx.GetSiteLatitude()
		slog.Info(fmt.Sprintf("Latitude: %f", lat))
		onstepx.SetSiteLatitude(lat)
		long, _ := onstepx.GetSiteLongitude()
		slog.Info(fmt.Sprintf("Longitude: %f", long))
		onstepx.SetSiteLongitude(long)
		elv, _ := onstepx.GetSiteElevation()
		slog.Info(fmt.Sprintf("Elevation: %f", elv))
		onstepx.SetSiteElevation(elv)
		utc, _ := onstepx.GetSiteUTCOffset()
		slog.Info(fmt.Sprintf("UTC Offset: %d", utc))
		onstepx.SetSiteUTCOffset(utc)
	}

	// Start HTTP server
	server.Start(onstepx, uiFS, "0.0.0")
}
