package bootstrap

import (
	"github.com/artcodefun/detective-game/backend/internal/application/commands"
)

type Commands struct {
	Scenario      *commands.ScenarioCommands
	Interrogation *commands.InterrogationCommands
	Evaluation    *commands.EvaluationCommands
	Actions       *commands.ActionCommands
	Notebook      *commands.NotebookCommands
}

func NewCommands(a *Adapters) *Commands {
	return &Commands{
		Scenario:      commands.NewScenarioCommands(a.Users, a.Sessions, a.LLM, a.Characters, a.Evidence, a.Prototypes),
		Interrogation: commands.NewInterrogationCommands(a.Sessions, a.Interrogations, a.Characters, a.Chat, a.LLM),
		Evaluation:    commands.NewEvaluationCommands(a.Sessions, a.LLM),
		Actions:       commands.NewActionCommands(a.Sessions, a.Reports, a.LLM),
		Notebook:      commands.NewNotebookCommands(a.Chronology),
	}
}
