# Image and Video Processor

## Overview

This Go-based Media Processor is a powerful tool designed to batch process both images and videos. It resizes, crops, and converts images to WebP format while optimizing their file size. For videos, it can convert them to a specified format and optimize their size. The processor is highly configurable and can handle multiple files concurrently, making it efficient for large-scale media processing tasks.

## Features

- Batch processing of images and videos
- Resizing images to specified dimensions
- Cropping images while maintaining aspect ratio
- Converting images to WebP format
- Converting videos to a specified format
- Optimizing image and video file sizes
- Concurrent processing for improved performance
- Configurable via YAML file
- Detailed logging for easy debugging

## Prerequisites

- Go 1.16 or higher
- libvips 8.10+ (for image processing)
- FFmpeg (for video processing)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/media-processor.git
   cd media-processor
   ```

2. Install the required Go packages:
   ```
   go mod tidy
   ```

3. Install libvips (if not already installed):
   - On Ubuntu/Debian: `sudo apt-get install libvips-dev`
   - On macOS: `brew install vips`
   - For other systems, refer to the [libvips installation guide](https://github.com/libvips/libvips/wiki/Installation)

4. Install FFmpeg (if not already installed):
   - On Ubuntu/Debian: `sudo apt-get install ffmpeg`
   - On macOS: `brew install ffmpeg`
   - For other systems, download from the [FFmpeg official site](https://ffmpeg.org/download.html)

## Configuration

Create a `config.yaml` file in the project root directory with the following structure:

```yaml
input_path: "/path/to/input/folder"
output_path: "/path/to/output/folder"
width: 1920
height: 1080
image_size_kb: 500
video_size_kb: 5000
video_format: "mp4"
```

Adjust the values according to your requirements.

## Usage

1. Ensure your `config.yaml` file is set up correctly.

2. Run the program:
   ```
   go run main.go
   ```

3. The program will process all images and videos in the specified input folder and save the results in the output folder.

4. Check the `conversion.log` file for detailed information about the processing.

## How It Works

The processor uses concurrent processing to handle multiple files simultaneously:

1. It reads the configuration from `config.yaml`.
2. It sets up logging to `conversion.log`.
3. It walks through the input directory to find all files.
4. It creates a worker pool based on the number of CPU cores available.
5. Each worker processes files concurrently:
   - For images: It resizes, crops, converts to WebP, and optimizes the file size.
   - For videos: It converts to the specified format and optimizes the file size.
6. Processed files are saved in the output directory, maintaining the original folder structure.

## Troubleshooting

If you encounter any issues:

1. Check the `conversion.log` file for error messages and warnings.
2. Ensure that the input files are in a supported format.
3. Verify that you have the necessary permissions to read from the input directory and write to the output directory.
4. Make sure libvips and FFmpeg are correctly installed and accessible.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
