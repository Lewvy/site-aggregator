package cli

import (
	"fmt"

	"github.com/Lewvy/aggregator/internal/config"
	"github.com/Lewvy/aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}
type Commands struct {
	CliCommand map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if f, ok := c.CliCommand[cmd.Name]; ok {
		return f(s, cmd)
	}
	return fmt.Errorf("Not a registered command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CliCommand[name] = f
}
