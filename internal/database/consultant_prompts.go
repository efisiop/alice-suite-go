package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/models"
	"github.com/google/uuid"
)

// InsertConsultantPrompt creates a new consultant prompt for a reader at a page/section
func InsertConsultantPrompt(p *models.ConsultantPrompt) error {
	p.ID = uuid.New().String()
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", now)
	p.UpdatedAt = p.CreatedAt

	var sectionNum interface{}
	if p.SectionNumber != nil {
		sectionNum = *p.SectionNumber
	}
	query := `INSERT INTO consultant_prompts (id, user_id, book_id, page_number, section_number, prompt_text, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, p.ID, p.UserID, p.BookID, p.PageNumber, sectionNum, p.PromptText, now, now)
	return err
}

// GetConsultantPromptsForReader returns all prompts for a reader (consultant inspector list), including dismissed_at and accepted_at for feedback
func GetConsultantPromptsForReader(userID string) ([]*models.ConsultantPrompt, error) {
	query := `SELECT id, user_id, book_id, page_number, section_number, prompt_text, created_at, updated_at, dismissed_at, accepted_at
	          FROM consultant_prompts WHERE user_id = ? ORDER BY page_number, COALESCE(section_number, 0), created_at DESC`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("consultant prompts query: %w", err)
	}
	defer rows.Close()
	var list []*models.ConsultantPrompt
	for rows.Next() {
		var p models.ConsultantPrompt
		var sectionNum sql.NullInt64
		var createdAt, updatedAt string
		var dismissedAt, acceptedAt sql.NullString
		err := rows.Scan(&p.ID, &p.UserID, &p.BookID, &p.PageNumber, &sectionNum, &p.PromptText, &createdAt, &updatedAt, &dismissedAt, &acceptedAt)
		if err != nil {
			return nil, err
		}
		if sectionNum.Valid {
			n := int(sectionNum.Int64)
			p.SectionNumber = &n
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		if dismissedAt.Valid && dismissedAt.String != "" {
			if t, e := time.Parse("2006-01-02 15:04:05", dismissedAt.String); e == nil {
				p.DismissedAt = &t
			}
		}
		if acceptedAt.Valid && acceptedAt.String != "" {
			if t, e := time.Parse("2006-01-02 15:04:05", acceptedAt.String); e == nil {
				p.AcceptedAt = &t
			}
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

// GetConsultantPromptsForReaderAtPage returns prompts for a reader at a specific page (and optionally section)
// Only returns prompts that have not been dismissed (dismissed_at IS NULL)
func GetConsultantPromptsForReaderAtPage(userID, bookID string, pageNumber int, sectionNumber *int) ([]*models.ConsultantPrompt, error) {
	var query string
	var args []interface{}
	andClosed := ` AND (dismissed_at IS NULL OR dismissed_at = '') AND (accepted_at IS NULL OR accepted_at = '')`
	if sectionNumber == nil {
		query = `SELECT id, user_id, book_id, page_number, section_number, prompt_text, created_at, updated_at
	          FROM consultant_prompts
	          WHERE user_id = ? AND book_id = ? AND page_number = ? AND section_number IS NULL` + andClosed + `
	          ORDER BY created_at DESC`
		args = []interface{}{userID, bookID, pageNumber}
	} else {
		query = `SELECT id, user_id, book_id, page_number, section_number, prompt_text, created_at, updated_at
	          FROM consultant_prompts
	          WHERE user_id = ? AND book_id = ? AND page_number = ?
	          AND (section_number IS NULL OR section_number = ?)` + andClosed + `
	          ORDER BY created_at DESC`
		args = []interface{}{userID, bookID, pageNumber, *sectionNumber}
	}
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("consultant prompts at page: %w", err)
	}
	defer rows.Close()
	var list []*models.ConsultantPrompt
	for rows.Next() {
		var p models.ConsultantPrompt
		var sectionNum sql.NullInt64
		var createdAt, updatedAt string
		err := rows.Scan(&p.ID, &p.UserID, &p.BookID, &p.PageNumber, &sectionNum, &p.PromptText, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		if sectionNum.Valid {
			n := int(sectionNum.Int64)
			p.SectionNumber = &n
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		list = append(list, &p)
	}
	return list, rows.Err()
}

// DeleteConsultantPrompt deletes a consultant prompt by ID
func DeleteConsultantPrompt(id string) error {
	_, err := DB.Exec("DELETE FROM consultant_prompts WHERE id = ?", id)
	return err
}

// DismissConsultantPrompt records that the reader dismissed this prompt (feedback to consultant)
func DismissConsultantPrompt(promptID, userID string) error {
	_, err := DB.Exec(`UPDATE consultant_prompts SET dismissed_at = datetime('now') WHERE id = ? AND user_id = ?`, promptID, userID)
	return err
}

// AcceptConsultantPrompt records that the reader clicked Open AI Help and interacted (feedback to consultant)
func AcceptConsultantPrompt(promptID, userID string) error {
	_, err := DB.Exec(`UPDATE consultant_prompts SET accepted_at = datetime('now') WHERE id = ? AND user_id = ?`, promptID, userID)
	return err
}

// ReTriggerConsultantPrompt clears dismissed_at and accepted_at so the prompt shows again to the reader (new cycle)
func ReTriggerConsultantPrompt(promptID string) error {
	_, err := DB.Exec(`UPDATE consultant_prompts SET dismissed_at = NULL, accepted_at = NULL WHERE id = ?`, promptID)
	return err
}
