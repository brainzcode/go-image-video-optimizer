package main

import (
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
	InputPath   string `yaml:"input_path"`
	OutputPath  string `yaml:"output_path"`
	Width       int    `yaml:"width"`
	Height      int    `yaml:"height"`
	ImageSizeKB int    `yaml:"image_size_kb"`
	VideoSizeKB int    `yaml:"video_size_kb"`
	VideoFormat string `yaml:"video_format"`
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
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				processFile(file, config)
			}
		}()
	}

	err := filepath.Walk(config.InputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			jobs <- path
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking through directory: %v\n", err)
	}

	close(jobs)
	wg.Wait()
}

func processFile(file string, config Config) {
	relPath, _ := filepath.Rel(config.InputPath, filepath.Dir(file))
	currentOutputPath := filepath.Join(config.OutputPath, relPath)

	switch getFileType(file) {
	case "image":
		if err := ProcessImage(file, currentOutputPath, config); err != nil {
			log.Printf("Error processing image %s: %v\n", file, err)
		}
	case "video":
		if err := ProcessVideo(file, currentOutputPath, config, int(imageIndex)); err != nil {
			log.Printf("Error processing video %s: %v\n", file, err)
		}
	}
}
