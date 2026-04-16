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

func middlewareLoggedIn(
	handler func(*state, command, database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("need one argument: time_between_reqs")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for range ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println("scrape error:", err)
		}
	}
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	if err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Println(feed.Channel.Title)
	for _, title := range feed.Channel.Item {
		fmt.Println(title.Title)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("need two arguments: name and url")
	}
	//userID, err := s.db.GetUserID(context.Background(), s.cfg.UserName)
	//if err != nil {
	//	return err
	//}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	if _, err := s.db.CreateFeed(context.Background(), newFeed); err != nil {
		return err
	}
	newFollowRow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("command 'follow' need one argument: URL")
	}

	feedID, err := s.db.GetFeedID(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	//userID, err := s.db.GetUserID(context.Background(), s.cfg.UserName)
	//if err != nil {
	//	return err
	//}

	newFollowRow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	followingList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s\nis following\n", user.Name)
	for _, item := range followingList {
		fmt.Printf("*%s\n", item)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("need argument: URL")
	}
	url := cmd.args[0]

	feedID, err := s.db.GetFeedID(context.Background(), url)
	if err != nil {
		return err
	}

	unfollowParams := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feedID,
	}

	if err = s.db.UnfollowFeed(context.Background(), unfollowParams); err != nil {
		return err
	}
	fmt.Printf("User: %s has unfollowed\n%s", user.Name, url)
	return nil
}
