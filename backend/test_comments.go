package main

import (
	"fmt"
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
)

func main() {
	database.InitDB()
	db := database.GetDB()
	
	commentRepo := repository.NewCommentRepository(db)
	
	comments, total, err := commentRepo.ListByPostID(1, 1, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Comments: %+v\n", comments)
	for _, c := range comments {
		fmt.Printf("  - ID: %d, Content: %s, Author: %s\n", c.ID, c.Content, c.AuthorName)
	}
}
