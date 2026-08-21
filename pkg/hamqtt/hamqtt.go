// Package hamqtt is a small, device-agnostic client for exposing entities to
// Home Assistant via MQTT discovery (https://www.home-assistant.io/integrations/mqtt/).
// It knows nothing about spas; callers describe entities and this package
// handles topics, discovery payloads, availability, and command dispatch.
package hamqtt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Device describes the physical device entities are grouped under in Home
// Assistant's device registry. All entities registered on a Client share
// the same Device.
type Device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	SWVersion    string   `json:"sw_version,omitempty"`
}

// Entity describes a single Home Assistant MQTT-discovered entity. Every
// entity gets a state topic; a command topic is added, and Handler
// subscribed, only when Handler is non-nil.
type Entity struct {
	// Component is the Home Assistant MQTT platform, e.g. "sensor",
	// "binary_sensor", "switch", "select", "number", "text", "button".
	Component string
	// ObjectID must be unique within the device (e.g. "water_temperature").
	ObjectID string
	Name     string
	// Config holds additional discovery fields specific to Component, e.g.
	// unit_of_measurement, device_class, options, min, max, step.
	Config map[string]any
	// Handler is invoked with the raw command payload when a message
	// arrives on this entity's command topic. Leave nil for read-only
	// entities (sensor, binary_sensor).
	Handler func(payload string) error
}

// Config configures a Client's connection to the broker and how it
// identifies itself and its entities' owning device to Home Assistant.
type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string

	// NodeID prefixes every entity's state/command topics and forms part
	// of each entity's unique_id; it must be unique per physical device.
	NodeID string
	// DiscoveryPrefix is Home Assistant's MQTT discovery topic prefix.
	// Defaults to "homeassistant".
	DiscoveryPrefix string

	Device Device
}

const (
	defaultDiscoveryPrefix = "homeassistant"
	waitTimeout            = 10 * time.Second
)

// Client publishes discovery-compatible entities to an MQTT broker and
// dispatches incoming command messages to their registered handlers.
type Client struct {
	cfg Config
	mqc mqtt.Client
}

// New builds a Client and its underlying MQTT connection options, but does
// not connect - call Connect to do that.
func New(cfg Config) *Client {
	if cfg.DiscoveryPrefix == "" {
		cfg.DiscoveryPrefix = defaultDiscoveryPrefix
	}

	c := &Client{cfg: cfg}

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetOrderMatters(false).
		SetBinaryWill(c.availabilityTopic(), []byte("offline"), 1, true).
		SetOnConnectHandler(func(mqtt.Client) {
			c.SetAvailable(true)
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			slog.Warn("mqtt connection lost", "error", err)
		})

	c.mqc = mqtt.NewClient(opts)

	return c
}

func (c *Client) availabilityTopic() string {
	return c.cfg.NodeID + "/status"
}

func (c *Client) stateTopic(objectID string) string {
	return fmt.Sprintf("%s/%s/state", c.cfg.NodeID, objectID)
}

func (c *Client) commandTopic(objectID string) string {
	return fmt.Sprintf("%s/%s/set", c.cfg.NodeID, objectID)
}

func (c *Client) discoveryTopic(component, objectID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/config", c.cfg.DiscoveryPrefix, component, c.cfg.NodeID, objectID)
}

// Topic returns the MQTT topic for objectID/suffix under this client's node
// id, e.g. Topic("water_heater", "current_temperature"). It's a lower-level
// building block for entities that don't fit Register's single
// state-topic/command-topic shape - see PublishRaw and SubscribeRaw.
func (c *Client) Topic(objectID, suffix string) string {
	return fmt.Sprintf("%s/%s/%s", c.cfg.NodeID, objectID, suffix)
}

// Connect opens the connection to the broker and blocks until it succeeds,
// fails, or waitTimeout elapses.
func (c *Client) Connect() error {
	token := c.mqc.Connect()
	if !token.WaitTimeout(waitTimeout) {
		return fmt.Errorf("timed out connecting to mqtt broker %s", c.cfg.Broker)
	}

	return token.Error()
}

