# 📥 Installation Guide

Complete guide for installing and setting up the Face Mask & Glasses Detector.

## Table of Contents
- [System Requirements](#system-requirements)
- [Installation Methods](#installation-methods)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## 🖥️ System Requirements

### Minimum Requirements
- **OS**: Windows 10/11, Linux (Ubuntu 18.04+), macOS 10.14+
- **Python**: 3.7 or higher
- **RAM**: 4 GB
- **Storage**: 2 GB free space
- **Processor**: Intel Core i3 or equivalent

### Recommended Requirements
- **OS**: Windows 10/11 64-bit
- **Python**: 3.8 or 3.9
- **RAM**: 8 GB or more
- **Storage**: 5 GB free space
- **Processor**: Intel Core i5 or better
- **GPU**: NVIDIA GPU with CUDA support (for faster processing)

## 📦 Installation Methods

### Method 1: Using the Executable (Easiest)

**For End Users - No Python Required**

1. Download the latest release from GitHub
2. Extract the ZIP file
3. Run `FaceMaskDetector.exe`
4. That's it! No installation needed.

```mermaid
graph LR
    A[Download Release] --> B[Extract ZIP]
    B --> C[Run EXE]
    C --> D[Start Detecting]
    
    style D fill:#90EE90
```

### Method 2: Python Installation (For Developers)

#### Step 1: Install Python

**Windows:**
1. Download Python from [python.org](https://www.python.org/downloads/)
2. Run installer
3. ✅ Check "Add Python to PATH"
4. Click "Install Now"

**Linux:**
```bash
sudo apt update
sudo apt install python3.8 python3-pip python3-venv
```

**macOS:**
```bash
brew install python@3.8
```

#### Step 2: Clone Repository

```bash
git clone https://github.com/elewashy/Tanta-University.git
cd Tanta-University/Level_3/Semester_1/Digital_Image_Processing/Project/FaceMask_Glasses_Detector
```

#### Step 3: Create Virtual Environment

**Windows:**
```bash
python -m venv venv
venv\Scripts\activate
```

**Linux/macOS:**
```bash
python3 -m venv venv
source venv/bin/activate
```

#### Step 4: Install Dependencies

```bash
# Basic installation
pip install -r requirements.txt

# For building executable
pip install -r requirements_exe.txt
```

#### Step 5: Verify Installation

```bash
python -c "import torch; print('PyTorch:', torch.__version__)"
python -c "import cv2; print('OpenCV:', cv2.__version__)"
python -c "import PyQt5; print('PyQt5: OK')"
```

### Method 3: GPU Support (Optional)

For faster processing with NVIDIA GPU:

#### Check CUDA Compatibility

```bash
nvidia-smi
```

#### Install PyTorch with CUDA

**CUDA 11.3:**
```bash
pip install torch==1.10.0+cu113 torchvision==0.11.1+cu113 -f https://download.pytorch.org/whl/cu113/torch_stable.html
```

**CUDA 11.6:**
```bash
pip install torch==1.12.0+cu116 torchvision==0.13.0+cu116 -f https://download.pytorch.org/whl/cu116/torch_stable.html
```

**CUDA 11.7:**
```bash
pip install torch==1.13.0+cu117 torchvision==0.14.0+cu117 -f https://download.pytorch.org/whl/cu117/torch_stable.html
```

## ✅ Verification

### Test Installation

Create a test script `test_installation.py`:

```python
import sys
import torch
import cv2
import numpy as np
from PyQt5.QtWidgets import QApplication

print("="*50)
print("Installation Verification")
print("="*50)

# Python version
print(f"Python: {sys.version}")

# PyTorch
print(f"PyTorch: {torch.__version__}")
print(f"CUDA Available: {torch.cuda.is_available()}")
if torch.cuda.is_available():
    print(f"CUDA Version: {torch.version.cuda}")
    print(f"GPU: {torch.cuda.get_device_name(0)}")

# OpenCV
print(f"OpenCV: {cv2.__version__}")

# NumPy
print(f"NumPy: {np.__version__}")

# PyQt5
app = QApplication(sys.argv)
print("PyQt5: OK")

print("="*50)
print("✅ All dependencies installed successfully!")
print("="*50)
```

Run the test:
```bash
python test_installation.py
```

### Quick Functionality Test

```bash
# Test GUI
python main.py

# Test CLI detection
python Detector.py
```

## 🔧 Troubleshooting

### Issue: "Python not found"

**Solution:**
- Ensure Python is installed
- Add Python to PATH
- Restart terminal/command prompt

**Windows - Add to PATH:**
1. Search "Environment Variables"
2. Edit "Path" variable
3. Add Python installation directory
4. Add Python Scripts directory

### Issue: "pip not found"

**Solution:**
```bash
python -m ensurepip --upgrade
```

### Issue: "torch not found" or Import Error

**Solution:**
```bash
pip uninstall torch torchvision
pip install torch torchvision
```

### Issue: "CUDA out of memory"

**Solution:**
- Reduce image size in detector
- Use CPU mode instead
- Close other GPU applications

### Issue: "OpenCV error: camera not found"

**Solution:**
- Check camera permissions
- Update camera drivers
- Try different camera index (0, 1, 2)

### Issue: "PyQt5 import error"

**Solution:**
```bash
pip uninstall PyQt5 PyQt5-sip
pip install PyQt5
```

### Issue: Slow Installation

**Solution:**
```bash
# Use faster mirror
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
```

## 🔄 Updating

### Update to Latest Version

```bash
cd FaceMask_Glasses_Detector
git pull origin main
pip install -r requirements.txt --upgrade
```

### Update Specific Package

```bash
pip install --upgrade torch
pip install --upgrade opencv-python
```

## 🗑️ Uninstallation

### Remove Virtual Environment

```bash
# Deactivate first
deactivate

# Remove directory
rm -rf venv  # Linux/macOS
rmdir /s venv  # Windows
```

### Remove All Packages

```bash
pip freeze > installed.txt
pip uninstall -r installed.txt -y
```

## 📞 Getting Help

If you encounter issues:

1. Check this guide first
2. Search existing GitHub issues
3. Create a new issue with:
   - Your OS and Python version
   - Error message
   - Steps to reproduce
   - Output of `pip list`

After successful installation:

1. Read the [README.md](../README.md) for project overview
2. Check [USER_GUIDE.md](USER_GUIDE.md) for detailed features
3. Try the example scripts
4. Build your own executable with [BUILD_GUIDE.md](BUILD_GUIDE.md)

---

**Happy Detecting! 🎭**
