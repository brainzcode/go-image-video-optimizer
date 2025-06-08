package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ProcessVideo(inputFile, outputPath string, config Config, index int) error {
	// fmt.Printf("Processing video: %s\n", inputFile)
	// fmt.Printf("Output path: %s\n", outputPath)
	// fmt.Printf("Config: %+v\n", config)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Construct output file name using the output folder name and index
	outputFile := fmt.Sprintf("%s/vid_potrait_tests_%d.mp4", outputPath, index)

	// If the target format is GIF, warn about audio loss
	if strings.ToLower(config.VideoFormat) == "gif" {
		log.Printf("Warning: Converting to GIF format will remove audio from the output file: %s", filepath.Base(inputFile))
		return processGif(inputFile, outputFile, config)
	}

	// Check if the input file has audio
	hasAudio, err := checkForAudio(inputFile)
	if err != nil {
		log.Printf("Warning: Could not determine if file has audio: %v", err)
		// Assume it has audio if we can't determine
		hasAudio = true
	}

	log.Printf("Processing video: %s (has audio: %v)", filepath.Base(inputFile), hasAudio)

	// Set up basic output options
	videoCodec := "libx264"
	audioCodec := "aac"
	outputOptions := ffmpeg.KwArgs{
		"c:v": videoCodec,
		"b:v": fmt.Sprintf("%dk", config.VideoSizeKB),
	}

	// Format-specific settings
	switch strings.ToLower(config.VideoFormat) {
	case "mp4":
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
		} else {
			outputOptions["an"] = ""
		}
	case "webm":
		videoCodec = "libvpx-vp9"
		audioCodec = "libopus"
		outputOptions["c:v"] = videoCodec
		outputOptions["crf"] = "30"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
			outputOptions["b:a"] = "128k"
		} else {
			outputOptions["an"] = ""
		}
	case "avi":
		videoCodec = "libxvid"
		audioCodec = "libmp3lame"
		outputOptions["c:v"] = videoCodec
		outputOptions["q:v"] = "5"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
			outputOptions["q:a"] = "3"
		} else {
			outputOptions["an"] = ""
		}
	case "mov":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		outputOptions["movflags"] = "+faststart"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
		} else {
			outputOptions["an"] = ""
		}
	case "mkv":
		videoCodec = "libx264"
		audioCodec = "libopus"
		outputOptions["c:v"] = videoCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
			outputOptions["b:a"] = "192k"
		} else {
			outputOptions["an"] = ""
		}
	case "flv":
		videoCodec = "flv"
		audioCodec = "libmp3lame"
		outputOptions["c:v"] = videoCodec
		if hasAudio {
			outputOptions["c:a"] = audioCodec
			outputOptions["ar"] = "44100"
		} else {
			outputOptions["an"] = ""
		}
	case "wmv":
		videoCodec = "wmv2"
		audioCodec = "wmav2"
		outputOptions["c:v"] = videoCodec
		if hasAudio {
			outputOptions["c:a"] = audioCodec
			outputOptions["b:a"] = "128k"
		} else {
			outputOptions["an"] = ""
		}
	case "m4v":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		outputOptions["movflags"] = "+faststart"
		if hasAudio {
			outputOptions["c:a"] = audioCodec
		} else {
			outputOptions["an"] = ""
		}
	case "3gp":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["preset"] = "veryfast"
		outputOptions["crf"] = "28"
		outputOptions["strict"] = -2
		if hasAudio {
			outputOptions["c:a"] = audioCodec
		} else {
			outputOptions["an"] = ""
		}
	default:
		return fmt.Errorf("unsupported video format: %s", config.VideoFormat)
	}

	// Ensure the output file has the correct extension
	if !strings.HasSuffix(outputFile, "."+strings.ToLower(config.VideoFormat)) && config.VideoFormat != "mp4" {
		outputFile = strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + "." + strings.ToLower(config.VideoFormat)
	}

	// Create a direct ffmpeg command with explicit mapping
	input := ffmpeg.Input(inputFile)
	scaled := input.Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", config.Width, config.Height)})

	// Create a command that properly handles audio
	var cmd *ffmpeg.Stream

	if hasAudio {
		// Use explicit mapping to preserve audio
		cmd = scaled.Output(outputFile, outputOptions).
			GlobalArgs("-map", "0:v", "-map", "0:a?") // Map video and audio streams
	} else {
		// No audio mapping needed
		cmd = scaled.Output(outputFile, outputOptions)
	}

	// Log the command for debugging
	audioOption := "-an"
	if hasAudio {
		audioOption = "-c:a " + audioCodec
	}
	log.Printf("compiled command: ffmpeg -i %s -filter_complex [0]scale=%d:%d[s0] -map [s0] %s %s %s -y",
		inputFile, config.Width, config.Height,
		audioOption,
		"-b:v "+fmt.Sprintf("%dk", config.VideoSizeKB),
		"-c:v "+videoCodec)

	// Run the command
	err = cmd.OverWriteOutput().Run()
	if err != nil {
		return fmt.Errorf("error processing video: %v", err)
	}

	// Verify the output has audio if the input had audio
	if hasAudio {
		outputHasAudio, verifyErr := checkForAudio(outputFile)
		if verifyErr != nil {
			log.Printf("Warning: Could not verify audio in output file: %v", verifyErr)
		} else if !outputHasAudio {
			log.Printf("Warning: Input had audio but output does not. This may indicate an issue with the conversion.")
		} else {
			log.Printf("Successfully preserved audio in output file: %s", filepath.Base(outputFile))
		}
	}

	return nil
}

