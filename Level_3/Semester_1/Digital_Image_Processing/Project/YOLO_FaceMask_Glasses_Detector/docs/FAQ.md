# ❓ Frequently Asked Questions (FAQ)

Common questions and answers about the Face Mask & Glasses Detector.

## 📋 Table of Contents
- [General Questions](#general-questions)
- [Installation & Setup](#installation--setup)
- [Usage & Features](#usage--features)
- [Performance & Optimization](#performance--optimization)
- [Building & Distribution](#building--distribution)
- [Troubleshooting](#troubleshooting)
- [Advanced Topics](#advanced-topics)

## 🌟 General Questions

### What is this project?

A real-time detection system for identifying people wearing face masks and glasses using YOLOv5 deep learning. It includes a GUI application and can be built as a standalone executable.

### What can I detect?

Currently, the model detects:
- Face masks
- Glasses

### Do I need Python installed?

**For end users**: No, if you use the standalone executable (.exe)
**For developers**: Yes, Python 3.7+ is required

### Is it free to use?

Yes! This project is licensed under MIT License, free for personal and commercial use.

### What platforms are supported?

- **Windows**: Full support (10/11)
- **Linux**: Python version supported
- **macOS**: Python version supported
- **Executable**: Windows only

### Do I need a GPU?

No, but recommended for better performance:
- **CPU**: Works fine, slower processing
- **GPU**: Much faster, especially for video/webcam

## 🔧 Installation & Setup

### How do I install this?

**Quick method:**
```bash
git clone https://github.com/elewashy/Tanta-University.git
cd Tanta-University/Level_3/Semester_1/Digital_Image_Processing/Project/FaceMask_Glasses_Detector
pip install -r requirements.txt
python main.py
```

See [Installation Guide](INSTALLATION_GUIDE.md) for details.

### What are the system requirements?

**Minimum:**
- Python 3.7+
- 4 GB RAM
- 2 GB storage

**Recommended:**
- Python 3.8+
- 8 GB RAM
- NVIDIA GPU with CUDA

### How do I install with GPU support?

```bash
# Check CUDA version
nvidia-smi

# Install PyTorch with CUDA
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu118
```

See [Installation Guide](INSTALLATION_GUIDE.md) → GPU Support.

### Installation fails with "No module named 'torch'"

```bash
# Reinstall PyTorch
pip uninstall torch torchvision
pip install torch torchvision
```

### Can I use a virtual environment?

Yes, recommended:
```bash
python -m venv venv
venv\Scripts\activate  # Windows
pip install -r requirements.txt
```

## 💻 Usage & Features

### How do I start the GUI?

**Windows:**
```bash
run_gui.bat
```

**Or:**
```bash
python main.py
```

### How do I use my webcam?

**Method 1 (GUI):**
1. Run `python main.py`
2. Click "Load Model"
3. Click "Detect on Webcam"

**Method 2 (Script):**
```bash
run_webcam.bat
```

### How do I process a video file?

**GUI:**
1. Click "Detect on Video"
2. Select your video file

**Script:**
Edit `Detection on Video.py`:
```python
video_path = 'path/to/your/video.mp4'
```

### How do I detect on a single image?

**GUI:**
1. Click "Detect on Image"
2. Select your image

**Script:**
Edit `Detector.py`:
```python
img = cv2.imread('path/to/your/image.jpg')
```

### Can I process multiple images at once?

Yes! See [User Guide](USER_GUIDE.md) → Batch Processing for code example.

### How do I save the results?

```python
# After detection
cv2.imwrite('output.jpg', result)
```

For video output, see [User Guide](USER_GUIDE.md) → Video Output.

### Can I adjust detection sensitivity?

Yes, in GUI or code:

**More sensitive** (more detections):
```python
confidence_thres=0.2
```

**Less sensitive** (fewer false positives):
```python
confidence_thres=0.6
```

### What do the colored boxes mean?

- **Red box**: Mask detected
- **Blue box**: Glasses detected (in some versions)
- **Label**: Shows the detected class name

### How accurate is the detection?

Accuracy depends on:
- Image quality
- Lighting conditions
- Distance from camera
- Model confidence threshold

Typical accuracy: 85-95% in good conditions.

## ⚡ Performance & Optimization

### Detection is too slow, how can I speed it up?

**Quick fixes:**
1. Reduce image size:
   ```python
   img_size=320  # Instead of 640
   ```

2. Disable augmentation:
   ```python
   augment=False
   ```

3. Use GPU if available

4. Lower video resolution

### How can I improve detection accuracy?

1. **Better lighting**: Use front-facing, even lighting
2. **Higher resolution**: Use `img_size=640` or higher
3. **Adjust confidence**: Try different thresholds
4. **Enable augmentation**: Set `augment=True`
5. **Better camera**: Use higher quality camera

### What's the best image size?

| Size | Speed | Accuracy | Use Case |
|------|-------|----------|----------|
| 320 | Fast | Lower | Real-time, low-end hardware |
| 640 | Medium | Good | Recommended, balanced |
| 1280 | Slow | Best | High accuracy needed |

### Can I use multiple cameras?

Yes, change camera index:
```python
cap = cv2.VideoCapture(0)  # First camera
cap = cv2.VideoCapture(1)  # Second camera
```

### How much RAM does it use?

Typical usage:
- **Model loading**: ~500 MB
- **Processing**: ~1-2 GB
- **GUI**: ~200 MB

Total: ~2-3 GB

## 📦 Building & Distribution

### How do I create an executable?

**Quick method:**
```bash
build_exe.bat
```

**Manual:**
```bash
pip install pyinstaller
python build_exe.py
```

See [Build Guide](BUILD_GUIDE.md) for details.

### Where is the executable created?

```
dist/FaceMaskDetector.exe
```

### How big is the executable?

Typically 300-500 MB (includes Python, PyTorch, OpenCV, and model).

### Can I reduce the executable size?

Yes:
1. Use `--onedir` instead of `--onefile`
2. Exclude unnecessary modules
3. Use UPX compression

See [Build Guide](BUILD_GUIDE.md) → Size Optimization.

### Can I distribute the executable?

Yes! The executable can run on any Windows 10/11 machine without Python.

**Distribution package:**
```
FaceMaskDetector/
├── FaceMaskDetector.exe
├── best.pt
├── classes.txt
└── README.txt
```

### Do users need to install anything?

No, if using the executable. Just:
1. Extract files
2. Run FaceMaskDetector.exe

### Can I create a Linux executable?

Yes, but requires building on Linux:
```bash
pyinstaller --onefile main.py
```

### Can I add a custom icon?

Yes, edit `build_exe.py`:
```python
'--icon=myicon.ico',
```

## 🐛 Troubleshooting

### "Model file not found" error

**Solution:**
- Ensure `best.pt` is in the same folder as the script
- Check the path in code:
  ```python
  weights='best.pt'  # or 'model/best.pt'
  ```

### "CUDA out of memory" error

**Solutions:**
1. Reduce image size: `img_size=320`
2. Use CPU mode: `device='cpu'`
3. Close other GPU applications
4. Process smaller batches

### Webcam not detected

**Solutions:**
1. Check camera permissions
2. Try different index:
   ```python
   cap = cv2.VideoCapture(1)  # Try 1, 2, etc.
   ```
3. Update camera drivers
4. Test camera in other apps

### "No module named 'cv2'" error

```bash
pip install opencv-python
```

### "No module named 'PyQt5'" error

```bash
pip install PyQt5
```

### GUI window is blank

**Solutions:**
1. Click "Load Model" first
2. Check if model loaded successfully
3. Try running from command line to see errors

### Detection shows no results

**Solutions:**
1. Lower confidence threshold
2. Check image quality
3. Ensure good lighting
4. Verify model is loaded

### "DLL load failed" error

**Solution:**
Install Visual C++ Redistributable from Microsoft.

### Executable won't start

**Solutions:**
1. Build with `--console` to see errors
2. Check if antivirus is blocking
3. Run as administrator
4. Rebuild executable

### Video playback is choppy

**Solutions:**
1. Reduce image size
2. Lower video resolution
3. Use GPU
4. Close other applications

## 🎓 Advanced Topics

### Can I train my own model?

Yes! You'll need:
1. Custom dataset with annotations
2. YOLOv5 training setup
3. GPU for training

See YOLOv5 documentation for training guide.

### Can I add more detection classes?

Yes, modify `Detector.py`:
```python
label = [
    'mask',
    'glasses',
    'face_shield',  # Add new class
    'helmet'        # Add new class
]
```

Then retrain the model with new classes.

### Can I integrate this into my application?

Yes! Import the detector:
```python
from Detector import YOLOV5_Detector

detector = YOLOV5_Detector(...)
result = detector.Detect(image)
```

### Can I use this with a web application?

Yes, you can:
1. Create Flask/FastAPI endpoint
2. Accept image uploads
3. Run detection
4. Return results as JSON

### Can I run this on Raspberry Pi?

Yes, but:
- Performance will be slower
- Use smaller image size (320)
- Disable augmentation
- Consider using lighter model

### How do I get detection coordinates?

Modify `Detector.py` to return coordinates:
```python
detections = []
for *xyxy, conf, cls in reversed(det):
    detections.append({
        'class': label[int(cls)],
        'confidence': float(conf),
        'bbox': [int(x) for x in xyxy]
    })
return detections
```

### Can I use this commercially?

Yes! MIT License allows commercial use. Just include the license file.

### How do I contribute?

1. Fork the repository
2. Create feature branch
3. Make changes
4. Submit pull request

See [README.md](README.md) → Contributing.

### Where can I get help?

1. Check documentation
2. Search [GitHub Issues](https://github.com/elewashy/Tanta-University/issues)
3. Create new issue with details

## 📊 Performance Benchmarks

### Typical Processing Speed

| Hardware | Image Size | FPS |
|----------|-----------|-----|
| CPU (i5) | 320 | 15-20 |
| CPU (i5) | 640 | 5-10 |
| GPU (GTX 1060) | 320 | 60+ |
| GPU (GTX 1060) | 640 | 30-40 |
| GPU (RTX 3060) | 640 | 60+ |

### Memory Usage

| Component | RAM Usage |
|-----------|-----------|
| Model | ~500 MB |
| Processing | ~1-2 GB |
| GUI | ~200 MB |
| **Total** | **~2-3 GB** |

## 🔗 Additional Resources

### Documentation
- [Quick Start Guide](QUICK_START.md)
- [Installation Guide](INSTALLATION_GUIDE.md)
- [User Guide](USER_GUIDE.md)
- [Build Guide](BUILD_GUIDE.md)
- [Documentation Index](DOCUMENTATION.md)

### External Resources
- [YOLOv5 Documentation](https://docs.ultralytics.com/)
- [PyTorch Documentation](https://pytorch.org/docs/)
- [OpenCV Documentation](https://docs.opencv.org/)

## ❓ Still Have Questions?

If your question isn't answered here:

1. **Search Documentation**: Check all guides
2. **GitHub Issues**: Search existing issues
3. **Create Issue**: Open new issue with:
   - Your question
   - What you've tried
   - System information
   - Error messages (if any)

---

**Can't find your answer?** [Open an issue](https://github.com/elewashy/Tanta-University/issues) and we'll help!
