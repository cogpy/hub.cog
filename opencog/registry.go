package opencog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Registry manages the collection of cognitive agents
type Registry struct {
	agents map[string]*Agent
	mu     sync.RWMutex
	file   string
}

// NewRegistry creates a new agent registry
func NewRegistry(configDir string) (*Registry, error) {
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config", "hub.cog")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	registry := &Registry{
		agents: make(map[string]*Agent),
		file:   filepath.Join(configDir, "agents.json"),
	}

	// Load existing agents from file
	if err := registry.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load agents: %w", err)
	}

	return registry, nil
}

// Register adds a new agent to the registry
func (r *Registry) Register(agent *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agent.ID]; exists {
		return fmt.Errorf("agent with ID %s already exists", agent.ID)
	}

	r.agents[agent.ID] = agent
	return r.save()
}

// Get retrieves an agent by ID
func (r *Registry) Get(id string) (*Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}

	return agent, nil
}

// GetByName retrieves an agent by name
func (r *Registry) GetByName(name string) (*Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, agent := range r.agents {
		if agent.Name == name {
			return agent, nil
		}
	}

	return nil, fmt.Errorf("agent with name %s not found", name)
}

// List returns all registered agents
func (r *Registry) List() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}

	return agents
}

// ListByType returns agents of a specific type
func (r *Registry) ListByType(agentType AgentType) []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*Agent, 0)
	for _, agent := range r.agents {
		if agent.Type == agentType {
			agents = append(agents, agent)
		}
	}

	return agents
}

// ListByStatus returns agents with a specific status
func (r *Registry) ListByStatus(status AgentStatus) []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*Agent, 0)
	for _, agent := range r.agents {
		if agent.Status == status {
			agents = append(agents, agent)
		}
	}

	return agents
}

// Update updates an existing agent in the registry
func (r *Registry) Update(agent *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agent.ID]; !exists {
		return fmt.Errorf("agent with ID %s not found", agent.ID)
	}

	r.agents[agent.ID] = agent
	return r.save()
}

// Unregister removes an agent from the registry
func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[id]; !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}

	delete(r.agents, id)
	return r.save()
}

// Count returns the total number of agents
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agents)
}

// load reads agents from the JSON file
func (r *Registry) load() error {
	data, err := os.ReadFile(r.file)
	if err != nil {
		return err
	}

	var agents []*Agent
	if err := json.Unmarshal(data, &agents); err != nil {
		return fmt.Errorf("failed to unmarshal agents: %w", err)
	}

	for _, agent := range agents {
		r.agents[agent.ID] = agent
	}

	return nil
}

// save writes agents to the JSON file
func (r *Registry) save() error {
	agents := make([]*Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}

	data, err := json.MarshalIndent(agents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agents: %w", err)
	}

	if err := os.WriteFile(r.file, data, 0644); err != nil {
		return fmt.Errorf("failed to write agents file: %w", err)
	}

	return nil
}
