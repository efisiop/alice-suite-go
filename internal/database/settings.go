package database

import "fmt"

// GetAdminSetting returns the value for the given key, or empty string if not found.
func GetAdminSetting(key string) (string, error) {
	var value string
	err := DB.QueryRow(Rebind("SELECT value FROM admin_settings WHERE key = ?"), key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("get admin setting %q: %w", key, err)
	}
	return value, nil
}

// SetAdminSetting sets the value for the given key.
func SetAdminSetting(key, value string) error {
	var query string
	if DriverName == "postgres" {
		query = `INSERT INTO admin_settings (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`
	} else {
		query = "INSERT OR REPLACE INTO admin_settings (key, value) VALUES (?, ?)"
	}
	_, err := DB.Exec(Rebind(query), key, value)
	if err != nil {
		return fmt.Errorf("set admin setting %q: %w", key, err)
	}
	return nil
}

// GetLoginEmailNotificationsEnabled returns true if login email notifications are enabled.
func GetLoginEmailNotificationsEnabled() (bool, error) {
	val, err := GetAdminSetting("login_email_notifications")
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// SetLoginEmailNotificationsEnabled enables or disables login email notifications.
func SetLoginEmailNotificationsEnabled(enabled bool) error {
	val := "0"
	if enabled {
		val = "1"
	}
	return SetAdminSetting("login_email_notifications", val)
}
