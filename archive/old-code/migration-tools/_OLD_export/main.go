package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "data/alice-suite.db"

func main() {
	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get book
	var book struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		TotalPages  int    `json:"total_pages"`
	}
	err = db.QueryRow("SELECT id, title, author, description, total_pages FROM books WHERE id = 'alice-in-wonderland'").Scan(
		&book.ID, &book.Title, &book.Author, &book.Description, &book.TotalPages)
	if err != nil {
		fmt.Printf("Error reading book: %v\n", err)
		os.Exit(1)
	}

	// Get pages (new page-based structure)
	pageRows, err := db.Query(`
		SELECT id, page_number, chapter_id, chapter_title, content, word_count
		FROM pages 
		ORDER BY page_number
	`)
	if err != nil {
		fmt.Printf("Error reading pages: %v\n", err)
		os.Exit(1)
	}
	defer pageRows.Close()

	// Define Section type first
	type Section struct {
		ID           string `json:"id"`
		PageID       string `json:"page_id"`
		PageNumber   int    `json:"page_number"`
		SectionNumber int  `json:"section_number"`
		Content      string `json:"content"`
		WordCount    int    `json:"word_count"`
	}

	type Page struct {
		ID          string  `json:"id"`
		PageNumber  int     `json:"page_number"`
		ChapterID   *string `json:"chapter_id"`
		ChapterTitle *string `json:"chapter_title"`
		Content     string  `json:"content"`
		WordCount   int     `json:"word_count"`
		Sections    []Section `json:"sections"`
	}

	pages := []Page{}
	for pageRows.Next() {
		var p Page
		var chapterID, chapterTitle sql.NullString
		pageRows.Scan(&p.ID, &p.PageNumber, &chapterID, &chapterTitle, &p.Content, &p.WordCount)
		if chapterID.Valid {
			p.ChapterID = &chapterID.String
		}
		if chapterTitle.Valid {
			p.ChapterTitle = &chapterTitle.String
		}
		pages = append(pages, p)
	}

	// Get sections for each page

	sectionRows, err := db.Query(`
		SELECT id, page_id, page_number, section_number, content, word_count
		FROM sections 
		ORDER BY page_number, section_number
	`)
	if err != nil {
		fmt.Printf("Error reading sections: %v\n", err)
		os.Exit(1)
	}
	defer sectionRows.Close()

	allSections := []Section{}
	for sectionRows.Next() {
		var sec Section
		sectionRows.Scan(&sec.ID, &sec.PageID, &sec.PageNumber, &sec.SectionNumber, &sec.Content, &sec.WordCount)
		allSections = append(allSections, sec)
	}

	// Attach sections to their pages
	for i := range pages {
		for _, sec := range allSections {
			if sec.PageID == pages[i].ID {
				pages[i].Sections = append(pages[i].Sections, sec)
			}
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"book":  book,
		"pages": pages,
	})
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// HTML template
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Alice Book Viewer - Standalone</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f5f5;
            padding: 20px;
            line-height: 1.6;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            border-bottom: 3px solid #667eea;
            padding-bottom: 10px;
        }
        .subtitle { color: #666; margin-bottom: 30px; }
        .sidebar {
            float: left;
            width: 250px;
            margin-right: 30px;
            background: #f9f9f9;
            padding: 20px;
            border-radius: 8px;
            max-height: 80vh;
            overflow-y: auto;
        }
        .content { margin-left: 280px; }
        .chapter-list { list-style: none; }
        .chapter-item {
            padding: 10px;
            margin-bottom: 5px;
            cursor: pointer;
            border-radius: 4px;
            transition: background 0.2s;
        }
        .chapter-item:hover { background: #e0e0e0; }
        .chapter-item.active { background: #667eea; color: white; }
        .section-list { margin-top: 10px; padding-left: 20px; }
        .section-item {
            padding: 8px;
            margin-bottom: 3px;
            cursor: pointer;
            border-radius: 4px;
            font-size: 0.9em;
            transition: background 0.2s;
        }
        .section-item:hover { background: #e0e0e0; }
        .section-item.active { background: #764ba2; color: white; }
        .book-content {
            background: #fff;
            padding: 30px;
            border-radius: 8px;
            border: 1px solid #ddd;
            min-height: 400px;
        }
        .section-header {
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #667eea;
        }
        .section-title {
            font-size: 1.5em;
            color: #333;
            margin-bottom: 5px;
        }
        .section-meta { color: #666; font-size: 0.9em; }
        .section-text {
            margin-top: 20px;
            font-size: 1.1em;
            line-height: 1.8;
            color: #333;
            white-space: pre-wrap;
        }
        .stats {
            background: #e8f4f8;
            padding: 15px;
            border-radius: 4px;
            margin-bottom: 20px;
            font-size: 0.9em;
        }
        .stats strong { color: #667eea; }
        @media (max-width: 768px) {
            .sidebar { float: none; width: 100%; margin-right: 0; margin-bottom: 20px; }
            .content { margin-left: 0; }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📖 Alice Book Viewer</h1>
        <p class="subtitle">Standalone Viewer - All Data Embedded</p>
        
        <div class="stats" id="stats"></div>

        <div class="sidebar">
            <h3>Pages</h3>
            <ul class="chapter-list" id="pageList"></ul>
        </div>

        <div class="content">
            <div class="book-content" id="bookContent">
                <div style="text-align: center; padding: 40px; color: #666;">Select a page to view content</div>
            </div>
        </div>
    </div>

    <script>
        const data = {{.JSONData}};

        // Update stats
        var totalSections = 0;
        data.pages.forEach(function(page) { totalSections += page.sections.length; });
        document.getElementById('stats').innerHTML = 
            '<strong>Book:</strong> ' + data.book.title + '<br>' +
            '<strong>Author:</strong> ' + data.book.author + '<br>' +
            '<strong>Pages:</strong> ' + data.pages.length + '<br>' +
            '<strong>Sections:</strong> ' + totalSections + '<br>' +
            '<strong>Total Pages:</strong> ' + data.book.total_pages;

        // Render pages
        const pageList = document.getElementById('pageList');
        pageList.innerHTML = data.pages.map(function(page) {
            var pageTitle = 'Page ' + page.page_number;
            if (page.chapter_title) {
                pageTitle = pageTitle + ' - ' + page.chapter_title;
            }
            return '<li class="chapter-item" data-page-id="' + page.id + '" data-page-number="' + page.page_number + '">' +
                   pageTitle + ' (' + page.sections.length + ' sections, ' + page.word_count + ' words)' +
                   '<ul class="section-list" id="sections-' + page.id + '" style="display:none;"></ul>' +
                   '</li>';
        }).join('');

        // Add page click handlers
        document.querySelectorAll('.chapter-item').forEach(function(item) {
            item.addEventListener('click', function(e) {
                e.stopPropagation();
                const pageId = item.dataset.pageId;
                const pageNumber = item.dataset.pageNumber;
                
                // Toggle active
                document.querySelectorAll('.chapter-item').forEach(function(i) { i.classList.remove('active'); });
                item.classList.add('active');

                // Get page data
                const page = data.pages.find(function(p) { return p.id === pageId; });
                if (!page) return;
                
                // Render sections for this page
                const sectionList = document.getElementById('sections-' + pageId);
                sectionList.innerHTML = page.sections.map(function(section) {
                    return '<li class="section-item" data-section-id="' + section.id + '">' +
                           'Section ' + section.section_number + ' (' + section.word_count + ' words)' +
                           '</li>';
                }).join('');

                // Toggle section list
                sectionList.style.display = sectionList.style.display === 'none' ? 'block' : 'none';

                // Add section click handlers
                sectionList.querySelectorAll('.section-item').forEach(function(sectionItem) {
                    sectionItem.addEventListener('click', function(e) {
                        e.stopPropagation();
                        const sectionId = sectionItem.dataset.sectionId;
                        
                        // Toggle active
                        document.querySelectorAll('.section-item').forEach(function(i) { i.classList.remove('active'); });
                        sectionItem.classList.add('active');

                        // Find and display section
                        const section = page.sections.find(function(s) { return s.id === sectionId; });
                        if (section) {
                            const contentDiv = document.getElementById('bookContent');
                            var header = '<div class="section-header">' +
                                '<div class="section-title">Page ' + page.page_number + ', Section ' + section.section_number + '</div>' +
                                '<div class="section-meta">' + section.word_count + ' words';
                            if (page.chapter_title) {
                                header += ' | ' + page.chapter_title;
                            }
                            header += '</div></div>';
                            contentDiv.innerHTML = header +
                                '<div class="section-text">' + escapeHtml(section.content) + '</div>';
                        }
                    });
                });
            });
        });

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    </script>
</body>
</html>`

	// Parse template
	tmpl, err := template.New("viewer").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outputPath := filepath.Join("static", "viewer-standalone.html")
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, map[string]interface{}{
		"JSONData": template.JS(string(jsonData)),
	})
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Standalone viewer created: %s\n", outputPath)
	fmt.Printf("📖 Book: %s\n", book.Title)
	fmt.Printf("📄 Pages: %d\n", len(pages))
	var totalSections int
	for _, p := range pages {
		totalSections += len(p.Sections)
	}
	fmt.Printf("📝 Sections: %d\n", totalSections)
	fmt.Printf("\n💡 Open this file directly in your browser:\n")
	
	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("   file://%s\n", absPath)
	
	// Try to open it
	exec.Command("open", absPath).Run()
}
