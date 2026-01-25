package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}

type ColumnInfo struct {
	Name     string
	Type     string
	NotNull  bool
	Default  string
	Primary  bool
}

func main() {
	cfg := config.Load()

	fmt.Println("🔍 Database Structure Comparison Tool")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("")

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Printf("✅ Database connected: %s\n", cfg.DBPath)
	fmt.Println("")

	// Get all tables
	tables, err := getAllTables()
	if err != nil {
		log.Fatalf("Error getting tables: %v", err)
	}

	fmt.Printf("📊 Found %d tables:\n", len(tables))
	for _, table := range tables {
		fmt.Printf("   - %s\n", table.Name)
	}
	fmt.Println("")

	// Get structure for each table
	fmt.Println("📋 Table Structures:")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	for _, table := range tables {
		columns, err := getTableColumns(table.Name)
		if err != nil {
			log.Printf("Error getting columns for %s: %v", table.Name, err)
			continue
		}
		table.Columns = columns
		
		fmt.Printf("\n📄 Table: %s\n", table.Name)
		fmt.Printf("   Columns: %d\n", len(columns))
		for _, col := range columns {
			notNull := ""
			if col.NotNull {
				notNull = " NOT NULL"
			}
			defaultVal := ""
			if col.Default != "" {
				defaultVal = fmt.Sprintf(" DEFAULT %s", col.Default)
			}
			primary := ""
			if col.Primary {
				primary = " PRIMARY KEY"
			}
			fmt.Printf("   - %s: %s%s%s%s\n", col.Name, col.Type, notNull, defaultVal, primary)
		}
	}

	// Check specific tables for sections issue
	fmt.Println("")
	fmt.Println("🔍 Sections Table Analysis:")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	var sectionsExists bool
	var sectionsSQL string
	err = database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='sections'
	`).Scan(&sectionsSQL)
	
	if err == nil {
		sectionsExists = true
		fmt.Println("✅ Sections table exists")
		fmt.Printf("   Structure: %s\n", sectionsSQL)
		
		// Check for key fields
		hasPageNumber := strings.Contains(sectionsSQL, "page_number")
		hasSectionNumber := strings.Contains(sectionsSQL, "section_number")
		hasPageID := strings.Contains(sectionsSQL, "page_id")
		
		fmt.Println("")
		fmt.Println("   Key Fields:")
		fmt.Printf("   - page_number: %v\n", hasPageNumber)
		fmt.Printf("   - section_number: %v\n", hasSectionNumber)
		fmt.Printf("   - page_id: %v\n", hasPageID)
		
		if hasPageNumber && hasSectionNumber {
			fmt.Println("   ✅ New structure detected (page_number, section_number)")
		} else if hasPageID {
			fmt.Println("   ⚠️  Old structure detected (page_id only)")
		}
	} else {
		sectionsExists = false
		fmt.Println("❌ Sections table does NOT exist")
	}

	// Check pages table
	fmt.Println("")
	fmt.Println("🔍 Pages Table Analysis:")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	var pagesExists bool
	var pagesSQL string
	err = database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='pages'
	`).Scan(&pagesSQL)
	
	if err == nil {
		pagesExists = true
		fmt.Println("✅ Pages table exists")
	} else {
		pagesExists = false
		fmt.Println("❌ Pages table does NOT exist")
	}

	// Check data counts
	fmt.Println("")
	fmt.Println("📊 Data Counts:")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	if sectionsExists {
		var sectionCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&sectionCount)
		fmt.Printf("   Sections: %d\n", sectionCount)
		
		var page1Count int
		database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&page1Count)
		fmt.Printf("   Sections for page 1: %d\n", page1Count)
	}
	
	if pagesExists {
		var pageCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM pages").Scan(&pageCount)
		fmt.Printf("   Pages: %d\n", pageCount)
	}

	// Summary
	fmt.Println("")
	fmt.Println("📋 Summary:")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Printf("   Total tables: %d\n", len(tables))
	fmt.Printf("   Sections table exists: %v\n", sectionsExists)
	fmt.Printf("   Pages table exists: %v\n", pagesExists)
	
	if sectionsExists {
		var page1Count int
		database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&page1Count)
		if page1Count >= 5 {
			fmt.Println("   ✅ Page 1 has correct number of sections (5+)")
		} else {
			fmt.Printf("   ⚠️  Page 1 has only %d sections (expected 5+)\n", page1Count)
			fmt.Println("   💡 Run: ./bin/fix-render")
		}
	}
	
	fmt.Println("")
}

func getAllTables() ([]TableInfo, error) {
	rows, err := database.DB.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, TableInfo{Name: name})
	}
	return tables, rows.Err()
}

func getTableColumns(tableName string) ([]ColumnInfo, error) {
	rows, err := database.DB.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var defaultVal sql.NullString
		
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultVal, &pk); err != nil {
			return nil, err
		}
		
		columns = append(columns, ColumnInfo{
			Name:    name,
			Type:    colType,
			NotNull: notNull == 1,
			Default: defaultVal.String,
			Primary: pk == 1,
		})
	}
	
	// Sort by primary key first, then by name
	sort.Slice(columns, func(i, j int) bool {
		if columns[i].Primary != columns[j].Primary {
			return columns[i].Primary
		}
		return columns[i].Name < columns[j].Name
	})
	
	return columns, rows.Err()
}
