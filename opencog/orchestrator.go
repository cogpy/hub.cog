package opencog

import (
	"fmt"
	"sync"
	"time"
)

// Orchestrator manages multi-agent coordination and communication
type Orchestrator struct {
	registry  *Registry
	channels  map[string]chan *Message
	mu        sync.RWMutex
	running   bool
	stopCh    chan struct{}
}

// Message represents communication between agents
type Message struct {
	ID        string                 `json:"id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Type      MessageType            `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}

// MessageType defines types of inter-agent messages
type MessageType string

const (
	MessageTypeCommand   MessageType = "command"
	MessageTypeQuery     MessageType = "query"
	MessageTypeResponse  MessageType = "response"
	MessageTypeKnowledge MessageType = "knowledge"
	MessageTypeHeartbeat MessageType = "heartbeat"
	MessageTypeError     MessageType = "error"
)

// NewOrchestrator creates a new multi-agent orchestrator
func NewOrchestrator(registry *Registry) *Orchestrator {
	return &Orchestrator{
		registry: registry,
		channels: make(map[string]chan *Message),
		stopCh:   make(chan struct{}),
	}
}

// Start begins the orchestrator's coordination loop
func (o *Orchestrator) Start() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.running {
		return fmt.Errorf("orchestrator is already running")
	}

	o.running = true
	
	// Start coordination goroutine
	go o.coordinationLoop()

	return nil
}

// Stop halts the orchestrator
func (o *Orchestrator) Stop() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.running {
		return fmt.Errorf("orchestrator is not running")
	}

	o.running = false
	close(o.stopCh)

	// Close all agent channels
	for _, ch := range o.channels {
		close(ch)
	}
	o.channels = make(map[string]chan *Message)

	return nil
}

// RegisterAgent registers an agent with the orchestrator
func (o *Orchestrator) RegisterAgent(agentID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.channels[agentID]; exists {
		return fmt.Errorf("agent %s is already registered", agentID)
	}

	o.channels[agentID] = make(chan *Message, 100)
	return nil
}

// UnregisterAgent removes an agent from the orchestrator
func (o *Orchestrator) UnregisterAgent(agentID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	ch, exists := o.channels[agentID]
	if !exists {
		return fmt.Errorf("agent %s is not registered", agentID)
	}

	close(ch)
	delete(o.channels, agentID)
	return nil
}

// SendMessage sends a message from one agent to another
func (o *Orchestrator) SendMessage(msg *Message) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	ch, exists := o.channels[msg.To]
	if !exists {
		return fmt.Errorf("agent %s is not registered", msg.To)
	}

	msg.Timestamp = time.Now()
	if msg.ID == "" {
		msg.ID = generateMessageID()
	}

	select {
	case ch <- msg:
		return nil
	default:
		return fmt.Errorf("agent %s message queue is full", msg.To)
	}
}

// BroadcastMessage sends a message to all registered agents
func (o *Orchestrator) BroadcastMessage(from string, msgType MessageType, payload map[string]interface{}) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	msg := &Message{
		ID:        generateMessageID(),
		From:      from,
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	for agentID, ch := range o.channels {
		if agentID == from {
			continue // Don't send to self
		}

		msgCopy := *msg
		msgCopy.To = agentID

		select {
		case ch <- &msgCopy:
			// Message sent
		default:
			// Queue full, skip this agent
			fmt.Printf("Warning: message queue full for agent %s\n", agentID)
		}
	}

	return nil
}

// coordinationLoop is the main coordination routine
func (o *Orchestrator) coordinationLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-o.stopCh:
			return
		case <-ticker.C:
			o.performHealthChecks()
		}
	}
}

// performHealthChecks checks the health of all agents
func (o *Orchestrator) performHealthChecks() {
	agents := o.registry.List()
	
	for _, agent := range agents {
		if agent.Status == StatusRunning {
			// Check if agent is still responding
			if agent.Metrics != nil {
				timeSinceHeartbeat := time.Since(agent.Metrics.LastHeartbeat)
				if timeSinceHeartbeat > 30*time.Second {
					// Agent may be unresponsive
					agent.Status = StatusError
					agent.UpdatedAt = time.Now()
					o.registry.Update(agent)
				}
			}
		}
	}
}

// GetAgentChannel returns the message channel for an agent
func (o *Orchestrator) GetAgentChannel(agentID string) (chan *Message, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	ch, exists := o.channels[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s is not registered", agentID)
	}

	return ch, nil
}

// generateMessageID generates a unique message identifier
func generateMessageID() string {
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}
