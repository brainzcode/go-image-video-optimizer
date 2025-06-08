#!/bin/bash

echo "Starting image/video optimizer with memory monitoring..."
echo "Monitor memory usage with: watch -n 1 'ps aux | grep app'"
echo ""

# Run the application
./app

echo ""
echo "Processing completed. Check conversion.log for details." 