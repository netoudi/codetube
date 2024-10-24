package converter

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"
)

func IsProcessed(db *sql.DB, videoId int) bool {
	var isProcessed bool
	query := "SELECT EXISTS(SELECT 1 FROM processed_videos WHERE video_id = $1 AND status = 'success')"
	err := db.QueryRow(query, videoId).Scan(&isProcessed)
	if err != nil {
		slog.Error("Error checking if video is processed", slog.Int("video_id", videoId))
		return false
	}
	return isProcessed
}

func MarkAsProcessed(db *sql.DB, videoId int) error {
	query := "INSERT INTO processed_videos (video_id, status, processed_at) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, videoId, "success", time.Now())
	if err != nil {
		slog.Error("Error marking video as processed", slog.Int("video_id", videoId))
		return err
	}
	return nil
}

func RegisterError(db *sql.DB, errorData map[string]any) {
	serializedError, _ := json.Marshal(errorData)
	query := "INSERT INTO process_errors_log (error_details, created_at) VALUES ($1, $2)"
	_, err := db.Exec(query, string(serializedError), time.Now())
	if err != nil {
		slog.Error("Error storing error log in database", slog.String("error", string(err.Error())))
		return
	}
	slog.Error("Error log stored successfully", slog.String("error", string(serializedError)))
}
