package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InputPath         string `yaml:"input_path"`
	OutputPath        string `yaml:"output_path"`
	Width             int    `yaml:"width"`
	Height            int    `yaml:"height"`
	ImageSizeKB       int    `yaml:"image_size_kb"`
	VideoSizeKB       int    `yaml:"video_size_kb"`
	VideoFormat       string `yaml:"video_format"`
	ConversionTimeout int    `yaml:"conversion_timeout"` // Timeout in milliseconds
	MaxBatchSize      int    `yaml:"max_batch_size"`     // Maximum files per batch
	MaxWorkers        int    `yaml:"max_workers"`        // Maximum concurrent workers
}

func main() {
	// Read configuration
	config, err := readConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Set up logging
	logFile, err := os.OpenFile("conversion.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Starting media conversion")
	start := time.Now()

	processFiles(config)

	duration := time.Since(start)
	seconds := duration.Seconds()
	log.Printf("Total conversion time: %.2f seconds\n", seconds)
	fmt.Printf("Total conversion time: %.2f seconds\n", seconds) // Print to console as well
	log.Println("Media conversion completed")
}

func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	return config, err
}

func processFiles(config Config) {
	// Collect all files first
	var allFiles []string
	err := filepath.Walk(config.InputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking through directory: %v\n", err)
		return
	}

	log.Printf("Found %d files to process\n", len(allFiles))

	// Process files in batches to prevent memory exhaustion
	batchSize := config.MaxBatchSize
	if batchSize <= 0 {
		batchSize = 5 // Default to 5 files at a time
	}

	numWorkers := config.MaxWorkers
	if numWorkers <= 0 {
		numWorkers = 1 // Default to 1 worker for image processing to avoid memory issues
	}

	for i := 0; i < len(allFiles); i += batchSize {
		end := i + batchSize
		if end > len(allFiles) {
			end = len(allFiles)
		}

		batch := allFiles[i:end]
		log.Printf("Processing batch %d-%d of %d files\n", i+1, end, len(allFiles))

		processBatch(batch, config, numWorkers)

		// Force garbage collection between batches
		runtime.GC()
		runtime.GC() // Call twice to ensure cleanup

		// Small delay to allow system cleanup
		time.Sleep(100 * time.Millisecond)
	}
}

func processBatch(files []string, config Config, numWorkers int) {
	var wg sync.WaitGroup
	jobs := make(chan string, len(files))
	var videoIndex int
	var videoIndexLock sync.Mutex

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				videoIndexLock.Lock()
				currentVideoIndex := videoIndex
				videoIndex++
				videoIndexLock.Unlock()

				processFileWithTimeout(file, config, currentVideoIndex)
			}
		}()
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
}

func processFileWithTimeout(file string, config Config, videoIndex int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConversionTimeout)*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		processFile(file, config, videoIndex)
		done <- true
	}()

	select {
	case <-ctx.Done():
		log.Printf("File processing timed out for %s\n", file)
	case <-done:
		// Successfully completed within the timeout
	}
}

func processFile(file string, config Config, videoIndex int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occurred while processing file %s: %v\n", file, r)
		}
	}()

	log.Printf("Processing file: %s (type: %s)\n", file, GetFileType(file))

	relPath, _ := filepath.Rel(config.InputPath, filepath.Dir(file))
	currentOutputPath := filepath.Join(config.OutputPath, relPath)

	switch GetFileType(file) {
	case "image":
		if err := ProcessImage(file, currentOutputPath, config); err != nil {
			log.Printf("Error processing image %s: %v\n", file, err)
		}
	case "video":
		// Pass the videoIndex to ensure unique file names
		if err := ProcessVideo(file, currentOutputPath, config, videoIndex); err != nil {
			log.Printf("Error processing video %s: %v\n", file, err)
		}
	case "unknown":
		log.Printf("Skipping unknown file type: %s\n", file)
	}
}
