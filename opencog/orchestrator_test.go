package opencog

import (
	"testing"
	"time"
)

func TestNewOrchestrator(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)
	if orchestrator == nil {
		t.Error("Orchestrator should not be nil")
	}

	if orchestrator.registry != registry {
		t.Error("Orchestrator registry mismatch")
	}
}

func TestOrchestratorStartStop(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)

	// Test start
	err = orchestrator.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !orchestrator.running {
		t.Error("Orchestrator should be running after Start()")
	}

	// Test double start
	err = orchestrator.Start()
	if err == nil {
		t.Error("Second Start() should return an error")
	}

	// Test stop
	err = orchestrator.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if orchestrator.running {
		t.Error("Orchestrator should not be running after Stop()")
	}

	// Test double stop
	err = orchestrator.Stop()
	if err == nil {
		t.Error("Second Stop() should return an error")
	}
}

func TestOrchestratorRegisterAgent(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)

	// Register an agent
	agentID := "test-agent-1"
	err = orchestrator.RegisterAgent(agentID)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	// Verify channel was created
	_, err = orchestrator.GetAgentChannel(agentID)
	if err != nil {
		t.Errorf("GetAgentChannel failed: %v", err)
	}

	// Test double registration
	err = orchestrator.RegisterAgent(agentID)
	if err == nil {
		t.Error("Second RegisterAgent should return an error")
	}
}

func TestOrchestratorUnregisterAgent(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)

	agentID := "test-agent-1"
	err = orchestrator.RegisterAgent(agentID)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	// Unregister agent
	err = orchestrator.UnregisterAgent(agentID)
	if err != nil {
		t.Fatalf("UnregisterAgent failed: %v", err)
	}

	// Verify channel was removed
	_, err = orchestrator.GetAgentChannel(agentID)
	if err == nil {
		t.Error("GetAgentChannel should fail after UnregisterAgent")
	}

	// Test double unregistration
	err = orchestrator.UnregisterAgent(agentID)
	if err == nil {
		t.Error("Second UnregisterAgent should return an error")
	}
}

func TestOrchestratorSendMessage(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)

	// Register agents
	agent1 := "agent-1"
	agent2 := "agent-2"
	orchestrator.RegisterAgent(agent1)
	orchestrator.RegisterAgent(agent2)

	// Send message
	msg := &Message{
		From:    agent1,
		To:      agent2,
		Type:    MessageTypeCommand,
		Payload: map[string]interface{}{"command": "test"},
	}

	err = orchestrator.SendMessage(msg)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Verify message was received
	ch, _ := orchestrator.GetAgentChannel(agent2)
	select {
	case receivedMsg := <-ch:
		if receivedMsg.From != agent1 {
			t.Errorf("Message from mismatch: expected %s, got %s", agent1, receivedMsg.From)
		}
		if receivedMsg.To != agent2 {
			t.Errorf("Message to mismatch: expected %s, got %s", agent2, receivedMsg.To)
		}
		if receivedMsg.Type != MessageTypeCommand {
			t.Errorf("Message type mismatch: expected %s, got %s", MessageTypeCommand, receivedMsg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestOrchestratorBroadcastMessage(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	orchestrator := NewOrchestrator(registry)

	// Register multiple agents
	agents := []string{"agent-1", "agent-2", "agent-3"}
	for _, agentID := range agents {
		orchestrator.RegisterAgent(agentID)
	}

	// Broadcast message
	sender := "agent-1"
	payload := map[string]interface{}{"broadcast": "test"}
	err = orchestrator.BroadcastMessage(sender, MessageTypeKnowledge, payload)
	if err != nil {
		t.Fatalf("BroadcastMessage failed: %v", err)
	}

	// Verify all agents except sender received the message
	receivedCount := 0
	for _, agentID := range agents {
		if agentID == sender {
			continue
		}

		ch, _ := orchestrator.GetAgentChannel(agentID)
		select {
		case msg := <-ch:
			if msg.From != sender {
				t.Errorf("Message from mismatch for %s: expected %s, got %s", agentID, sender, msg.From)
			}
			if msg.Type != MessageTypeKnowledge {
				t.Errorf("Message type mismatch for %s", agentID)
			}
			receivedCount++
		case <-time.After(1 * time.Second):
			t.Errorf("Timeout waiting for broadcast message to %s", agentID)
		}
	}

	expectedCount := len(agents) - 1 // All except sender
	if receivedCount != expectedCount {
		t.Errorf("Expected %d agents to receive broadcast, got %d", expectedCount, receivedCount)
	}
}

func TestMessageTypes(t *testing.T) {
	types := []MessageType{
		MessageTypeCommand,
		MessageTypeQuery,
		MessageTypeResponse,
		MessageTypeKnowledge,
		MessageTypeHeartbeat,
		MessageTypeError,
	}

	for _, msgType := range types {
		msg := &Message{
			From:    "agent-1",
			To:      "agent-2",
			Type:    msgType,
			Payload: map[string]interface{}{},
		}

		if msg.Type != msgType {
			t.Errorf("Message type mismatch: expected %s, got %s", msgType, msg.Type)
		}
	}
}

func TestMessageIDGeneration(t *testing.T) {
	id1 := generateMessageID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	id2 := generateMessageID()

	if id1 == "" || id2 == "" {
		t.Error("Generated message IDs should not be empty")
	}

	if id1 == id2 {
		t.Error("Generated message IDs should be unique")
	}
}
