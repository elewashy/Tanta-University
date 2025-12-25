"""
Face Mask & Glasses Detector - Main Application
A user-friendly GUI application for detecting masks and glasses in images and videos
"""

import sys
import os
import cv2
import torch
from pathlib import Path
from PyQt5.QtWidgets import (QApplication, QMainWindow, QPushButton, QLabel, 
                             QVBoxLayout, QHBoxLayout, QWidget, QFileDialog, 
                             QComboBox, QSlider, QGroupBox, QMessageBox)
from PyQt5.QtCore import Qt, QThread, pyqtSignal
from PyQt5.QtGui import QImage, QPixmap

from Detector import YOLOV5_Detector


class VideoThread(QThread):
    """Thread for processing video frames"""
    change_pixmap_signal = pyqtSignal(object)
    
    def __init__(self, detector, video_source):
        super().__init__()
        self.detector = detector
        self.video_source = video_source
        self.running = True
        
    def run(self):
        cap = cv2.VideoCapture(self.video_source)
        while self.running:
            ret, frame = cap.read()
            if ret:
                result = self.detector.Detect(frame)
                self.change_pixmap_signal.emit(result)
            else:
                break
        cap.release()
        
    def stop(self):
        self.running = False
        self.wait()


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.detector = None
        self.video_thread = None
        self.init_ui()
        
    def init_ui(self):
        self.setWindowTitle('Face Mask & Glasses Detector')
        self.setGeometry(100, 100, 1200, 800)
        
        # Central widget
        central_widget = QWidget()
        self.setCentralWidget(central_widget)
        main_layout = QHBoxLayout(central_widget)
        
        # Left panel - Controls
        left_panel = self.create_control_panel()
        main_layout.addWidget(left_panel, 1)
        
        # Right panel - Display
        right_panel = self.create_display_panel()
        main_layout.addWidget(right_panel, 3)
        
    def create_control_panel(self):
        panel = QWidget()
        layout = QVBoxLayout(panel)
        
        # Model Settings
        model_group = QGroupBox("Model Settings")
        model_layout = QVBoxLayout()
        
        # Model selection
        model_layout.addWidget(QLabel("Model:"))
        self.model_combo = QComboBox()
        self.model_combo.addItems(['best.pt', 'model/best.pt'])
        model_layout.addWidget(self.model_combo)
        
        # Confidence threshold
        model_layout.addWidget(QLabel("Confidence Threshold:"))
        self.conf_slider = QSlider(Qt.Horizontal)
        self.conf_slider.setMinimum(1)
        self.conf_slider.setMaximum(100)
        self.conf_slider.setValue(25)
        self.conf_label = QLabel("0.25")
        self.conf_slider.valueChanged.connect(
            lambda v: self.conf_label.setText(f"{v/100:.2f}")
        )
        model_layout.addWidget(self.conf_slider)
        model_layout.addWidget(self.conf_label)
        
        # IOU threshold
        model_layout.addWidget(QLabel("IOU Threshold:"))
        self.iou_slider = QSlider(Qt.Horizontal)
        self.iou_slider.setMinimum(1)
        self.iou_slider.setMaximum(100)
        self.iou_slider.setValue(45)
        self.iou_label = QLabel("0.45")
        self.iou_slider.valueChanged.connect(
            lambda v: self.iou_label.setText(f"{v/100:.2f}")
        )
        model_layout.addWidget(self.iou_slider)
        model_layout.addWidget(self.iou_label)
        
        # Load model button
        self.load_model_btn = QPushButton("Load Model")
        self.load_model_btn.clicked.connect(self.load_model)
        model_layout.addWidget(self.load_model_btn)
        
        model_group.setLayout(model_layout)
        layout.addWidget(model_group)
        
        # Detection Options
        detection_group = QGroupBox("Detection Options")
        detection_layout = QVBoxLayout()
        
        self.image_btn = QPushButton("Detect on Image")
        self.image_btn.clicked.connect(self.detect_image)
        self.image_btn.setEnabled(False)
        detection_layout.addWidget(self.image_btn)
        
        self.video_btn = QPushButton("Detect on Video")
        self.video_btn.clicked.connect(self.detect_video)
        self.video_btn.setEnabled(False)
        detection_layout.addWidget(self.video_btn)
        
        self.webcam_btn = QPushButton("Detect on Webcam")
        self.webcam_btn.clicked.connect(self.detect_webcam)
        self.webcam_btn.setEnabled(False)
        detection_layout.addWidget(self.webcam_btn)
        
        self.stop_btn = QPushButton("Stop Detection")
        self.stop_btn.clicked.connect(self.stop_detection)
        self.stop_btn.setEnabled(False)
        detection_layout.addWidget(self.stop_btn)
        
        detection_group.setLayout(detection_layout)
        layout.addWidget(detection_group)
        
        # Status
        self.status_label = QLabel("Status: Ready")
        self.status_label.setWordWrap(True)
        layout.addWidget(self.status_label)
        
        layout.addStretch()
        return panel
        
    def create_display_panel(self):
        panel = QWidget()
        layout = QVBoxLayout(panel)
        
        self.image_label = QLabel("Load a model to start detection")
        self.image_label.setAlignment(Qt.AlignCenter)
        self.image_label.setStyleSheet("QLabel { background-color: #2b2b2b; color: white; }")
        self.image_label.setMinimumSize(800, 600)
        layout.addWidget(self.image_label)
        
        return panel
        
    def load_model(self):
        try:
            self.status_label.setText("Status: Loading model...")
            QApplication.processEvents()
            
            model_path = self.model_combo.currentText()
            if not os.path.exists(model_path):
                QMessageBox.warning(self, "Error", f"Model file not found: {model_path}")
                self.status_label.setText("Status: Model not found")
                return
                
            conf_thres = self.conf_slider.value() / 100
            iou_thres = self.iou_slider.value() / 100
            
            self.detector = YOLOV5_Detector(
                weights=model_path,
                img_size=640,
                confidence_thres=conf_thres,
                iou_thresh=iou_thres,
                agnostic_nms=True,
                augment=True
            )
            
            self.image_btn.setEnabled(True)
            self.video_btn.setEnabled(True)
            self.webcam_btn.setEnabled(True)
            self.status_label.setText("Status: Model loaded successfully")
            
        except Exception as e:
            QMessageBox.critical(self, "Error", f"Failed to load model: {str(e)}")
            self.status_label.setText(f"Status: Error - {str(e)}")
            
    def detect_image(self):
        if not self.detector:
            return
            
        file_path, _ = QFileDialog.getOpenFileName(
            self, "Select Image", "", "Image Files (*.png *.jpg *.jpeg *.bmp)"
        )
        
        if file_path:
            try:
                self.status_label.setText("Status: Processing image...")
                QApplication.processEvents()
                
                img = cv2.imread(file_path)
                result = self.detector.Detect(img)
                self.display_image(result)
                self.status_label.setText("Status: Image processed successfully")
                
            except Exception as e:
                QMessageBox.critical(self, "Error", f"Failed to process image: {str(e)}")
                self.status_label.setText(f"Status: Error - {str(e)}")
                
    def detect_video(self):
        if not self.detector:
            return
            
        file_path, _ = QFileDialog.getOpenFileName(
            self, "Select Video", "", "Video Files (*.mp4 *.avi *.mov)"
        )
        
        if file_path:
            self.start_video_detection(file_path)
            
    def detect_webcam(self):
        if not self.detector:
            return
        self.start_video_detection(0)
        
    def start_video_detection(self, source):
        try:
            self.status_label.setText("Status: Starting video detection...")
            self.video_thread = VideoThread(self.detector, source)
            self.video_thread.change_pixmap_signal.connect(self.display_image)
            self.video_thread.start()
            
            self.image_btn.setEnabled(False)
            self.video_btn.setEnabled(False)
            self.webcam_btn.setEnabled(False)
            self.stop_btn.setEnabled(True)
            self.status_label.setText("Status: Video detection running")
            
        except Exception as e:
            QMessageBox.critical(self, "Error", f"Failed to start video detection: {str(e)}")
            self.status_label.setText(f"Status: Error - {str(e)}")
            
    def stop_detection(self):
        if self.video_thread:
            self.video_thread.stop()
            self.video_thread = None
            
        self.image_btn.setEnabled(True)
        self.video_btn.setEnabled(True)
        self.webcam_btn.setEnabled(True)
        self.stop_btn.setEnabled(False)
        self.status_label.setText("Status: Detection stopped")
        
    def display_image(self, cv_img):
        rgb_image = cv2.cvtColor(cv_img, cv2.COLOR_BGR2RGB)
        h, w, ch = rgb_image.shape
        bytes_per_line = ch * w
        qt_image = QImage(rgb_image.data, w, h, bytes_per_line, QImage.Format_RGB888)
        
        scaled_pixmap = QPixmap.fromImage(qt_image).scaled(
            self.image_label.size(), Qt.KeepAspectRatio, Qt.SmoothTransformation
        )
        self.image_label.setPixmap(scaled_pixmap)
        
    def closeEvent(self, event):
        if self.video_thread:
            self.video_thread.stop()
        event.accept()


def main():
    app = QApplication(sys.argv)
    app.setStyle('Fusion')
    window = MainWindow()
    window.show()
    sys.exit(app.exec_())


if __name__ == '__main__':
    main()
