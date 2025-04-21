package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period.
	pingPeriod = 30 * time.Second
)

// pingConnection sends a ping control message to the given websocket connection.
func (m *Manager) pingConnection(conn *websocket.Conn) error {
	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write ping message directly
	if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to write ping message: %w", err)
	}

	return nil
}

// startHostPingLoop starts a goroutine that periodically sends ping messages to a connected host.
// It stops when the host's done channel is closed.
func (m *Manager) startHostPingLoop(host *connectedHost) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-host.done:
			return
		case <-ticker.C:
			if err := m.pingConnection(host.conn); err != nil {
				// Consider retrying the connection a few times before disconnecting the host immediately.

				logger.Logger.Sugar().Errorf("Failed to ping host %s: %v", host.hostInstance.ID, err)

				host.hostInstance.Status = core.HostInstanceStatusUnreachable

				tx, err := m.db.Beginx()
				if err != nil {
					logger.Logger.Sugar().Errorf("Failed to begin transaction: %v", err)
					continue
				}

				if err := m.db.UpdateHostInstance(context.Background(), tx, host.hostInstance); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host status: %v", err)
				}

				if err := tx.Rollback(); err != nil {
					logger.Logger.Sugar().Errorf("Failed to rollback transaction: %v", err)
				}

				m.DisconnectHost(host.hostInstance.ID) // Use manager method
				return
			}

			logger.Logger.Sugar().Debugf("Successfully pinged host %s", host.hostInstance.ID)
		}
	}
}

// startClientPingLoop starts a goroutine that periodically sends ping messages to a connected client.
// It stops when the client's done channel is closed.
func (m *Manager) startClientPingLoop(client *connectedClient) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-client.done:
			return
		case <-ticker.C:
			if err := m.pingConnection(client.conn); err != nil {
				// Consider retrying the connection a few times before disconnecting the client immediately.

				logger.Logger.Sugar().Errorf("Failed to ping client %s: %v", client.session.ID, err)

				tx, err := m.db.Beginx()
				if err != nil {
					logger.Logger.Sugar().Errorf("Failed to begin transaction: %v", err)
					continue
				}

				if err := m.db.DeleteSession(context.Background(), tx, client.session); err != nil {
					logger.Logger.Sugar().Errorf("Failed to delete session %s: %v", client.session.ID, err)
				}

				if err := tx.Rollback(); err != nil {
					logger.Logger.Sugar().Errorf("Failed to rollback transaction: %v", err)
				}

				m.DisconnectClient(client.session.ID) // Use manager method
				return
			}

			logger.Logger.Sugar().Debugf("Successfully pinged client %s", client.session.ID)
		}
	}
}
