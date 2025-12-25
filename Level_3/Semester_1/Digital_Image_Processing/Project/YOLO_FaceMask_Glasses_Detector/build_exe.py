"""
Build script for creating executable
Run: python build_exe.py
"""

import PyInstaller.__main__
import os
import sys

def build_executable():
    """Build the executable using PyInstaller"""
    
    # Get the absolute path to the project directory
    project_dir = os.path.dirname(os.path.abspath(__file__))
    
    # PyInstaller arguments
    args = [
        'main.py',  # Main script
        '--name=FaceMaskDetector',  # Name of the executable
        '--onefile',  # Create a single executable file
        '--windowed',  # No console window (GUI only)
        '--icon=NONE',  # Add icon path if you have one
        
        # Add data files
        '--add-data=best.pt;.',
        '--add-data=classes.txt;.',
        '--add-data=models;models',
        '--add-data=utils;utils',
        
        # Hidden imports
        '--hidden-import=torch',
        '--hidden-import=torchvision',
        '--hidden-import=cv2',
        '--hidden-import=numpy',
        '--hidden-import=PIL',
        '--hidden-import=yaml',
        '--hidden-import=scipy',
        '--hidden-import=matplotlib',
        '--hidden-import=PyQt5',
        
        # Exclude unnecessary modules
        '--exclude-module=pytest',
        '--exclude-module=notebook',
        '--exclude-module=IPython',
        
        # Clean build
        '--clean',
        
        # Output directory
        '--distpath=dist',
        '--workpath=build',
        '--specpath=.',
    ]
    
    print("Building executable...")
    print("This may take several minutes...")
    
    try:
        PyInstaller.__main__.run(args)
        print("\n" + "="*50)
        print("Build completed successfully!")
        print(f"Executable location: {os.path.join(project_dir, 'dist', 'FaceMaskDetector.exe')}")
        print("="*50)
    except Exception as e:
        print(f"\nError during build: {str(e)}")
        sys.exit(1)


if __name__ == '__main__':
    build_executable()
