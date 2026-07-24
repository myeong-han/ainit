package command

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/myeong-han/ainit/pkg/config"
	"github.com/myeong-han/ainit/pkg/generator"
	"github.com/myeong-han/ainit/pkg/git"
)

type ActionType int

const (
	ActionNone ActionType = iota
	ActionSetConfs
	ActionGenDocs
	ActionGenCodes
	ActionGitInit
	ActionHelp
)

type Result struct {
	Action  ActionType
	Message string
}

type SlashOption struct {
	Name        string
	Description string
	Example     string
}

type CommandEngine struct {
	cfg *config.Config
}

func NewCommandEngine(cfg *config.Config) *CommandEngine {
	return &CommandEngine{cfg: cfg}
}

// GetAvailableSlashCommands returns the list of supported slash commands for UI autocomplete dropdown
func GetAvailableSlashCommands() []SlashOption {
	return []SlashOption{
		{Name: "/git-init", Description: "Clone existing repo or initialize new git repository & update work-dir", Example: "/git-init myeong-han/ainit"},
		{Name: "/set-confs", Description: "Configure AI Provider, Arch, Git & Conventions", Example: "/set-confs --provider openai --arch msa"},
		{Name: "/gen-docs", Description: "Generate Architecture Spec & Mermaid Diagrams", Example: "/gen-docs"},
		{Name: "/gen-codes", Description: "Generate Agent Context Rules & Scaffolding", Example: "/gen-codes"},
		{Name: "/help", Description: "Show available Slash Commands & Usage", Example: "/help"},
	}
}

// IsSlashCommand returns true if the input string starts with '/'
func IsSlashCommand(input string) bool {
	return strings.HasPrefix(strings.TrimSpace(input), "/")
}

// Execute parses and runs a slash command string
func (e *CommandEngine) Execute(input string) (*Result, error) {
	trimmed := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmed, "/") {
		return nil, errors.New("command must start with '/'")
	}

	parts := strings.Fields(trimmed)
	cmdName := strings.ToLower(parts[0])

	switch cmdName {
	case "/git-init":
		return e.handleGitInit(parts[1:])
	case "/set-confs":
		return e.handleSetConfs(parts[1:])
	case "/gen-docs":
		return e.handleGenDocs()
	case "/gen-codes":
		return e.handleGenCodes()
	case "/help":
		return e.handleHelp()
	default:
		return nil, fmt.Errorf("unknown slash command '%s'. Type '/help' for available commands", cmdName)
	}
}

func (e *CommandEngine) handleGitInit(args []string) (*Result, error) {
	repoPath := "myeong-han/ainit"
	if len(args) > 0 {
		repoPath = args[0]
	}

	provider := e.cfg.Step2.GitProvider
	if provider == "" {
		provider = "github"
	}

	repoName := repoPath
	if parts := strings.Split(repoPath, "/"); len(parts) == 2 {
		repoName = parts[1]
	}
	targetDir := filepath.Join(".", repoName)

	res, err := git.InitOrCloneRepository(provider, repoPath, targetDir)
	if err != nil {
		return nil, fmt.Errorf("git-init failed: %v", err)
	}

	msg := fmt.Sprintf("🐙 [%s] %s\nUpdated Working Directory: %s", strings.ToUpper(res.Action), res.Message, res.WorkDir)

	return &Result{
		Action:  ActionGitInit,
		Message: msg,
	}, nil
}

func (e *CommandEngine) handleSetConfs(args []string) (*Result, error) {
	var updated []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") && i+1 < len(args) {
			key := strings.TrimPrefix(arg, "--")
			val := args[i+1]
			i++

			switch key {
			case "provider":
				e.cfg.Step0.ProviderID = val
				updated = append(updated, "provider="+val)
			case "model":
				e.cfg.Step0.PrimaryModel = val
				updated = append(updated, "model="+val)
			case "arch":
				e.cfg.Step1.ArchitectureStyle = val
				updated = append(updated, "arch="+val)
			case "repo":
				e.cfg.Step1.RepoStructure = val
				updated = append(updated, "repo="+val)
			case "git":
				e.cfg.Step2.GitProvider = val
				updated = append(updated, "git="+val)
			case "commit":
				e.cfg.Step3.CommitConvention = val
				updated = append(updated, "commit="+val)
			case "tdd":
				e.cfg.Step3.TDDMode = (val == "true" || val == "1" || val == "yes")
				updated = append(updated, fmt.Sprintf("tdd=%v", e.cfg.Step3.TDDMode))
			}
		}
	}

	msg := "Updated configs: " + strings.Join(updated, ", ")
	if len(updated) == 0 {
		msg = "No configs specified. Usage: /set-confs --provider anthropic --arch msa --git github --commit conventional"
	}

	return &Result{
		Action:  ActionSetConfs,
		Message: msg,
	}, nil
}

func (e *CommandEngine) handleGenDocs() (*Result, error) {
	err := generator.GenerateHarnessProject(".", e.cfg, "Architecture Spec & Mermaid diagrams generated via /gen-docs slash command")
	if err != nil {
		return nil, fmt.Errorf("failed to generate docs: %v", err)
	}

	return &Result{
		Action:  ActionGenDocs,
		Message: "📄 Successfully generated docs/ARCHITECTURE_SPEC.md with Mermaid diagrams!",
	}, nil
}

func (e *CommandEngine) handleGenCodes() (*Result, error) {
	err := generator.GenerateAgentContextFiles(".")
	if err != nil {
		return nil, fmt.Errorf("failed to generate codes & rules: %v", err)
	}

	return &Result{
		Action:  ActionGenCodes,
		Message: "💻 Successfully generated AGENTS.md, CLAUDE.md, .cursorrules & agent context rules!",
	}, nil
}

func (e *CommandEngine) handleHelp() (*Result, error) {
	helpMsg := `Available Slash Commands:
• /git-init [owner/repo] : Clone existing repo or initialize new git repository & update work-dir
• /set-confs --provider <id> --arch <msa|monolith> --git <github|bitbucket> --commit <conventional|gitmoji>
• /gen-docs  : Generate docs/ARCHITECTURE_SPEC.md & 4 Mermaid diagrams
• /gen-codes : Generate AGENTS.md, CLAUDE.md & cross-agent context rules
• /help      : Show available slash commands`

	return &Result{
		Action:  ActionHelp,
		Message: helpMsg,
	}, nil
}
