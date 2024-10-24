package converter

import (
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
)

type VideoTask struct {
	VideoId int    `json:"video_id"`
	Path    string `json:"path"`
}

type VideoConvert struct{}

func NewVideoConvert() *VideoConvert {
	return &VideoConvert{}
}

func (c *VideoConvert) Handle(msg []byte) error {
	var task VideoTask
	err := json.Unmarshal(msg, &task)
	if err != nil {
		c.logError(task, "failed to unmarshal task", err)
		return err
	}
	err = c.processVideo(&task)
	if err != nil {
		c.logError(task, "failed to process video", err)
		return err
	}
	return nil
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
	// TODO: register error on database
}