// Close publishes the offline availability message and disconnects.
func (c *Client) Close() {
	c.SetAvailable(false)

	// Give the offline message a moment to actually reach the broker
	// before we tear the connection down.
	c.mqc.Disconnect(250)
}

// SetAvailable publishes this device's availability. Register'd entities
// all share one availability topic, so this affects every entity at once -
// callers that can tell the underlying device has stopped responding
// (e.g. a bridge whose polls are failing) should call this to reflect that
// in Home Assistant, rather than leaving entities showing stale state as
// if nothing were wrong.
func (c *Client) SetAvailable(online bool) {
	payload := "offline"
	if online {
		payload = "online"
	}

	c.mqc.Publish(c.availabilityTopic(), 1, true, payload)
}

// Register publishes retained discovery config for e and, if it has a
// Handler, subscribes to its command topic.
func (c *Client) Register(e Entity) error {
	payload := map[string]any{
		"name":        e.Name,
		"state_topic": c.stateTopic(e.ObjectID),
	}

	if e.Handler != nil {
		payload["command_topic"] = c.commandTopic(e.ObjectID)
	}

	maps.Copy(payload, e.Config)

	if err := c.PublishDiscovery(e.Component, e.ObjectID, payload); err != nil {
		return err
	}

	if e.Handler == nil {
		return nil
	}

	return c.SubscribeRaw(c.commandTopic(e.ObjectID), e.Handler)
}

// PublishDiscovery publishes a fully custom discovery payload for
// component/objectID, filling in unique_id, availability_topic, and device
// unless payload already sets them. Unlike Register, it doesn't assume a
// single state/command topic pair - use it (with Topic, PublishRaw, and
// SubscribeRaw) for entities with more than one, e.g. water_heater's
// separate current/target temperature topics.
func (c *Client) PublishDiscovery(component, objectID string, payload map[string]any) error {
	full := map[string]any{
		"unique_id":          c.cfg.NodeID + "_" + objectID,
		"availability_topic": c.availabilityTopic(),
		"device":             c.cfg.Device,
	}

	maps.Copy(full, payload)

	blob, err := json.Marshal(full)
	if err != nil {
		return fmt.Errorf("marshal discovery payload for %s: %w", objectID, err)
	}

	if err := c.publishSync(c.discoveryTopic(component, objectID), 1, blob); err != nil {
		return fmt.Errorf("publish discovery for %s: %w", objectID, err)
	}

	return nil
}

// PublishState publishes value, retained, to objectID's state topic.
func (c *Client) PublishState(objectID, value string) {
	c.mqc.Publish(c.stateTopic(objectID), 0, true, value)
}

// PublishRaw publishes value, retained, to an arbitrary topic - see Topic.
func (c *Client) PublishRaw(topic, value string) {
	c.mqc.Publish(topic, 0, true, value)
}

// SubscribeRaw subscribes handler to an arbitrary command topic - see
// Topic. A handler error is logged, not returned to the caller, since it
// happens asynchronously long after Subscribe itself returns.
func (c *Client) SubscribeRaw(topic string, handler func(payload string) error) error {
	token := c.mqc.Subscribe(topic, 1, func(_ mqtt.Client, m mqtt.Message) {
		if err := handler(string(m.Payload())); err != nil {
			slog.Warn("mqtt command failed", "topic", topic, "payload", string(m.Payload()), "error", err)
		}
	})
	if !token.WaitTimeout(waitTimeout) {
		return fmt.Errorf("timed out subscribing to %s", topic)
	}

	return token.Error()
}

func (c *Client) publishSync(topic string, qos byte, payload any) error {
	token := c.mqc.Publish(topic, qos, true, payload)
	if !token.WaitTimeout(waitTimeout) {
		return fmt.Errorf("timed out publishing to %s", topic)
	}

	return token.Error()
}
