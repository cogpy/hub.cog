# OpenCog Multi-Agent Orchestration Workbench

## Overview

The hub.cog tool provides autonomous multi-agent orchestration capabilities for OpenCog cognitive architectures. It enables distributed cognitive systems to coordinate through GitHub's infrastructure.

## Features

### Agent Types

The workbench supports various cognitive agent types:

- **atomspace**: Knowledge representation and storage using OpenCog's AtomSpace
- **pln**: Probabilistic Logic Networks reasoning engine
- **ecan**: Economic Attention Networks for attention allocation
- **openpsi**: Goal-driven behavior and action selection
- **patternminer**: Pattern mining and discovery
- **metalearning**: Meta-learning and optimization
- **reflection**: Self-reflection and monitoring
- **orchestrator**: Multi-agent coordination and management
- **broker**: Message routing and coordination between agents
- **custom**: User-defined custom agents

### Core Capabilities

1. **Agent Lifecycle Management**
   - Create cognitive agents with specific types and configurations
   - Start, stop, and monitor agent states
   - Persistent agent registry with JSON storage

2. **Multi-Agent Orchestration**
   - Message passing between agents
   - Broadcast messaging for knowledge sharing
   - Health monitoring and automatic failure detection

3. **Registry & Persistence**
   - Automatic persistence of agent configurations
   - Query agents by type, status, name, or ID
   - Configuration stored in `~/.config/hub.cog/agents.json`

## Usage

### Creating Agents

```bash
# Create an AtomSpace agent for knowledge storage
$ hub agent create --name my-atomspace --type atomspace

# Create a PLN reasoning agent with repository
$ hub agent create --name my-reasoner --type pln --repo github.com/opencog/pln

# Create an ECAN attention allocation agent
$ hub agent create --name attention-mgr --type ecan
```

### Managing Agents

```bash
# List all agents
$ hub agent list

# List agents with details
$ hub agent list --verbose

# Filter agents by type
$ hub agent list --type atomspace

# Filter agents by status
$ hub agent list --status running
```

### Controlling Agent Lifecycle

```bash
# Start an agent
$ hub agent start my-atomspace

# Stop an agent
$ hub agent stop my-atomspace

# Check agent status (JSON output)
$ hub agent status my-atomspace
```

### Agent Information

```bash
# List available agent types
$ hub agent types

# Remove an agent
$ hub agent remove my-atomspace
```

## Architecture

### Agent Model

Each agent has:
- **ID**: Unique identifier
- **Name**: Human-readable name
- **Type**: Agent type (atomspace, pln, ecan, etc.)
- **Status**: Current state (created, running, stopped, error)
- **Repository**: Optional Git repository URL
- **Branch**: Git branch (default: main)
- **Config**: Custom configuration as key-value pairs
- **Metrics**: Performance and health metrics
- **Timestamps**: Created, updated, started, stopped times

### Registry

The registry provides:
- Thread-safe agent storage and retrieval
- Automatic persistence to JSON file
- Query capabilities by ID, name, type, and status
- Atomic updates and deletions

### Orchestrator

The orchestrator enables:
- Agent registration and lifecycle management
- Inter-agent message passing
- Broadcast messaging for knowledge sharing
- Health monitoring with periodic checks
- Automatic failure detection (30s heartbeat timeout)

### Message Types

Agents can exchange different message types:
- **command**: Action requests
- **query**: Information requests
- **response**: Query responses
- **knowledge**: Knowledge sharing
- **heartbeat**: Health checks
- **error**: Error notifications

## Integration with OpenCog

This workbench is designed to integrate with OpenCog cognitive architectures:

1. **AtomSpace Integration**: Agents can share a distributed AtomSpace for knowledge representation
2. **PLN Reasoning**: Coordinate probabilistic logic reasoning across multiple agents
3. **ECAN Attention**: Distribute attention allocation across cognitive components
4. **Meta-Cognitive Coordination**: Enable meta-learning and self-optimization

## Storage

Agent configurations are stored in:
```
~/.config/hub.cog/agents.json
```

This file contains all registered agents and is automatically loaded on startup.

## Example Workflow

```bash
# 1. Create a multi-agent cognitive system
$ hub agent create --name knowledge-base --type atomspace
$ hub agent create --name reasoner --type pln
$ hub agent create --name attention --type ecan
$ hub agent create --name orchestrator --type orchestrator

# 2. Start all agents
$ hub agent start knowledge-base
$ hub agent start reasoner
$ hub agent start attention
$ hub agent start orchestrator

# 3. Monitor system status
$ hub agent list --status running

# 4. Check individual agent details
$ hub agent status knowledge-base

# 5. Stop agents when done
$ hub agent stop knowledge-base
$ hub agent stop reasoner
$ hub agent stop attention
$ hub agent stop orchestrator
```

## Development

### Running Tests

```bash
# Test the opencog package
$ go test ./opencog/... -v

# Run all tests
$ make test
```

### Building from Source

```bash
# Build the binary
$ make bin/hub

# Or use go build directly
$ go build -o bin/hub
```

## Future Enhancements

Potential future features:
- Docker/Kubernetes deployment of agents
- RESTful API for remote agent management
- WebSocket-based real-time communication
- Prometheus metrics export
- Agent dependency management
- Workflow automation and pipelines
- Integration with GitHub Actions for CI/CD
- Distributed consensus mechanisms
- Knowledge graph visualization
- Performance profiling and optimization

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

## License

See LICENSE file for details.
