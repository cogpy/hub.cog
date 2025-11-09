package opencog

import (
	"encoding/json"
	"fmt"
	"time"
)

// AgentType represents different types of cognitive agents
type AgentType string

const (
	// Core cognitive agent types
	AtomSpaceAgent    AgentType = "atomspace"    // Knowledge representation and storage
	PLNAgent          AgentType = "pln"          // Probabilistic Logic Networks reasoning
	ECANAgent         AgentType = "ecan"         // Economic Attention Networks
	OpenPsiAgent      AgentType = "openpsi"      // Goal-driven behavior
	PatternMinerAgent AgentType = "patternminer" // Pattern mining and discovery
	
	// Meta-cognitive agents
	MetaLearningAgent AgentType = "metalearning"  // Meta-learning and optimization
	ReflectionAgent   AgentType = "reflection"    // Self-reflection and monitoring
	
	// Coordination agents
	OrchestratorAgent AgentType = "orchestrator"  // Multi-agent coordination
	BrokerAgent       AgentType = "broker"        // Message routing and coordination
	
	// Custom agent type
	CustomAgent AgentType = "custom" // User-defined agents
)

// AgentStatus represents the current state of an agent
type AgentStatus string

const (
	StatusCreated   AgentStatus = "created"
	StatusStarting  AgentStatus = "starting"
	StatusRunning   AgentStatus = "running"
	StatusPaused    AgentStatus = "paused"
	StatusStopping  AgentStatus = "stopping"
	StatusStopped   AgentStatus = "stopped"
	StatusError     AgentStatus = "error"
)

// Agent represents a cognitive agent in the OpenCog system
type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        AgentType              `json:"type"`
	Status      AgentStatus            `json:"status"`
	Repository  string                 `json:"repository"`
	Branch      string                 `json:"branch"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	StoppedAt   *time.Time             `json:"stopped_at,omitempty"`
	Endpoint    string                 `json:"endpoint,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metrics     *AgentMetrics          `json:"metrics,omitempty"`
}

// AgentMetrics contains performance and health metrics for an agent
type AgentMetrics struct {
	CPUUsage      float64   `json:"cpu_usage"`
	MemoryUsage   int64     `json:"memory_usage"`
	RequestCount  int64     `json:"request_count"`
	ErrorCount    int64     `json:"error_count"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Uptime        int64     `json:"uptime"` // seconds
}

// AgentConfig defines configuration options for creating an agent
type AgentConfig struct {
	Name       string                 `json:"name"`
	Type       AgentType              `json:"type"`
	Repository string                 `json:"repository,omitempty"`
	Branch     string                 `json:"branch,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// Validate checks if the agent configuration is valid
func (ac *AgentConfig) Validate() error {
	if ac.Name == "" {
		return fmt.Errorf("agent name is required")
	}
	if ac.Type == "" {
		return fmt.Errorf("agent type is required")
	}
	return nil
}

// NewAgent creates a new agent instance
func NewAgent(config AgentConfig) (*Agent, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	agent := &Agent{
		ID:         generateAgentID(),
		Name:       config.Name,
		Type:       config.Type,
		Status:     StatusCreated,
		Repository: config.Repository,
		Branch:     config.Branch,
		Config:     config.Config,
		CreatedAt:  now,
		UpdatedAt:  now,
		Tags:       config.Tags,
	}

	return agent, nil
}

// ToJSON converts the agent to JSON string
func (a *Agent) ToJSON() (string, error) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// generateAgentID generates a unique agent identifier
func generateAgentID() string {
	return fmt.Sprintf("agent-%d", time.Now().UnixNano())
}
