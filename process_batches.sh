#!/bin/bash

# Script to process image batches sequentially
BATCH_BASE_DIR="worldover-img-batches"
OUTPUT_BASE_DIR="worldover-landscape-batches"
CONFIG_FILE="config.yaml"
CONFIG_BACKUP="config.yaml.backup"

# Check if batch directory exists
if [ ! -d "$BATCH_BASE_DIR" ]; then
    echo "Error: Batch directory '$BATCH_BASE_DIR' not found!"
    echo "Please run ./divide_images.sh first"
    exit 1
fi

# Create backup of original config
cp "$CONFIG_FILE" "$CONFIG_BACKUP"
echo "Created backup of config.yaml as config.yaml.backup"

# Count total batches
total_batches=$(ls -1 "$BATCH_BASE_DIR" | grep "batch_" | wc -l)
echo "Found $total_batches batches to process"

# Process each batch
for batch_dir in "$BATCH_BASE_DIR"/batch_*; do
    if [ -d "$batch_dir" ]; then
        batch_name=$(basename "$batch_dir")
        echo ""
        echo "============================================"
        echo "Processing $batch_name..."
        echo "============================================"
        
        # Count images in this batch
        image_count=$(ls -1 "$batch_dir"/*.{jpg,jpeg,png,webp,tiff,bmp,dng,svg} 2>/dev/null | wc -l)
        echo "Images in batch: $image_count"
        
        # Update config for this batch
        sed "s|input_path: \".*\"|input_path: \"./$batch_dir\"|g" "$CONFIG_BACKUP" > "$CONFIG_FILE"
        sed -i '' "s|output_path: \".*\"|output_path: \"./$OUTPUT_BASE_DIR/$batch_name\"|g" "$CONFIG_FILE"
        
        # Show current config
        echo "Input:  ./$batch_dir"
        echo "Output: ./$OUTPUT_BASE_DIR/$batch_name"
        
        # Run the image processor
        echo "Starting processing..."
        start_time=$(date +%s)
        
        ./app
        
        end_time=$(date +%s)
        duration=$((end_time - start_time))
        
        echo "Batch $batch_name completed in ${duration} seconds"
        
        # Small delay between batches to allow system cleanup
        sleep 2
    fi
done

# Restore original config
cp "$CONFIG_BACKUP" "$CONFIG_FILE"
echo ""
echo "============================================"
echo "All batches processed successfully!"
echo "Original config.yaml restored"
echo "Check $OUTPUT_BASE_DIR for processed images"
echo "============================================" 