package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/github/hub/v2/opencog"
	"github.com/github/hub/v2/ui"
)

var cmdAgent = &Command{
	Run:   agent,
	Usage: "agent <command> [options]",
	Long: `Manage OpenCog cognitive agents for autonomous multi-agent orchestration.

## Commands:

	create     Create a new cognitive agent
	list       List all registered agents
	start      Start an agent
	stop       Stop an agent
	status     Show agent status and metrics
	remove     Remove an agent
	types      List available agent types

## Examples:

	# Create a new AtomSpace agent
	$ hub agent create --name my-atomspace --type atomspace

	# List all agents
	$ hub agent list

	# Show status of a specific agent
	$ hub agent status my-atomspace

	# Start an agent
	$ hub agent start my-atomspace

	# Stop an agent
	$ hub agent stop my-atomspace

	# Remove an agent
	$ hub agent remove my-atomspace

	# List available agent types
	$ hub agent types
`,
	KnownFlags: `
	--name <NAME>
		Agent name (required for create)

	--type <TYPE>
		Agent type (required for create)

	--repo <URL>
		Git repository URL (optional for create)

	--branch <BRANCH>
		Git branch (optional for create, default: main)

	--verbose
		Show detailed information (for list)
`,
}

func init() {
	CmdRunner.Use(cmdAgent)
}

func agent(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		ui.Errorln("Usage: hub agent <command> [options]")
		ui.Errorln("Run 'hub agent --help' for more information")
		os.Exit(1)
	}

	subCommand := args.FirstParam()
	args.Params = args.Params[1:]

	switch subCommand {
	case "create":
		agentCreate(cmd, args)
	case "list", "ls":
		agentList(cmd, args)
	case "start":
		agentStart(cmd, args)
	case "stop":
		agentStop(cmd, args)
	case "status":
		agentStatus(cmd, args)
	case "remove", "rm":
		agentRemove(cmd, args)
	case "types":
		agentTypes(cmd, args)
	default:
		ui.Errorf("Unknown agent command: %s\n", subCommand)
		ui.Errorln("Run 'hub agent --help' for usage information")
		os.Exit(1)
	}
}

