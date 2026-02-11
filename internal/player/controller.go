package player

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

const socketPath = "/tmp/ymusic-mpv.sock"

type State struct {
	Playing   bool
	Position  float64
	Duration  float64
	Volume    float64
	Idle      bool
	TrackURL  string
}

type Event struct {
	Type  string // "property-change", "end-file", "idle"
	Name  string
	Value interface{}
}

type Controller struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	conn     net.Conn
	reqID    atomic.Int64
	state    State
	Events   chan Event
	started  bool
}

func NewController(volume float64) *Controller {
	if volume <= 0 {
		volume = 70
	}
	return &Controller{
		Events: make(chan Event, 64),
		state:  State{Volume: volume},
	}
}

func (c *Controller) Start() error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	os.Remove(socketPath)

	cmd := exec.Command("mpv",
		"--idle",
		"--no-video",
		"--no-terminal",
		"--input-ipc-server="+socketPath,
	)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start mpv: %w", err)
	}

	// Wait for socket (without holding mutex)
	var conn net.Conn
	for i := 0; i < 50; i++ {
		nc, err := net.Dial("unix", socketPath)
		if err == nil {
			conn = nc
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if conn == nil {
		cmd.Process.Kill()
		return fmt.Errorf("mpv socket not ready")
	}

	c.mu.Lock()
	c.cmd = cmd
	c.conn = conn
	c.started = true

	// Observe properties (mutex held, use locked version)
	c.observe("time-pos", 1)
	c.observe("duration", 2)
	c.observe("pause", 3)
	c.observe("volume", 4)
	c.observe("idle-active", 5)

	// Set initial volume
	c.sendCommandLocked(map[string]interface{}{
		"command": []interface{}{"set_property", "volume", c.state.Volume},
	})
	c.mu.Unlock()

	go c.readLoop()
	return nil
}

func (c *Controller) observe(prop string, id int) {
	c.sendCommandLocked(map[string]interface{}{
		"command":    []interface{}{"observe_property", id, prop},
		"request_id": id + 1000,
	})
}

// sendCommandLocked writes to mpv socket without acquiring mutex.
// Caller must hold c.mu.
func (c *Controller) sendCommandLocked(cmd map[string]interface{}) (json.RawMessage, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	data, _ := json.Marshal(cmd)
	data = append(data, '\n')
	_, err := c.conn.Write(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *Controller) sendCommand(cmd map[string]interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sendCommandLocked(cmd)
}

func (c *Controller) readLoop() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		var msg map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue
		}

		eventType, _ := msg["event"].(string)
		if eventType == "property-change" {
			name, _ := msg["name"].(string)
			value := msg["data"]
			c.updateState(name, value)
			select {
			case c.Events <- Event{Type: "property-change", Name: name, Value: value}:
			default:
			}
		} else if eventType == "end-file" {
			reason, _ := msg["reason"].(string)
			select {
			case c.Events <- Event{Type: "end-file", Name: reason}:
			default:
			}
		} else if eventType == "idle" {
			select {
			case c.Events <- Event{Type: "idle"}:
			default:
			}
		}
	}
}

func (c *Controller) updateState(name string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch name {
	case "time-pos":
		if v, ok := value.(float64); ok {
			c.state.Position = v
		}
	case "duration":
		if v, ok := value.(float64); ok {
			c.state.Duration = v
		}
	case "pause":
		if v, ok := value.(bool); ok {
			c.state.Playing = !v
		}
	case "volume":
		if v, ok := value.(float64); ok {
			c.state.Volume = v
		}
	case "idle-active":
		if v, ok := value.(bool); ok {
			c.state.Idle = v
		}
	}
}

func (c *Controller) GetState() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

func (c *Controller) LoadURL(url string) error {
	c.mu.Lock()
	c.state.TrackURL = url
	_, err := c.sendCommandLocked(map[string]interface{}{
		"command": []interface{}{"loadfile", url},
	})
	c.mu.Unlock()
	return err
}

func (c *Controller) TogglePause() error {
	_, err := c.sendCommand(map[string]interface{}{
		"command": []interface{}{"cycle", "pause"},
	})
	return err
}

func (c *Controller) Seek(seconds float64) error {
	_, err := c.sendCommand(map[string]interface{}{
		"command": []interface{}{"seek", seconds, "relative"},
	})
	return err
}

func (c *Controller) SeekAbsolute(seconds float64) error {
	_, err := c.sendCommand(map[string]interface{}{
		"command": []interface{}{"seek", seconds, "absolute"},
	})
	return err
}

func (c *Controller) SetVolume(vol float64) error {
	_, err := c.sendCommand(map[string]interface{}{
		"command": []interface{}{"set_property", "volume", vol},
	})
	return err
}

func (c *Controller) Stop() error {
	_, err := c.sendCommand(map[string]interface{}{
		"command": []interface{}{"stop"},
	})
	return err
}

func (c *Controller) Quit() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		c.cmd.Process.Kill()
		c.cmd.Wait()
	}
	os.Remove(socketPath)
}
