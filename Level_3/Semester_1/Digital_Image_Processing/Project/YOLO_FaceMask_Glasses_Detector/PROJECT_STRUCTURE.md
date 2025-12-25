# 📁 Project Structure

Face Mask & Glasses Detector - Complete project organization.

## 🗂️ Directory Structure

```
FaceMask_Glasses_Detector/
│
├── 📄 README.md                      # Main project documentation
├── 📄 PROJECT_STRUCTURE.md           # This file
│
├── 🎯 Application Files
│   ├── main.py                       # GUI application
│   ├── Detector.py                   # Core detection class
│   ├── new_Detector.py               # Alternative detector
│   └── Detection on Video.py         # Video/webcam script
│
├── 🔨 Build System
│   ├── build_exe.py                  # Build script
│   ├── build_exe.bat                 # Windows automation
│   ├── requirements.txt              # Core dependencies
│   └── requirements_exe.txt          # Build dependencies
│
├── 🚀 Helper Scripts
│   ├── run_gui.bat                   # Launch GUI
│   └── run_webcam.bat                # Launch webcam
│
├── 🤖 Model Files
│   ├── best.pt                       # Trained model weights
│   └── classes.txt                   # Detection classes
│
├── 📚 Documentation (docs/)
│   ├── README.md                     # Documentation index
│   ├── QUICK_START.md               # 5-minute guide
│   ├── INSTALLATION_GUIDE.md        # Setup instructions
│   ├── USER_GUIDE.md                # Complete user manual
│   ├── BUILD_GUIDE.md               # Executable creation
│   ├── DOCUMENTATION.md             # API reference
│   ├── PROJECT_SUMMARY.md           # Project overview
│   └── FAQ.md                       # Common questions
│
├── 🏗️ YOLOv5 Components
│   ├── models/                       # Model architectures
│   │   ├── common.py
│   │   ├── experimental.py
│   │   ├── yolo.py
│   │   └── yolov5*.yaml
│   │
│   ├── utils/                        # Utility functions
│   │   ├── datasets.py
│   │   ├── general.py
│   │   ├── torch_utils.py
│   │   └── ...
│   │
│   └── data/                         # Configuration files
│       └── *.yaml
│
└── 📊 Dataset/
    ├── images/                       # Training images
    └── labels/                       # Annotations
```

## 📖 Documentation Files

### Core Documentation
- **README.md** - Main project documentation with overview, features, and quick start
- **docs/README.md** - Documentation index and navigation

### User Guides
- **docs/QUICK_START.md** - Get running in 5 minutes
- **docs/INSTALLATION_GUIDE.md** - Detailed installation instructions
- **docs/USER_GUIDE.md** - Complete feature documentation

### Developer Guides
- **docs/BUILD_GUIDE.md** - Create standalone executable
- **docs/DOCUMENTATION.md** - API reference and code examples

### Reference
- **docs/PROJECT_SUMMARY.md** - Project overview and architecture
- **docs/FAQ.md** - Frequently asked questions

## 🎯 Quick Navigation

### For End Users
```
README.md → docs/QUICK_START.md → docs/USER_GUIDE.md
```

### For Developers
```
README.md → docs/INSTALLATION_GUIDE.md → docs/DOCUMENTATION.md
```

### For Builders
```
README.md → docs/BUILD_GUIDE.md
```

## 🔗 Repository

**GitHub**: [Tanta University - Digital Image Processing Project](https://github.com/elewashy/Tanta-University/tree/main/Level_3/Semester_1/Digital_Image_Processing/Project/FaceMask_Glasses_Detector)