func agentCreate(cmd *Command, args *Args) {
	name := args.Flag.Value("--name")
	agentType := args.Flag.Value("--type")
	repository := args.Flag.Value("--repo")
	branch := args.Flag.Value("--branch")
	if branch == "" {
		branch = "main"
	}

	if name == "" {
		ui.Errorln("Error: --name is required")
		os.Exit(1)
	}

	if agentType == "" {
		ui.Errorln("Error: --type is required")
		os.Exit(1)
	}

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	config := opencog.AgentConfig{
		Name:       name,
		Type:       opencog.AgentType(agentType),
		Repository: repository,
		Branch:     branch,
		Config:     make(map[string]interface{}),
	}

	agent, err := opencog.NewAgent(config)
	if err != nil {
		ui.Errorf("Error: failed to create agent: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Register(agent); err != nil {
		ui.Errorf("Error: failed to register agent: %v\n", err)
		os.Exit(1)
	}

	ui.Printf("Created agent: %s (ID: %s, Type: %s)\n", agent.Name, agent.ID, agent.Type)
}

func agentList(cmd *Command, args *Args) {
	typeFilter := args.Flag.Value("--type")
	statusFilter := args.Flag.Value("--status")
	verbose := args.Flag.Bool("--verbose")

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	var agents []*opencog.Agent
	if typeFilter != "" {
		agents = registry.ListByType(opencog.AgentType(typeFilter))
	} else if statusFilter != "" {
		agents = registry.ListByStatus(opencog.AgentStatus(statusFilter))
	} else {
		agents = registry.List()
	}

	if len(agents) == 0 {
		ui.Println("No agents found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if verbose {
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tSTATUS\tCREATED\tREPOSITORY")
		for _, agent := range agents {
			repo := agent.Repository
			if repo == "" {
				repo = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				agent.ID, agent.Name, agent.Type, agent.Status,
				agent.CreatedAt.Format("2006-01-02 15:04:05"), repo)
		}
	} else {
		fmt.Fprintln(w, "NAME\tTYPE\tSTATUS")
		for _, agent := range agents {
			fmt.Fprintf(w, "%s\t%s\t%s\n", agent.Name, agent.Type, agent.Status)
		}
	}
	w.Flush()
}

func agentStart(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		ui.Errorln("Error: agent name is required")
		ui.Errorln("Usage: hub agent start <name>")
		os.Exit(1)
	}

	agentName := args.FirstParam()

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	agent, err := registry.GetByName(agentName)
	if err != nil {
		ui.Errorf("Error: %v\n", err)
		os.Exit(1)
	}

	if agent.Status == opencog.StatusRunning {
		ui.Printf("Agent %s is already running\n", agentName)
		return
	}

	agent.Status = opencog.StatusRunning
	now := time.Now()
	agent.StartedAt = &now
	agent.UpdatedAt = now

	if err := registry.Update(agent); err != nil {
		ui.Errorf("Error: failed to update agent: %v\n", err)
		os.Exit(1)
	}

	ui.Printf("Started agent: %s\n", agentName)
}

func agentStop(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		ui.Errorln("Error: agent name is required")
		ui.Errorln("Usage: hub agent stop <name>")
		os.Exit(1)
	}

	agentName := args.FirstParam()

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	agent, err := registry.GetByName(agentName)
	if err != nil {
		ui.Errorf("Error: %v\n", err)
		os.Exit(1)
	}

	if agent.Status == opencog.StatusStopped {
		ui.Printf("Agent %s is already stopped\n", agentName)
		return
	}

	agent.Status = opencog.StatusStopped
	now := time.Now()
	agent.StoppedAt = &now
	agent.UpdatedAt = now

	if err := registry.Update(agent); err != nil {
		ui.Errorf("Error: failed to update agent: %v\n", err)
		os.Exit(1)
	}

	ui.Printf("Stopped agent: %s\n", agentName)
}

func agentStatus(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		ui.Errorln("Error: agent name is required")
		ui.Errorln("Usage: hub agent status <name>")
		os.Exit(1)
	}

	agentName := args.FirstParam()

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	agent, err := registry.GetByName(agentName)
	if err != nil {
		ui.Errorf("Error: %v\n", err)
		os.Exit(1)
	}

	jsonStr, err := agent.ToJSON()
	if err != nil {
		ui.Errorf("Error: failed to convert agent to JSON: %v\n", err)
		os.Exit(1)
	}

	ui.Println(jsonStr)
}

func agentRemove(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		ui.Errorln("Error: agent name is required")
		ui.Errorln("Usage: hub agent remove <name>")
		os.Exit(1)
	}

	agentName := args.FirstParam()

	registry, err := opencog.NewRegistry("")
	if err != nil {
		ui.Errorf("Error: failed to create registry: %v\n", err)
		os.Exit(1)
	}

	agent, err := registry.GetByName(agentName)
	if err != nil {
		ui.Errorf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Unregister(agent.ID); err != nil {
		ui.Errorf("Error: failed to remove agent: %v\n", err)
		os.Exit(1)
	}

	ui.Printf("Removed agent: %s\n", agentName)
}

func agentTypes(cmd *Command, args *Args) {
	types := []struct {
		Type        string
		Description string
	}{
		{"atomspace", "Knowledge representation and storage"},
		{"pln", "Probabilistic Logic Networks reasoning"},
		{"ecan", "Economic Attention Networks"},
		{"openpsi", "Goal-driven behavior"},
		{"patternminer", "Pattern mining and discovery"},
		{"metalearning", "Meta-learning and optimization"},
		{"reflection", "Self-reflection and monitoring"},
		{"orchestrator", "Multi-agent coordination"},
		{"broker", "Message routing and coordination"},
		{"custom", "User-defined agents"},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tDESCRIPTION")
	for _, t := range types {
		fmt.Fprintf(w, "%s\t%s\n", t.Type, t.Description)
	}
	w.Flush()
}
