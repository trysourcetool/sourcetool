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

	// Maximum time to wait for connection recovery before deleting the host instance.
	maxRecoveryWait = 6 * time.Hour
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
// It handles connection recovery attempts and deletes the host instance if recovery fails.
// It stops when the host's done channel is closed.
func (m *Manager) startHostPingLoop(host *connectedHost) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	var firstFailureTime *time.Time
	for {
		select {
		case <-host.done:
			logger.Logger.Sugar().Infof("Ping loop received done signal for host %s", host.hostInstance.ID)
			if firstFailureTime != nil && time.Since(*firstFailureTime) > maxRecoveryWait {
				logger.Logger.Sugar().Infof("Host %s was already unrecoverable (%v elapsed > %v) when done signal received, proceeding with deletion.", host.hostInstance.ID, time.Since(*firstFailureTime), maxRecoveryWait)
				if delErr := m.db.HostInstance().Delete(context.Background(), host.hostInstance); delErr != nil {
					logger.Logger.Sugar().Errorf("Failed to delete unrecoverable host instance %s during done signal handling: %v", host.hostInstance.ID, delErr)
				}
				m.DisconnectHost(host.hostInstance.ID)
			}
			return
		case <-ticker.C:
			if err := m.pingConnection(host.conn); err != nil {
				logger.Logger.Sugar().Errorf("Failed to ping host %s: %v", host.hostInstance.ID, err)

				if firstFailureTime == nil {
					now := time.Now()
					firstFailureTime = &now
					logger.Logger.Sugar().Infof("Recording first ping failure time for host %s: %v", host.hostInstance.ID, now)
				}

				if time.Since(*firstFailureTime) > maxRecoveryWait {
					logger.Logger.Sugar().Infof("Connection for host %s unrecoverable after %v, deleting instance.", host.hostInstance.ID, maxRecoveryWait)
					if delErr := m.db.HostInstance().Delete(context.Background(), host.hostInstance); delErr != nil {
						logger.Logger.Sugar().Errorf("Failed to delete unrecoverable host instance %s: %v", host.hostInstance.ID, delErr)
					}
					m.DisconnectHost(host.hostInstance.ID)
					return
				}

				if host.hostInstance.Status != core.HostInstanceStatusUnreachable {
					host.hostInstance.Status = core.HostInstanceStatusUnreachable // Update local state first
					if err := m.db.HostInstance().Update(context.Background(), host.hostInstance); err != nil {
						logger.Logger.Sugar().Errorf("Failed to update host %s status to unreachable: %v", host.hostInstance.ID, err)
					} else {
						logger.Logger.Sugar().Infof("Updated host %s status to unreachable due to ping failure.", host.hostInstance.ID)
					}
				}
				continue
			}

			logger.Logger.Sugar().Debugf("Successfully pinged host %s", host.hostInstance.ID)

			if firstFailureTime != nil {
				logger.Logger.Sugar().Infof("Connection recovered for host %s, resetting failure time", host.hostInstance.ID)
				firstFailureTime = nil
			}

			if host.hostInstance.Status != core.HostInstanceStatusOnline {
				host.hostInstance.Status = core.HostInstanceStatusOnline // Update local state first
				if err := m.db.HostInstance().Update(context.Background(), host.hostInstance); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host %s status to online after recovery: %v", host.hostInstance.ID, err)
				} else {
					logger.Logger.Sugar().Infof("Updated host %s status to online after successful ping.", host.hostInstance.ID)
				}
			}
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

				if err := m.db.Session().Delete(context.Background(), client.session); err != nil {
					logger.Logger.Sugar().Errorf("Failed to delete session %s: %v", client.session.ID, err)
				}

				m.DisconnectClient(client.session.ID) // Use manager method
				return
			}

			logger.Logger.Sugar().Debugf("Successfully pinged client %s", client.session.ID)
		}
	}
}
