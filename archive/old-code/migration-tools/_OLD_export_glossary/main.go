package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "data/alice-suite.db"

type GlossaryTerm struct {
	Term       string
	Definition string
	Example    string
}

type Section struct {
	ID           string
	PageNumber   int
	SectionNumber int
	Content      string
	WordCount    int
	GlossaryTerms []GlossaryTerm
}

type Page struct {
	ID           string
	PageNumber   int
	ChapterID    *string
	ChapterTitle *string
	Content      string
	WordCount    int
	Sections     []Section
}

type BookData struct {
	Book struct {
		ID          string
		Title       string
		Author      string
		Description string
		TotalPages  int
	}
	Pages []Page
}

func main() {
	fmt.Println("📚 Exporting Book Content with Glossary Links")
	fmt.Println("🔗 Creating viewer with glossary term highlighting")
	fmt.Println("")

	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		fmt.Printf("❌ Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get book info
	var bookData BookData
	err = db.QueryRow(`
		SELECT id, title, author, description, total_pages
		FROM books
		WHERE id = 'alice-in-wonderland'
	`).Scan(&bookData.Book.ID, &bookData.Book.Title, &bookData.Book.Author,
		&bookData.Book.Description, &bookData.Book.TotalPages)
	if err != nil {
		fmt.Printf("❌ Failed to get book info: %v\n", err)
		os.Exit(1)
	}

	// Get pages with sections
	pagesRows, err := db.Query(`
		SELECT id, page_number, chapter_id, chapter_title, content, word_count
		FROM pages
		ORDER BY page_number
	`)
	if err != nil {
		fmt.Printf("❌ Failed to query pages: %v\n", err)
		os.Exit(1)
	}
	defer pagesRows.Close()

	for pagesRows.Next() {
		var page Page
		var chapterID, chapterTitle sql.NullString
		pagesRows.Scan(&page.ID, &page.PageNumber, &chapterID, &chapterTitle,
			&page.Content, &page.WordCount)
		
		if chapterID.Valid {
			page.ChapterID = &chapterID.String
		}
		if chapterTitle.Valid {
			page.ChapterTitle = &chapterTitle.String
		}

		// Get sections for this page
		sectionsRows, err := db.Query(`
			SELECT id, page_number, section_number, content, word_count
			FROM sections
			WHERE page_id = ?
			ORDER BY section_number
		`, page.ID)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to query sections for page %d: %v\n", page.PageNumber, err)
			continue
		}

		for sectionsRows.Next() {
			var section Section
			sectionsRows.Scan(&section.ID, &section.PageNumber, &section.SectionNumber,
				&section.Content, &section.WordCount)

			// Get glossary terms for this section
			glossaryRows, err := db.Query(`
				SELECT DISTINCT g.term, g.definition, COALESCE(g.example, '') as example
				FROM glossary_section_links gs
				JOIN alice_glossary g ON gs.glossary_id = g.id
				WHERE gs.section_id = ?
				ORDER BY g.term
			`, section.ID)
			if err == nil {
				for glossaryRows.Next() {
					var term GlossaryTerm
					glossaryRows.Scan(&term.Term, &term.Definition, &term.Example)
					section.GlossaryTerms = append(section.GlossaryTerms, term)
				}
				glossaryRows.Close()
			}

			page.Sections = append(page.Sections, section)
		}
		sectionsRows.Close()

		bookData.Pages = append(bookData.Pages, page)
	}

	fmt.Printf("✅ Loaded %d pages with glossary links\n", len(bookData.Pages))

	// Generate HTML
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Book.Title}} - Glossary Linked Viewer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Georgia', 'Times New Roman', serif;
            line-height: 1.8;
            color: #333;
            background: #fafafa;
            padding: 20px;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            border-bottom: 3px solid #3498db;
            padding-bottom: 10px;
            margin-bottom: 30px;
        }
        .book-info {
            background: #ecf0f1;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 30px;
        }
        .page {
            margin-bottom: 50px;
            page-break-after: always;
        }
        .page-header {
            background: #34495e;
            color: white;
            padding: 15px;
            border-radius: 5px 5px 0 0;
            margin-bottom: 0;
        }
        .chapter-title {
            font-size: 1.5em;
            color: #e74c3c;
            margin-bottom: 20px;
            font-weight: bold;
        }
        .section {
            margin-bottom: 30px;
            padding: 20px;
            background: #fff;
            border-left: 4px solid #3498db;
        }
        .section-header {
            color: #7f8c8d;
            font-size: 0.9em;
            margin-bottom: 15px;
            font-weight: bold;
        }
        .section-content {
            font-size: 1.1em;
            text-align: justify;
            margin-bottom: 15px;
        }
        .glossary-term {
            background: #fff3cd;
            color: #856404;
            padding: 2px 6px;
            border-radius: 3px;
            cursor: pointer;
            border-bottom: 2px dotted #ffc107;
            position: relative;
            display: inline-block;
        }
        .glossary-term:hover {
            background: #ffc107;
            color: #000;
        }
        .glossary-tooltip {
            display: none;
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            background: #2c3e50;
            color: white;
            padding: 10px;
            border-radius: 5px;
            width: 300px;
            font-size: 0.9em;
            z-index: 1000;
            margin-bottom: 5px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
        .glossary-term:hover .glossary-tooltip {
            display: block;
        }
        .glossary-tooltip strong {
            color: #ffc107;
            display: block;
            margin-bottom: 5px;
        }
        .glossary-list {
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid #ecf0f1;
        }
        .glossary-list-title {
            font-weight: bold;
            color: #7f8c8d;
            margin-bottom: 10px;
            font-size: 0.9em;
        }
        .glossary-item {
            margin-bottom: 10px;
            padding: 10px;
            background: #f8f9fa;
            border-radius: 3px;
        }
        .glossary-item-term {
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 5px;
        }
        .glossary-item-def {
            color: #555;
            font-size: 0.95em;
        }
        .stats {
            background: #e8f5e9;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 30px;
            text-align: center;
        }
        .nav {
            position: fixed;
            top: 20px;
            right: 20px;
            background: white;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.2);
        }
        .nav a {
            display: block;
            padding: 5px 10px;
            color: #3498db;
            text-decoration: none;
        }
        .nav a:hover {
            background: #ecf0f1;
        }
    </style>
