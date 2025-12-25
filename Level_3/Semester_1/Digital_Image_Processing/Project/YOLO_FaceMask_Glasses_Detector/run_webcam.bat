@echo off
echo ========================================
echo Face Mask Detector - Webcam Detection
echo ========================================
echo.

echo Starting webcam detection...
echo Press 'q' to quit
echo.

python "Detection on Video.py"

if errorlevel 1 (
    echo.
    echo ERROR: Failed to start detection
    echo.
    echo Possible solutions:
    echo 1. Check webcam connection
    echo 2. Install requirements: pip install -r requirements.txt
    echo 3. See INSTALLATION_GUIDE.md for help
    echo.
    pause
)
