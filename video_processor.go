package main

import (
	"fmt"
	"os"
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

	if strings.ToLower(config.VideoFormat) == "gif" {
		return processGif(inputFile, outputFile, config)
	}

	stream := ffmpeg.Input(inputFile)
	stream = stream.Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", config.Width, config.Height)})

	videoCodec := "libx264"
	audioCodec := "aac"
	outputOptions := ffmpeg.KwArgs{
		"c:v": videoCodec,
		"b:v": fmt.Sprintf("%dk", config.VideoSizeKB),
		"c:a": audioCodec,
	}

	switch strings.ToLower(config.VideoFormat) {
	case "mp4":
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
	case "webm":
		videoCodec = "libvpx-vp9"
		audioCodec = "libopus"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["crf"] = "30"
		outputOptions["b:a"] = "128k"
	case "avi":
		videoCodec = "libxvid"
		audioCodec = "libmp3lame"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["q:v"] = "5"
		outputOptions["q:a"] = "3"
	case "mov":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		outputOptions["movflags"] = "+faststart"
	case "mkv":
		videoCodec = "libx264"
		audioCodec = "libopus"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		outputOptions["b:a"] = "192k"
	case "flv":
		videoCodec = "flv"
		audioCodec = "libmp3lame"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["ar"] = "44100"
	case "wmv":
		videoCodec = "wmv2"
		audioCodec = "wmav2"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["b:a"] = "128k"
	case "m4v":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["preset"] = "medium"
		outputOptions["crf"] = "23"
		outputOptions["movflags"] = "+faststart"
	case "3gp":
		videoCodec = "libx264"
		audioCodec = "aac"
		outputOptions["c:v"] = videoCodec
		outputOptions["c:a"] = audioCodec
		outputOptions["preset"] = "veryfast"
		outputOptions["crf"] = "28"
		outputOptions["strict"] = -2
	default:
		return fmt.Errorf("unsupported video format: %s", config.VideoFormat)
	}

	stream = stream.Output(outputFile, outputOptions)
	return stream.OverWriteOutput().Run()
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
