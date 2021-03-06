package commands

import (
	"flag"
	"fmt"
	"io"

	"github.com/martinrue/cosmo/config"
)

const maxCommandLength = 50

// CommandSteps lists the steps of a task.
type CommandSteps struct {
	TaskName string
	Task     config.Task
	Writer   io.Writer
}

func (cmd *CommandSteps) printSteps(serverName string, steps []config.Step) {
	if len(steps) == 0 {
		return
	}

	fmt.Fprintf(cmd.Writer, "\n%s\n", serverName)

	for _, step := range steps {
		fmt.Fprintf(cmd.Writer, "  %s\n", step.Exec)
	}
}

// Exec runs the subcommand.
func (cmd *CommandSteps) Exec() error {
	fmt.Fprintf(cmd.Writer, "Steps for task '%s' on server '%s':\n", cmd.TaskName, cmd.Task.ServerName)

	cmd.printSteps("local", cmd.Task.Local)
	cmd.printSteps("remote", cmd.Task.Remote)

	return nil
}

// NewCommandSteps creates a new 'steps' subcommand.
func NewCommandSteps(config config.Config, args []string, writer io.Writer) (Command, error) {
	flags := flag.NewFlagSet("tasks", flag.ContinueOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(writer, "Usage: cosmo steps <task> [--server=<name>]")
	}

	if len(args) == 0 {
		flags.Usage()
		return nil, ErrNoTask
	}

	if err := flags.Parse(args[1:]); err != nil {
		return nil, ErrFlagParse
	}

	taskName := args[0]

	task, _, err := config.Servers.FindTask(taskName, *server)
	if err != nil {
		fmt.Fprintln(writer, err)
		return nil, ErrFindTask
	}

	return &CommandSteps{
		TaskName: taskName,
		Task:     task,
		Writer:   writer,
	}, nil
}
