#!/bin/bash

# Script to divide images into smaller folders of 4 images each
SOURCE_DIR="worldover-img"
OUTPUT_BASE_DIR="worldover-img-batches"

# Check if source directory exists
if [ ! -d "$SOURCE_DIR" ]; then
    echo "Error: Source directory '$SOURCE_DIR' not found!"
    exit 1
fi

# Create output base directory
mkdir -p "$OUTPUT_BASE_DIR"

echo "Dividing images from '$SOURCE_DIR' into batches of 4..."

# Count total images
total_images=$(find "$SOURCE_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.webp" -o -iname "*.tiff" -o -iname "*.bmp" -o -iname "*.dng" -o -iname "*.svg" \) | wc -l)

echo "Found $total_images image files"

if [ $total_images -eq 0 ]; then
    echo "No image files found in '$SOURCE_DIR'"
    exit 1
fi

# Calculate number of batches needed
batches_needed=$(( (total_images + 3) / 4 ))  # Round up division
echo "Will create $batches_needed batch folders"

# Initialize counters
batch_num=1
image_count=0

# Create first batch directory
current_batch_dir="$OUTPUT_BASE_DIR/batch_$(printf "%03d" $batch_num)"
mkdir -p "$current_batch_dir"

# Process each image file
find "$SOURCE_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.webp" -o -iname "*.tiff" -o -iname "*.bmp" -o -iname "*.dng" -o -iname "*.svg" \) | while read -r image_file; do
    
    # Copy image to current batch directory
    cp "$image_file" "$current_batch_dir/"
    
    image_count=$((image_count + 1))
    
    echo "Copied $(basename "$image_file") to batch_$(printf "%03d" $batch_num) ($image_count/4)"
    
    # Check if current batch is full
    if [ $((image_count % 4)) -eq 0 ]; then
        echo "Batch $batch_num completed with 4 images"
        batch_num=$((batch_num + 1))
        
        # Create next batch directory if there are more images
        if [ $image_count -lt $total_images ]; then
            current_batch_dir="$OUTPUT_BASE_DIR/batch_$(printf "%03d" $batch_num)"
            mkdir -p "$current_batch_dir"
        fi
    fi
done

echo ""
echo "Image division completed!"
echo "Created $batches_needed batch folders in '$OUTPUT_BASE_DIR'"
echo "Each batch contains up to 4 images"
echo ""
echo "Batch folders:"
ls -la "$OUTPUT_BASE_DIR"

echo ""
echo "To process each batch separately, update your config.yaml input_path to point to individual batch folders:"
echo "Example: input_path: \"./worldover-img-batches/batch_001\"" 