package main

import (
	"embed"
	"fmt"
	"io/fs"
	"onstepx-alpaca-proxy/connection"
	"time"
	//"onstepx-alpaca-proxy/alpaca"
	//"onstepx-alpaca-proxy/config"
	//"onstepx-alpaca-proxy/logger"
	//"onstepx-alpaca-proxy/onstepx"
	//"onstepx-alpaca-proxy/server"
)

//go:embed webui/dist
var embeddedFS embed.FS

var (
	uiFS fs.FS
)

func main() {

	osxSerial := connection.NewSerialConnection("COM4", 57600, 1000)
	osx := connection.NewOnstepX(osxSerial)
	if err := osx.Connect(); err != nil {
		fmt.Printf("Error connecting to OnStepX: %v\n", err)
	} else {
		if v, err := osx.GetVersionFull(); err == nil {
			fmt.Printf("Version : %v\n", v)
		}
		if v, err := osx.GetVersionNumber(); err == nil {
			fmt.Printf("VersionNumber : %v\n", v)
		}
		if v, err := osx.GetVersionFull(); err == nil {
			fmt.Printf("VersionFull : %v\n", v)
		}
		if l, err := osx.GetSiteLatitude(); err == nil {
			fmt.Printf("SiteLatitude : %v\n", l)
		} else {
			fmt.Printf("Error getting site latitude: %v\n", err)
		}
		if err := osx.SetSiteLatitude(52.7); err != nil {
			fmt.Printf("Error setting site latitude: %v\n", err)
		}
		if l, err := osx.GetSiteLatitude(); err == nil {
			fmt.Printf("SiteLatitude : %v\n", l)
		} else {
			fmt.Printf("Error getting site latitude: %v\n", err)
		}
		if st, err := osx.GetSiderealTime(); err == nil {
			fmt.Printf("SiderealTime : %v\n", st)
		} else {
			fmt.Printf("Error getting sidereal time: %v\n", err)
		}
		if err := osx.SetLocalDateTime(time.Date(2026, 6, 12, 15, 23, 45, 0, time.UTC)); err != nil {
			fmt.Printf("Error setting local date/time: %v\n", err)
		}
		if ldt, err := osx.GetLocalDateTime(); err == nil {
			fmt.Printf("LocalTime : %v\n", ldt)
		} else {
			fmt.Printf("Error getting local date/time: %v\n", err)
		}
		if err := osx.SetSiteUTCOffset(-5); err != nil {
			fmt.Printf("Error setting UTC offset: %v\n", err)
		}
		if ldt, err := osx.GetLocalDateTime(); err == nil {
			fmt.Printf("LocalTime : %v\n", ldt)
		} else {
			fmt.Printf("Error getting local date/time: %v\n", err)
		}
		osx.Disconnect()
	}
}

/*
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
*/
