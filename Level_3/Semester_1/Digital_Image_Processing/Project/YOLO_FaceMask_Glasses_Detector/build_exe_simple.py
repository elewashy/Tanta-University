"""
Simplified build script for creating executable
Run: python build_exe_simple.py
"""

import PyInstaller.__main__
import os
import sys

def build_executable():
    """Build the executable using PyInstaller with minimal dependencies"""
    
    # Get the absolute path to the project directory
    project_dir = os.path.dirname(os.path.abspath(__file__))
    
    # PyInstaller arguments - simplified
    args = [
        'main.py',
        '--name=FaceMaskDetector',
        '--onefile',
        '--windowed',
        
        # Add only essential data files
        f'--add-data=best.pt{os.pathsep}.',
        f'--add-data=models{os.pathsep}models',
        f'--add-data=utils{os.pathsep}utils',
        
        # Essential hidden imports only
        '--hidden-import=torch',
        '--hidden-import=torchvision',
        '--hidden-import=cv2',
        '--hidden-import=numpy',
        '--hidden-import=PyQt5',
        '--hidden-import=PIL',
        '--hidden-import=yaml',
        
        # Exclude unnecessary modules to reduce size
        '--exclude-module=matplotlib',
        '--exclude-module=pandas',
        '--exclude-module=scipy',
        '--exclude-module=tensorflow',
        '--exclude-module=transformers',
        '--exclude-module=sklearn',
        '--exclude-module=pytest',
        '--exclude-module=notebook',
        '--exclude-module=IPython',
        '--exclude-module=jupyter',
        '--exclude-module=streamlit',
        '--exclude-module=gradio',
        '--exclude-module=flask',
        '--exclude-module=django',
        
        # Clean build
        '--clean',
        '--noconfirm',
        
        # Output directories
        '--distpath=dist',
        '--workpath=build',
        '--specpath=.',
    ]
    
    print("="*60)
    print("Building Face Mask Detector Executable")
    print("="*60)
    print("\nThis will take several minutes...")
    print("Building with minimal dependencies for smaller size...\n")
    
    try:
        PyInstaller.__main__.run(args)
        print("\n" + "="*60)
        print("✅ Build completed successfully!")
        print(f"📦 Executable location: {os.path.join(project_dir, 'dist', 'FaceMaskDetector.exe')}")
        print("="*60)
    except Exception as e:
        print(f"\n❌ Error during build: {str(e)}")
        sys.exit(1)


if __name__ == '__main__':
    build_executable()
