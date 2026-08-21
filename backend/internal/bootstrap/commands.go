package bootstrap

import (
	"github.com/artcodefun/detective-game/backend/internal/application/commands"
)

type Commands struct {
	User          *commands.UserCommands
	Scenario      *commands.ScenarioCommands
	Interrogation *commands.InterrogationCommands
	Evaluation    *commands.EvaluationCommands
	Actions       *commands.ActionCommands
	Notebook      *commands.NotebookCommands
}

func NewCommands(a *Adapters) *Commands {
	return &Commands{
		User:          commands.NewUserCommands(a.Users),
		Scenario:      commands.NewScenarioCommands(a.Sessions, a.LLM, a.Characters, a.Evidence, a.Chronology),
		Interrogation: commands.NewInterrogationCommands(a.Sessions, a.Interrogations, a.Characters, a.Chat, a.LLM, a.Chronology),
		Evaluation:    commands.NewEvaluationCommands(a.Sessions, a.Interrogations, a.LLM, a.Chronology),
		Actions:       commands.NewActionCommands(a.Sessions, a.Reports, a.Evidence, a.Characters, a.LLM, a.Chronology),
		Notebook:      commands.NewNotebookCommands(a.Chronology),
	}
}
