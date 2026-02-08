package alpaca

import (
	"fmt"
	"log/slog"
	"net"

	"onstepx-alpaca-proxy/config"
)

// RespondToDiscovery listens for Alpaca discovery packets on UDP port 32227
// and responds with the server's listening port.
func RespondToDiscovery() {
	listenAddr := config.Get().ListenAddress
	udpAddress := fmt.Sprintf("%s:32227", listenAddr)

	addr, err := net.ResolveUDPAddr("udp4", udpAddress)
	if err != nil {
		slog.Error("Discovery: Could not resolve UDP address '%s': %v", "udpAddress", udpAddress, "error", err)
		return
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		slog.Error("Discovery: Could not listen on UDP address", "udpAddress", udpAddress, "error", err)
		slog.Info("HINT: This may be caused by another Alpaca application running, or a permissions issue.")
		return
	}
	defer conn.Close()
	slog.Info("Alpaca discovery responder started", "udpAddress", udpAddress)

	discoveryMsg := []byte("alpacadiscovery1")
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			slog.Warn("Discovery: Error reading from UDP", "error", err)
			continue
		}

		if string(buffer[:n]) == string(discoveryMsg) {
			slog.Debug("Discovery: Request received", "remoteAddr", remoteAddr)

			// Get the current network port from the config
			port := config.Get().NetworkPort
			response := fmt.Sprintf(`{"AlpacaPort": %d}`, port)

			_, err := conn.WriteToUDP([]byte(response), remoteAddr)
			if err != nil {
				slog.Error("Discovery: Failed to send response", "remoteAddr", remoteAddr, "error", err)
			} else {
				slog.Debug("Discovery: Sent response", "response", response, "remoteAddr", remoteAddr)
			}
		}
	}
}
