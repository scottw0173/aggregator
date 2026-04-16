package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scottw0173/aggregator/internal/config"
	"github.com/scottw0173/aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type handlerFunc func(*state, command) error

type commands struct {
	handlers map[string]handlerFunc
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires exactly one argument")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user: %s is not registered with db", username)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("username successfully set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing name")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	newUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	s.cfg.SetUser(cmd.args[0])

	fmt.Printf("user successfully created: %s\n", cmd.args[0])
	fmt.Println(newUser)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unregistered command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, exists := c.handlers[name]; exists {
		return fmt.Errorf("duplicate command")
	}
	c.handlers[name] = f
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("database successfully reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	usersList, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range usersList {
		if user == s.cfg.UserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Print(rss)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("need two arguments: name and url")
	}
	userID, err := s.db.GetUserID(context.Background(), s.cfg.UserName)
	if err != nil {
		return err
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userID,
	}

	if _, err := s.db.CreateFeed(context.Background(), newFeed); err != nil {
		return err
	}
	newFollowRow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    newFeed.ID,
	}
	newFeedFollowRow, err := s.db.CreateFeedFollow(context.Background(), newFollowRow)
	if err != nil {
		return err
	}
	fmt.Println(newFeedFollowRow.FeedName)
	fmt.Println(newFeedFollowRow.UserName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsInfo(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("FeedName: %s\n", feed.Name)
		fmt.Printf("FeedURL: %s\n", feed.Url)
		fmt.Printf("UserName: %s\n", feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("command 'follow' need one argument: URL")
	}

	feedID, err := s.db.GetFeedID(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	userID, err := s.db.GetUserID(context.Background(), s.cfg.UserName)
	if err != nil {
		return err
	}

	newFollowRow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feedID,
	}
	newFeedFollowRow, err := s.db.CreateFeedFollow(context.Background(), newFollowRow)
	if err != nil {
		return err
	}

	fmt.Println(newFeedFollowRow.FeedName)
	fmt.Println(newFeedFollowRow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	var userName string
	if len(cmd.args) == 0 {
		userName = s.cfg.UserName
	} else {
		userName = cmd.args[0]
	}

	userID, err := s.db.GetUserID(context.Background(), userName)
	if err != nil {
		return err
	}

	followingList, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s\nis following\n", userName)
	for _, item := range followingList {
		fmt.Printf("*%s\n", item)
	}
	return nil
}
