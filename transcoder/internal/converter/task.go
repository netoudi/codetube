package converter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
	"transcoder/internal/rabbitmq"

	"github.com/streadway/amqp"
)

type VideoTask struct {
	VideoId int    `json:"video_id"`
	Path    string `json:"path"`
}

type VideoConvert struct {
	db             *sql.DB
	rabbitmqClient *rabbitmq.RabbitClient
}

func NewVideoConvert(db *sql.DB, rabbitmqClient *rabbitmq.RabbitClient) *VideoConvert {
	return &VideoConvert{
		db:             db,
		rabbitmqClient: rabbitmqClient,
	}
}

func (c *VideoConvert) Handle(d amqp.Delivery, conversionExchange, confirmationKey, confirmationQueue string) {
	var task VideoTask
	err := json.Unmarshal(d.Body, &task)
	if err != nil {
		c.logError(task, "failed to unmarshal task", err)
		return
	}
	if IsProcessed(c.db, task.VideoId) {
		slog.Info("Video already processed", slog.Int("video_id", task.VideoId))
		d.Ack(false)
		return
	}
	err = c.processVideo(&task)
	if err != nil {
		c.logError(task, "failed to process video", err)
		d.Ack(false)
		return
	}
	err = MarkAsProcessed(c.db, task.VideoId)
	if err != nil {
		c.logError(task, "failed to mark video as processed", err)
		return
	}
	d.Ack(false)
	slog.Info("Video marked as processed", slog.Int("video_id", task.VideoId))
	confirmationMessage := []byte(fmt.Sprintf(`{"video_id": %d, "path": "%s"}`, task.VideoId, task.Path))
	err = c.rabbitmqClient.PublishMessage(conversionExchange, confirmationKey, confirmationQueue, confirmationMessage)
	if err != nil {
		slog.Info("Failed to publish confirmation message", slog.String("error", err.Error()))
		c.logError(task, "failed to publish confirmation message", err)
	}
	slog.Info("Published confirmation message", slog.String("message", string(confirmationMessage)))
}

func (c *VideoConvert) processVideo(task *VideoTask) error {
	mergedFile := filepath.Join(task.Path, "merged.mp4")
	mpegDashPath := filepath.Join(task.Path, "mpeg-dash")

	slog.Info("Merging chunks", slog.String("path", task.Path))
	err := c.mergeChunks(task.Path, mergedFile)
	if err != nil {
		c.logError(*task, "failed to merge chunks", err)
		return err
	}

	slog.Info("Creating mpeg-dash dir", slog.String("path", mpegDashPath))
	err = os.MkdirAll(mpegDashPath, os.ModePerm)
	if err != nil {
		c.logError(*task, "failed to create mpeg-dash directory", err)
		return err
	}

	slog.Info("Converting video to mpeg-dash", slog.String("path", mpegDashPath))
	// ffmpegCmd := exec.Command("ffmpeg", "-i", mergedFile, "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", filepath.Join(mpegDashPath, "segment-%d.ts"))
	ffmpegCmd := exec.Command("ffmpeg", "-i", mergedFile, "-f", "dash", filepath.Join(mpegDashPath, "output.mpd"))
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		c.logError(*task, "failed to convert video to mpeg-dash, output: "+string(output), err)
		return err
	}
	slog.Info("Video converted to mpeg-dash", slog.String("path", mpegDashPath))

	slog.Info("Removing merged file", slog.String("path", mpegDashPath))
	err = os.Remove(mergedFile)
	if err != nil {
		c.logError(*task, "failed to remove merged file", err)
		return err
	}

	return nil
}

func (c *VideoConvert) extractNumber(filename string) int {
	re := regexp.MustCompile(`\d+`)
	munStr := re.FindString(filepath.Base(filename))
	num, err := strconv.Atoi(munStr)
	if err != nil {
		return -1
	}
	return num
}

func (c *VideoConvert) mergeChunks(inputDir string, outputFile string) error {
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk"))
	if err != nil {
		return fmt.Errorf("failed to find chunks: %v", err)
	}

	sort.Slice(chunks, func(i, j int) bool {
		return c.extractNumber(chunks[i]) < c.extractNumber(chunks[j])
	})

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	for _, chunk := range chunks {
		input, err := os.Open(chunk)
		if err != nil {
			return fmt.Errorf("failed to open chunk: %v", err)
		}
		_, err = output.ReadFrom(input)
		if err != nil {
			return fmt.Errorf("failed to write chunk %s to merged file: %v", chunk, err)
		}
		input.Close()
	}

	return nil
}

func (c *VideoConvert) logError(task VideoTask, message string, err error) {
	errorData := map[string]any{
		"video_id": task.VideoId,
		"error":    message,
		"details":  err.Error(),
		"time":     time.Now(),
	}
	serializedError, _ := json.Marshal(errorData)
	slog.Error("Processing error", slog.String("error_details", string(serializedError)))
	RegisterError(c.db, errorData)
}
