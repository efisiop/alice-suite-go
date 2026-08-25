package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SetPrimaryConsultant transfers a reader's active assignment for one book.
func SetPrimaryConsultant(readerID, bookID, consultantID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(Rebind(`UPDATE consultant_assignments SET active = 0 WHERE user_id = ? AND book_id = ?`), readerID, bookID); err != nil {
		return fmt.Errorf("deactivate primary consultant: %w", err)
	}
	_, err = tx.Exec(Rebind(`INSERT INTO consultant_assignments (id, consultant_id, user_id, book_id, active, created_at)
		VALUES (?, ?, ?, ?, 1, ?)
		ON CONFLICT(consultant_id, user_id, book_id) DO UPDATE SET active = 1`), uuid.NewString(), consultantID, readerID, bookID, time.Now())
	if err != nil {
		return fmt.Errorf("set primary consultant: %w", err)
	}
	return tx.Commit()
}
