// Package mqttbridge exposes a spa, over an existing safespa.SafeSpa, as a
// set of Home Assistant MQTT-discovered entities: sensors and selects for
// everything the spa reports, and command topics wired to the same Spanet
// setters the HTTP API uses.
package mqttbridge

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/freman/spanet/pkg/hamqtt"
	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

// Config configures a Bridge.
type Config struct {
	Broker   string
	Username string
	Password string
	// NodeID uniquely identifies this spa to Home Assistant; used as both
	// the MQTT client ID and the topic/unique_id prefix for every entity.
	NodeID string
	// PollInterval is how often the bridge re-reads spa status and
	// republishes entity state. Defaults to 15s.
	PollInterval time.Duration
}

// unavailableAfter is how many consecutive failed polls it takes before
// the bridge marks the device unavailable in Home Assistant, rather than
// leaving every entity showing its last known (now stale) state.
const unavailableAfter = 3

// Bridge polls a spa on an interval and republishes its state to Home
// Assistant over MQTT, while accepting commands back over the same broker.
type Bridge struct {
	cfg      Config
	spa      *safespa.SafeSpa
	mqc      *hamqtt.Client
	entities []entity

	consecutiveFailures int
}

// initialStatusRetry is how long to wait between attempts to read the
// spa's initial status before Start can finish setting up. The spa being
// briefly unreachable at boot is a real, observed condition (its wifi
// bridge only tolerates one connection and can get stuck holding a stale
// one) - Start retries through that rather than giving up.
const initialStatusRetry = 5 * time.Second

// NewBridge builds a Bridge ready to Start. It does no I/O - unlike
// talking to the spa or the broker, this can't hang or fail, so callers
// don't need to guard the rest of their startup against it.
func NewBridge(cfg Config, spa *safespa.SafeSpa) *Bridge {
	if cfg.NodeID == "" {
		cfg.NodeID = "spanet"
	}

	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 15 * time.Second
	}

	mqc := hamqtt.New(hamqtt.Config{
		Broker:   cfg.Broker,
		ClientID: cfg.NodeID,
		Username: cfg.Username,
		Password: cfg.Password,
		NodeID:   cfg.NodeID,
		Device: hamqtt.Device{
			Identifiers:  []string{cfg.NodeID},
			Name:         "Spa",
			Manufacturer: "Spanet",
		},
	})

	return &Bridge{
		cfg: cfg,
		spa: spa,
		mqc: mqc,
	}
}

// Start reads the spa's initial status (retrying until it succeeds or ctx
// is done, so a spa that's unreachable at boot doesn't block anything
// else in the process from starting), connects to the broker, publishes
// discovery config and initial state for every entity, subscribes command
// topics, and then polls spa status on cfg.PollInterval until ctx is done,
// at which point it publishes the offline availability message and
// disconnects.
func (b *Bridge) Start(ctx context.Context) error {
	status, err := b.awaitInitialStatus(ctx)
	if err != nil {
		return err
	}

	b.entities = entities(b.spa, status.Pumps)

	if err := b.mqc.Connect(); err != nil {
		return fmt.Errorf("connecting to mqtt broker: %w", err)
	}

	for _, e := range b.entities {
		if err := b.mqc.Register(e.Entity); err != nil {
			b.mqc.Close()

			return fmt.Errorf("registering %s: %w", e.ObjectID, err)
		}
	}

	if err := registerWaterHeater(b.mqc, b.spa); err != nil {
		b.mqc.Close()

		return fmt.Errorf("registering water heater: %w", err)
	}

	b.poll()

	ticker := time.NewTicker(b.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			b.mqc.Close()

			return nil
		case <-ticker.C:
			b.poll()
		}
	}
}

// awaitInitialStatus retries reading spa status until it succeeds or ctx
// is done. The first attempt is silent; only repeated failures are logged,
// since a single miss during a normal reconnect is not noteworthy.
func (b *Bridge) awaitInitialStatus(ctx context.Context) (spanet.Status, error) {
	attempt := 0

	for {
		var status spanet.Status

		err := b.spa.Do(func(s *spanet.Spanet) error {
			var err error
			status, err = s.GetStatus()

			return err
		})
		if err == nil {
			return status, nil
		}

		attempt++
		if attempt > 1 {
			slog.Warn("mqtt bridge: waiting for spa to become reachable", "error", err, "attempt", attempt)
		}

		select {
		case <-ctx.Done():
			return spanet.Status{}, ctx.Err()
		case <-time.After(initialStatusRetry):
		}
	}
}

func (b *Bridge) poll() {
	var status spanet.Status

	err := b.spa.Do(func(s *spanet.Spanet) error {
		var err error
		status, err = s.GetStatus()

		return err
	})
	if err != nil {
		slog.Warn("mqtt bridge: failed to read spa status", "error", err)

		b.consecutiveFailures++
		if b.consecutiveFailures == unavailableAfter {
			b.mqc.SetAvailable(false)
		}

		return
	}

	if b.consecutiveFailures >= unavailableAfter {
		b.mqc.SetAvailable(true)
	}

	b.consecutiveFailures = 0

	for _, e := range b.entities {
		if e.State == nil {
			continue
		}

		value, ok := e.State(status)
		if !ok {
			continue
		}

		b.mqc.PublishState(e.ObjectID, value)
	}

	publishWaterHeaterState(b.mqc, status)
}
