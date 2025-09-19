package ffmpeg

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strconv"

	"github.com/Nadim147c/material/color"
)

// GetPixels decodes media using ffmpeg and returns slices of pixels for given
// maxFrames numbers
func GetPixels(ctx context.Context, path string, maxFrames int) ([]color.ARGB, error) {
	var pixels []color.ARGB

	meta, err := CheckMediaType(path)
	if err != nil {
		return pixels, fmt.Errorf("Failed to get media metadata: %w", err)
	}

	if meta.Type != "image" && meta.Type != "video" {
		return pixels, fmt.Errorf("Invalid media type: %s", meta.Type)
	}

	if meta.Type == "image" {
		ffmpeg := exec.CommandContext(ctx, "ffmpeg",
			"-i", path,
			"-vframes", "1",
			"-f", "rawvideo",
			"-pix_fmt", "rgb24",
			"-")

		out, err := ffmpeg.Output()
		if err != nil {
			return pixels, fmt.Errorf("failed to get pixels from ffmpeg command: %w", err)
		}

		totalBytes := len(out)

		pixels := make([]color.ARGB, 0, totalBytes/3)
		for i := 0; i+2 < totalBytes; i += 3 {
			c := color.ARGBFromRGB(out[i], out[i+1], out[i+2])
			pixels = append(pixels, c)
		}

		return pixels, nil
	}

	fps := float64(maxFrames) / meta.Duration

	if math.Floor(meta.Duration) < float64(maxFrames) {
		fps = 1
	}

	ffmpeg := exec.CommandContext(ctx, "ffmpeg",
		"-i", path,
		"-vf", fmt.Sprintf("fps=%.8f", fps),
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-")

	out, err := ffmpeg.Output()
	if err != nil {
		return pixels, fmt.Errorf("failed to process image: %w", err)
	}

	totalBytes := len(out)

	pixels = make([]color.ARGB, 0, totalBytes/3)
	for i := 0; i+2 < totalBytes; i += 3 {
		c := color.ARGBFromRGB(out[i], out[i+1], out[i+2])
		pixels = append(pixels, c)
	}

	return pixels, nil
}

// MediaInfo holds the file type and duration
type MediaInfo struct {
	Type     string  // "image", "video", or "unknown"
	Duration float64 // Duration in seconds
}

// CheckMediaType runs ffprobe to determine if the file is an image or video and
// returns its duration
func CheckMediaType(filePath string) (MediaInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)
	output, err := cmd.Output()
	if err != nil {
		return MediaInfo{Type: "unknown", Duration: 0}, err
	}

	// Define structures to parse ffprobe JSON output
	type Stream struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
		NbFrames  string `json:"nb_frames"`
	}
	type Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
	}
	type FFProbeOutput struct {
		Streams []Stream `json:"streams"`
		Format  Format   `json:"format"`
	}

	// Parse JSON output
	var data FFProbeOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return MediaInfo{Type: "unknown", Duration: 0}, err
	}

	// Extract video and audio streams
	var videoStreams, audioStreams []Stream
	for _, stream := range data.Streams {
		switch stream.CodecType {
		case "video":
			videoStreams = append(videoStreams, stream)
		case "audio":
			audioStreams = append(audioStreams, stream)
		}
	}

	// Parse duration
	duration := 0.0
	if data.Format.Duration != "" {
		if d, err := strconv.ParseFloat(data.Format.Duration, 64); err == nil {
			duration = d
		}
	}

	// Determine file type
	mediaType := "unknown"
	if len(videoStreams) == 0 {
		return MediaInfo{Type: "unknown", Duration: 0}, nil
	}

	// Check for image-specific formats or single-frame video
	if data.Format.FormatName == "image2" || data.Format.FormatName == "png" || data.Format.FormatName == "jpeg" {
		mediaType = "image"
	} else if len(videoStreams) == 1 && len(audioStreams) == 0 {
		nbFrames := 0
		if videoStreams[0].NbFrames != "" {
			if n, err := strconv.Atoi(videoStreams[0].NbFrames); err == nil {
				nbFrames = n
			}
		}
		if nbFrames <= 1 && videoStreams[0].CodecName != "gif" {
			mediaType = "image"
		} else {
			mediaType = "video"
		}
	} else if len(videoStreams) > 0 && (duration > 0 || len(audioStreams) > 0) {
		mediaType = "video"
	}

	return MediaInfo{Type: mediaType, Duration: duration}, nil
}
