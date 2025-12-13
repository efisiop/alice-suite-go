package main

import (
	"fmt"
	"log"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

func main() {
	// Initialize database
	if err := database.InitDB("data/alice-suite.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Test GetPageByNumber
	page, err := database.GetPageByNumber("alice-in-wonderland", 1)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if page == nil {
		log.Fatal("Page is nil!")
	}

	fmt.Printf("✅ Page loaded successfully!\n")
	fmt.Printf("Page ID: %s\n", page.ID)
	fmt.Printf("Page Number: %d\n", page.PageNumber)
	fmt.Printf("Chapter Title: %v\n", page.ChapterTitle)
	fmt.Printf("Word Count: %d\n", page.WordCount)
	fmt.Printf("Sections: %d\n", len(page.Sections))
	
	for i, section := range page.Sections {
		fmt.Printf("  Section %d: %d words\n", section.SectionNumber, section.WordCount)
		if i >= 2 {
			fmt.Printf("  ... and %d more sections\n", len(page.Sections)-3)
			break
		}
	}
}


