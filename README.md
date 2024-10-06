# Image Processor

## Overview

This Go-based Image Processor is a powerful tool designed to batch process images. It resizes, crops, and converts images to WebP format while optimizing their file size. The processor is highly configurable and can handle multiple images concurrently, making it efficient for large-scale image processing tasks.

## Features

- Batch processing of images
- Resizing images to specified dimensions
- Cropping images while maintaining aspect ratio
- Converting images to WebP format
- Optimizing image file size
- Concurrent processing for improved performance
- Configurable via YAML file
- Detailed logging for easy debugging

## Prerequisites

- Go 1.16 or higher
- libvips 8.10+

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/image-processor.git
   cd image-processor
   ```

2. Install the required Go packages:
   ```
   go mod tidy
   ```

3. Install libvips (if not already installed):
   - On Ubuntu/Debian: `sudo apt-get install libvips-dev`
   - On macOS: `brew install vips`
   - For other systems, refer to the [libvips installation guide](https://github.com/libvips/libvips/wiki/Installation)

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

3. The program will process all images in the specified input folder and save the results in the output folder.

4. Check the `conversion.log` file for detailed information about the processing.

## Troubleshooting

If you encounter any issues:

1. Check the `conversion.log` file for error messages and warnings.
2. Ensure that the input images are in a supported format.
3. Verify that you have the necessary permissions to read from the input directory and write to the output directory.
4. Make sure libvips is correctly installed and accessible.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
