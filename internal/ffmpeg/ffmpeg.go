package ffmpeg

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Nadim147c/material/color"
	"github.com/gabriel-vasile/mimetype"
)

// GetPixels decodes media using ffmpeg and returns slices of pixels for given
// maxFrames numbers
func GetPixels(
	ctx context.Context,
	path string,
	maxFrames int,
) ([]color.ARGB, error) {
	var pixels []color.ARGB

	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to get media type: %w", err)
	}

	kind := mtype.String()

	if !strings.HasPrefix(kind, "video") && !strings.HasPrefix(kind, "image") {
		return pixels, fmt.Errorf("Invalid media type: %s", kind)
	}

	if strings.HasPrefix(kind, "image") {
		ffmpeg := exec.CommandContext(ctx, "ffmpeg",
			"-i", path,
			"-vframes", "1",
			"-f", "rawvideo",
			"-pix_fmt", "rgb24",
			"-")

		out, err := ffmpeg.Output()
		if err != nil {
			return pixels, fmt.Errorf(
				"failed to get pixels from ffmpeg command: %w",
				err,
			)
		}

		totalBytes := len(out)

		pixels := make([]color.ARGB, 0, totalBytes/3)
		for i := 0; i+2 < totalBytes; i += 3 {
			c := color.ARGBFromRGB(out[i], out[i+1], out[i+2])
			pixels = append(pixels, c)
		}

		return pixels, nil
	}

	duration, err := GetDuration(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get duration: %v", err)
	}

	fps := float64(maxFrames) / duration

	if math.Floor(duration) < float64(maxFrames) {
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

// GetDuration runs ffprobe to determine if the file is an image or video and
// returns its duration
func GetDuration(src string) (float64, error) {
	out, err := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		src,
	).Output()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}

// GeneratePreview generates preview thumnail for given media
func GeneratePreview(ctx context.Context, src, dst string) error {
	dur, err := GetDuration(src)
	if err != nil {
		return err
	}
	seek := fmt.Sprintf("%.2f", dur*0.10) // 10%

	cmd := exec.CommandContext(
		ctx, "ffmpeg",
		"-ss", seek,
		"-i", src,
		"-vframes", "1",
		"-y",
		dst,
	)
	return cmd.Run()
}
