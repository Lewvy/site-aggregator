package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lewvy/aggregator/internal/database"
	"github.com/Lewvy/aggregator/rss"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid arguments. want=1, got=%d", len(cmd.Args))
	}
	name := cmd.Args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.Db.GetUser(ctx, name)
	if err != nil {
		return fmt.Errorf("User doesn't exist: %q", err)
	}
	if err = s.Cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user: %q", err)
	}
	log.Println("User switched successfully")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid arguments. want=1, got=%d", len(cmd.Args))
	}
	name := cmd.Args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	uid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	new_user := &database.CreateUserParams{
		Name:      name,
		ID:        uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user, err := s.Db.CreateUser(ctx, *new_user)
	if err != nil {
		return err
	}
	log.Println("User created successfully", user.ID, user.Name)
	if err = s.Cfg.SetUser(new_user.Name); err != nil {
		return err
	}
	return nil
}
func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil

}

func HandlerList(s *State, cmd Command) error {
	users, err := s.Db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error fetching registered users: %q\n", err)
	}
	loggedInUser := s.Cfg.USER
	for _, user := range users {
		if user == loggedInUser {
			fmt.Println("*", user, "(current)")
		} else {
			fmt.Println("*", user)

		}
	}
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if err := s.Db.DropRows(context.Background()); err != nil {
		return fmt.Errorf("Error resetting: %q", err)
	}
	log.Print("Table reset performed successfully")
	return nil

}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Invalid number of args: want 2, got=%d", len(cmd.Args))
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	curr_user, err := s.Db.GetUser(context.Background(), s.Cfg.USER)
	if err != nil {
		return fmt.Errorf("Error fetching user from the db: %w", err)
	}
	AddFeedParams := &database.AddFeedAndFollowParams{}
	AddFeedParams.ID = uuid.New()
	AddFeedParams.Url = url
	AddFeedParams.Name = name
	dbFeed, err := s.Db.AddFeedAndFollow(context.Background(), *AddFeedParams)
	if err != nil {
		return fmt.Errorf("Error inserting into feed: %w", err)
	}
	fmt.Println(dbFeed)
	return nil
}

func ListFeeds(s *State, cmd Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rows, err := s.Db.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rows")
	}
	for _, row := range rows {
		fmt.Println(row.Name, row.Url, row.Name_2)
	}
	return nil
}

func ListFollowing(s *State, cmd Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := s.Db.GetUser(ctx, s.Cfg.USER)
	if err != nil {
		return fmt.Errorf("Error gettine user: %q", err)
	}
	followed_feeds, err := s.Db.GetFeedNamesUserIsFollowing(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Println(followed_feeds)
	return nil
}
