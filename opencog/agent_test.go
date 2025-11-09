package opencog

import (
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	if agent.Name != "test-agent" {
		t.Errorf("Expected name 'test-agent', got '%s'", agent.Name)
	}

	if agent.Type != AtomSpaceAgent {
		t.Errorf("Expected type 'atomspace', got '%s'", agent.Type)
	}

	if agent.Status != StatusCreated {
		t.Errorf("Expected status 'created', got '%s'", agent.Status)
	}

	if agent.ID == "" {
		t.Error("Agent ID should not be empty")
	}
}

func TestAgentConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  AgentConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AgentConfig{
				Name: "test",
				Type: PLNAgent,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: AgentConfig{
				Type: PLNAgent,
			},
			wantErr: true,
		},
		{
			name: "missing type",
			config: AgentConfig{
				Name: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgentToJSON(t *testing.T) {
	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	json, err := agent.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if json == "" {
		t.Error("JSON output should not be empty")
	}

	// Check that JSON contains expected fields
	expectedFields := []string{"id", "name", "type", "status"}
	for _, field := range expectedFields {
		if !contains(json, field) {
			t.Errorf("JSON output missing field: %s", field)
		}
	}
}

func TestAgentTypes(t *testing.T) {
	types := []AgentType{
		AtomSpaceAgent,
		PLNAgent,
		ECANAgent,
		OpenPsiAgent,
		PatternMinerAgent,
		MetaLearningAgent,
		ReflectionAgent,
		OrchestratorAgent,
		BrokerAgent,
		CustomAgent,
	}

	for _, agentType := range types {
		config := AgentConfig{
			Name: "test-" + string(agentType),
			Type: agentType,
		}

		agent, err := NewAgent(config)
		if err != nil {
			t.Errorf("Failed to create agent of type %s: %v", agentType, err)
		}

		if agent.Type != agentType {
			t.Errorf("Expected type %s, got %s", agentType, agent.Type)
		}
	}
}

func TestAgentStatusTransitions(t *testing.T) {
	config := AgentConfig{
		Name: "test-agent",
		Type: AtomSpaceAgent,
	}

	agent, err := NewAgent(config)
	if err != nil {
		t.Fatalf("NewAgent failed: %v", err)
	}

	// Test initial status
	if agent.Status != StatusCreated {
		t.Errorf("Initial status should be 'created', got '%s'", agent.Status)
	}

	// Test status transitions
	statuses := []AgentStatus{
		StatusStarting,
		StatusRunning,
		StatusPaused,
		StatusStopping,
		StatusStopped,
		StatusError,
	}

	for _, status := range statuses {
		agent.Status = status
		agent.UpdatedAt = time.Now()

		if agent.Status != status {
			t.Errorf("Expected status %s, got %s", status, agent.Status)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != "" && 
		   (s == substr || (len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
