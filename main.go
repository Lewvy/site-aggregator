package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/Lewvy/aggregator/internal/cli"
	"github.com/Lewvy/aggregator/internal/config"
	"github.com/Lewvy/aggregator/internal/database"
)

func main() {
	db_url := "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", db_url)
	dbQueries := database.New(db)
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %q", err)
	}

	progState := &cli.State{
		Cfg: &cfg,
		Db:  dbQueries,
	}

	cmds := cli.Commands{
		CliCommand: make(map[string]func(*cli.State, cli.Command) error),
	}

	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)
	cmds.Register("users", cli.HandlerList)
	cmds.Register("agg", cli.HandlerAgg)
	cmds.Register("addfeed", cli.HandlerAddFeed)
	cmds.Register("feeds", cli.ListFeeds)
	cmds.Register("follow", cli.ListFollowing)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [...args]")
		return
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err = cmds.Run(progState, cli.Command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
