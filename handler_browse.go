package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/The-Legolas/aggregator_program/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s [optional] <post_limit_num>", cmd.name)
	} else if len(cmd.args) == 1 {
		specifiedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = specifiedLimit
	}

	postDB_args := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	userPosts, err := s.db.GetPostsForUser(context.Background(), postDB_args)
	if err != nil {
		return fmt.Errorf("Error getting %s's post: %w", user.Name, err)
	}

	fmt.Printf("Fetching %d post for user %s:\n\n", len(userPosts), user.Name)
	for _, post := range userPosts {
		printUserPost(post)
		fmt.Println("===================")
		fmt.Println()
	}

	return nil
}

func printUserPost(post database.GetPostsForUserRow) {
	fmt.Printf("* Post from:		%s\n", post.FeedName)
	fmt.Printf("* Post title:		%s\n", post.Title)
	fmt.Printf("* Published date:	%s\n", post.PublishedAt)
	fmt.Printf("* URL:			%s\n", post.Url)
	if post.Description.Valid {
		fmt.Printf("* Description: 		%v\n", post.Description.String)
	} else {
		fmt.Printf("* Description: 	N/A\n")
	}
}