</head>
<body>
    <div class="nav">
        <strong>Navigation</strong>
        {{range .Pages}}
        <a href="#page-{{.PageNumber}}">Page {{.PageNumber}}</a>
        {{end}}
    </div>

    <div class="container">
        <h1>{{.Book.Title}}</h1>
        <div class="book-info">
            <strong>Author:</strong> {{.Book.Author}}<br>
            <strong>Description:</strong> {{.Book.Description}}<br>
            <strong>Total Pages:</strong> {{.Book.TotalPages}}
        </div>

        <div class="stats">
            <strong>📊 Statistics:</strong><br>
            Pages: {{len .Pages}} | 
            Total Glossary Links: <span id="total-links">0</span>
        </div>

        {{range .Pages}}
        <div class="page" id="page-{{.PageNumber}}">
            <div class="page-header">
                📄 Page {{.PageNumber}} 
                {{if .ChapterTitle}}
                | {{.ChapterTitle}}
                {{end}}
                ({{.WordCount}} words)
            </div>

            {{range .Sections}}
            <div class="section">
                <div class="section-header">
                    Section {{.SectionNumber}} ({{.WordCount}} words)
                </div>
                <div class="section-content" id="section-{{.ID}}">
                    {{.Content}}
                </div>

                {{if .GlossaryTerms}}
                <div class="glossary-list">
                    <div class="glossary-list-title">
                        📚 Glossary Terms in this Section ({{len .GlossaryTerms}}):
                    </div>
                    {{range .GlossaryTerms}}
                    <div class="glossary-item">
                        <div class="glossary-item-term">{{.Term}}</div>
                        <div class="glossary-item-def">{{.Definition}}</div>
                        {{if .Example}}
                        <div style="font-size: 0.85em; color: #7f8c8d; margin-top: 5px; font-style: italic;">
                            Example: {{.Example}}
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
        {{end}}
    </div>

    <script>
        // Highlight glossary terms in content
        const pages = {{.Pages}};
        let totalLinks = 0;

        pages.forEach(page => {
            page.Sections.forEach(section => {
                section.GlossaryTerms.forEach(term => {
                    totalLinks++;
                    const sectionEl = document.getElementById('section-' + section.ID);
                    if (sectionEl) {
                        const content = sectionEl.innerHTML;
                        // Create regex to match the term (case-insensitive, word boundaries)
                        const regex = new RegExp('\\b' + term.Term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + '\\b', 'gi');
                        sectionEl.innerHTML = content.replace(regex, (match) => {
                            return '<span class="glossary-term" title="' + 
                                term.Definition.replace(/"/g, '&quot;') + '">' + match + 
                                '<span class="glossary-tooltip"><strong>' + term.Term + '</strong>' + 
                                term.Definition + 
                                (term.Example ? '<br><em>' + term.Example + '</em>' : '') +
                                '</span></span>';
                        });
                    }
                });
            });
        });

        document.getElementById('total-links').textContent = totalLinks;
    </script>
</body>
</html>`

	tmpl, err := template.New("viewer").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("❌ Failed to parse template: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outputFile := "static/viewer-glossary.html"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("❌ Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, bookData)
	if err != nil {
		fmt.Printf("❌ Failed to execute template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Viewer created: %s\n", outputFile)
	fmt.Printf("📖 Open in browser: file://%s/%s\n", getCurrentDir(), outputFile)
	fmt.Println("")
	fmt.Println("💡 Features:")
	fmt.Println("   - All book content with pages and sections")
	fmt.Println("   - Glossary terms highlighted in yellow")
	fmt.Println("   - Hover over terms to see definitions")
	fmt.Println("   - Glossary list shown below each section")
}

func getCurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

