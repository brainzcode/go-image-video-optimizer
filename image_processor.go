package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"github.com/h2non/bimg"
)

var imageIndex uint32

func ProcessImage(inputPath, outputPath string, config Config) error {
	// Add defer for cleanup and recovery
	defer func() {
		// Force garbage collection after processing each image
		runtime.GC()
	}()

	// Validate input parameters
	if inputPath == "" || outputPath == "" {
		return fmt.Errorf("input path and output path cannot be empty")
	}

	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Get the output folder name
	outputFolderName := filepath.Base(outputPath)

	// Generate new file name
	newFileName := fmt.Sprintf("%s_%04d.webp", outputFolderName, atomic.AddUint32(&imageIndex, 1))
	outputFilePath := filepath.Join(outputPath, newFileName)

	// Read the input image
	buffer, err := bimg.Read(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read image: %v", err)
	}

	// Validate buffer
	if len(buffer) == 0 {
		return fmt.Errorf("empty image buffer for file: %s", inputPath)
	}

	// Create a new image from buffer
	newImage := bimg.NewImage(buffer)

	// Get the current image size and metadata
	size, err := newImage.Size()
	if err != nil {
		return fmt.Errorf("failed to get image size: %v", err)
	}

	metadata, err := newImage.Metadata()
	if err != nil {
		return fmt.Errorf("failed to get image metadata: %v", err)
	}

	// Determine if we need to swap width and height based on orientation
	width, height := size.Width, size.Height
	if metadata.Orientation >= 5 && metadata.Orientation <= 8 {
		width, height = height, width
	}

	// Calculate the target aspect ratio
	targetRatio := float64(config.Width) / float64(config.Height)

	// Calculate dimensions for cropping
	var cropWidth, cropHeight int
	if float64(width)/float64(height) > targetRatio {
		// Image is wider than target ratio, crop width
		cropHeight = height
		cropWidth = int(float64(cropHeight) * targetRatio)
	} else {
		// Image is taller than target ratio, crop height
		cropWidth = width
		cropHeight = int(float64(cropWidth) / targetRatio)
	}

	// Ensure crop dimensions are not larger than the original image
	cropWidth = min(cropWidth, width)
	cropHeight = min(cropHeight, height)

	// Calculate crop offsets to center the crop
	left := (width - cropWidth) / 2
	top := (height - cropHeight) / 2

	// Ensure crop area is within image bounds
	if left < 0 || top < 0 || left+cropWidth > width || top+cropHeight > height {
		return fmt.Errorf("invalid crop area: image size: %dx%d, crop area: %d,%d %dx%d",
			width, height, left, top, cropWidth, cropHeight)
	}

	log.Printf("Image size: %dx%d, Crop area: %d,%d %dx%d",
		width, height, left, top, cropWidth, cropHeight)

	// Create options for processing
	options := bimg.Options{
		Width:   cropWidth,
		Height:  cropHeight,
		Crop:    true,
		Gravity: bimg.GravityCentre,
		Quality: 90,
		Type:    bimg.WEBP,
	}

	// Process the image
	processedBuffer, err := newImage.Process(options)
	if err != nil {
		return fmt.Errorf("failed to process image: %v", err)
	}

	// Resize the processed image if necessary
	if cropWidth != config.Width || cropHeight != config.Height {
		resizeOptions := bimg.Options{
			Width:  config.Width,
			Height: config.Height,
		}
		processedBuffer, err = bimg.NewImage(processedBuffer).Process(resizeOptions)
		if err != nil {
			return fmt.Errorf("failed to resize image: %v", err)
		}
	}

	// Optimize the image size
	for len(processedBuffer) > config.ImageSizeKB*1024 && options.Quality > 20 {
		options.Quality -= 5
		processedBuffer, err = bimg.NewImage(buffer).Process(options)
		if err != nil {
			return fmt.Errorf("failed to optimize image: %v", err)
		}
	}

	// Save the processed image
	if err := bimg.Write(outputFilePath, processedBuffer); err != nil {
		return fmt.Errorf("failed to save processed image: %v", err)
	}

	// Clear buffers to help with memory management
	buffer = nil
	processedBuffer = nil

	log.Printf("Successfully processed image: %s -> %s\n", inputPath, outputFilePath)
	return nil
}
