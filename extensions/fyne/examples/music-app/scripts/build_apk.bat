@echo off
setlocal enabledelayedexpansion

REM Go to project root (examples/music-app/)
cd /D "%~dp0.."

echo ============================================
echo Origami Music - APK Build (arm64)
echo ============================================
echo.

set ANDROID_HOME=C:\android-sdk
set ANDROID_NDK_HOME=C:\android-sdk\ndk\27.0.12077973

REM Step 1: Generate icon if missing
if not exist "Icon.png" (
    echo [1/3] Generating icon...
    go run scripts\gen_icon.go
)

REM Step 2: Compile PHP to Go source
echo [2/3] Compiling PHP to Go...
cd /D "%~dp0..\..\..\..\.."
if exist "extensions\fyne\examples\music-app\build\gen" rmdir /S /Q "extensions\fyne\examples\music-app\build\gen"
go run zy.go compile extensions/fyne/examples/music-app --output="extensions/fyne/examples/music-app/build/gen" --pkg=musicapp --entry="extensions/fyne/examples/music-app/app.php"
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo === PHP COMPILE FAILED ===
    pause
    exit /b 1
)
del /Q "extensions\fyne\examples\music-app\build\gen\go.mod" >nul 2>&1

REM Step 3: Build APK
cd /D "%~dp0.."
echo [3/3] Building APK...
fyne package -os android/arm64 -appID com.origami.music -icon Icon.png .
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo === BUILD FAILED ===
    pause
    exit /b 1
)

if not exist "build\android" mkdir "build\android"
move /Y music_app.apk "build\android\music_app.apk" >nul 2>&1

echo.
echo === BUILD SUCCESS ===
echo APK: build\android\music_app.apk
echo.
echo Install: adb install -r build\android\music_app.apk
echo.
pause
