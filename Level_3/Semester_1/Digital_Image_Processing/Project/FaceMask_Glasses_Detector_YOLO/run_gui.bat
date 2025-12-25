@echo off
echo ========================================
echo Face Mask Detector - GUI Application
echo ========================================
echo.

echo Starting application...
python main.py

if errorlevel 1 (
    echo.
    echo ERROR: Failed to start application
    echo.
    echo Possible solutions:
    echo 1. Install requirements: pip install -r requirements.txt
    echo 2. Check Python installation: python --version
    echo 3. See INSTALLATION_GUIDE.md for help
    echo.
    pause
)