// Check if the input file has an audio stream
func checkForAudio(inputFile string) (bool, error) {
	// Use ffprobe to check for audio streams
	probeData, err := ffmpeg.Probe(inputFile)
	if err != nil {
		return false, fmt.Errorf("error running ffprobe: %v", err)
	}

	// If the output contains "codec_type\":\"audio", then the file has an audio stream
	return strings.Contains(probeData, "codec_type\":\"audio"), nil
}

func processGif(inputFile, outputFile string, config Config) error {
	// Create a temporary file for the scaled video
	scaledOutput, err := os.CreateTemp(os.TempDir(), "scaled_*.mp4")
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(scaledOutput.Name())

	// Step 1: Scale the input video
	err = ffmpeg.Input(inputFile).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", config.Width, config.Height)}).
		Output(scaledOutput.Name()).
		OverWriteOutput().
		Run()
	if err != nil {
		return fmt.Errorf("error scaling video: %v", err)
	}

	// Create a temporary file for the palette
	paletteFile, err := os.CreateTemp(os.TempDir(), "palette_*.png")
	if err != nil {
		return fmt.Errorf("error creating palette temp file: %v", err)
	}
	defer os.Remove(paletteFile.Name())

	// Step 2: Generate a color palette from the scaled video
	err = ffmpeg.Input(scaledOutput.Name()).
		Filter("fps", ffmpeg.Args{"fps=10"}). // Set FPS
		Filter("palettegen", ffmpeg.Args{}).  // Generate palette
		Output(paletteFile.Name()).
		OverWriteOutput().
		Run()
	if err != nil {
		return fmt.Errorf("error generating palette: %v", err)
	}

	// Ensure the output file has .gif extension
	if !strings.HasSuffix(outputFile, ".gif") {
		outputFile = strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".gif"
	}

	// Step 3: Create GIF using the palette
	err = ffmpeg.Input(scaledOutput.Name()).
		Filter("fps", ffmpeg.Args{"fps=10"}). // Set FPS for the GIF
		Filter("paletteuse", ffmpeg.Args{}).  // Use the generated palette
		Output(outputFile, ffmpeg.KwArgs{
			"loop": "0", // 0 means loop forever
		}).
		OverWriteOutput().
		Run()
	if err != nil {
		return fmt.Errorf("error converting to GIF: %v", err)
	}

	return nil
}
