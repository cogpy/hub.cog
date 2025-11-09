package opencog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if registry == nil {
		t.Error("Registry should not be nil")
	}

	if registry.Count() != 0 {
		t.Errorf("New registry should have 0 agents, got %d", registry.Count())
	}
}

func TestRegistryRegisterAndGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	// Test registration
	err = registry.Register(agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if registry.Count() != 1 {
		t.Errorf("Registry should have 1 agent, got %d", registry.Count())
	}

	// Test retrieval by ID
	retrieved, err := registry.Get(agent.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.ID != agent.ID {
		t.Errorf("Retrieved agent ID mismatch: expected %s, got %s", agent.ID, retrieved.ID)
	}

	// Test retrieval by name
	retrievedByName, err := registry.GetByName(agent.Name)
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	if retrievedByName.Name != agent.Name {
		t.Errorf("Retrieved agent name mismatch: expected %s, got %s", agent.Name, retrievedByName.Name)
	}
}

func TestRegistryList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Create and register multiple agents
	agents := []AgentConfig{
		{Name: "agent1", Type: AtomSpaceAgent},
		{Name: "agent2", Type: PLNAgent},
		{Name: "agent3", Type: ECANAgent},
	}

	for _, config := range agents {
		agent, err := NewAgent(config)
		if err != nil {
			t.Fatalf("NewAgent failed: %v", err)
		}
		err = registry.Register(agent)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
	}

	// Test list all
	allAgents := registry.List()
	if len(allAgents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(allAgents))
	}
}

func TestRegistryListByType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Create agents of different types
	configs := []AgentConfig{
		{Name: "atomspace1", Type: AtomSpaceAgent},
		{Name: "atomspace2", Type: AtomSpaceAgent},
		{Name: "pln1", Type: PLNAgent},
	}

	for _, config := range configs {
		agent, err := NewAgent(config)
		if err != nil {
			t.Fatalf("NewAgent failed: %v", err)
		}
		err = registry.Register(agent)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
	}

	// Test list by type
	atomspaceAgents := registry.ListByType(AtomSpaceAgent)
	if len(atomspaceAgents) != 2 {
		t.Errorf("Expected 2 atomspace agents, got %d", len(atomspaceAgents))
	}

	plnAgents := registry.ListByType(PLNAgent)
	if len(plnAgents) != 1 {
		t.Errorf("Expected 1 PLN agent, got %d", len(plnAgents))
	}
}

func TestRegistryUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	err = registry.Register(agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Update agent status
	agent.Status = StatusRunning
	err = registry.Update(agent)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrieved, err := registry.Get(agent.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Status != StatusRunning {
		t.Errorf("Expected status 'running', got '%s'", retrieved.Status)
	}
}

func TestRegistryUnregister(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	err = registry.Register(agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Unregister agent
	err = registry.Unregister(agent.ID)
	if err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	if registry.Count() != 0 {
		t.Errorf("Registry should have 0 agents after unregister, got %d", registry.Count())
	}

	// Verify agent is gone
	_, err = registry.Get(agent.ID)
	if err == nil {
		t.Error("Get should fail for unregistered agent")
	}
}

func TestRegistryPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create registry and add agent
	registry1, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	config := AgentConfig{
		Name: "persistent-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	err = registry1.Register(agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Create new registry from same directory
	registry2, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry (2nd) failed: %v", err)
	}

	// Verify agent was loaded
	if registry2.Count() != 1 {
		t.Errorf("Registry should have 1 agent after reload, got %d", registry2.Count())
	}

	retrieved, err := registry2.GetByName("persistent-agent")
	if err != nil {
		t.Fatalf("GetByName failed after reload: %v", err)
	}

	if retrieved.Name != agent.Name {
		t.Errorf("Agent name mismatch after reload: expected %s, got %s", agent.Name, retrieved.Name)
	}

	// Verify file exists
	agentsFile := filepath.Join(tmpDir, "agents.json")
	if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
		t.Error("agents.json file should exist")
	}
}
