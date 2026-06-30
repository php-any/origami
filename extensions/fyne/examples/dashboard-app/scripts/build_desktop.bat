@echo off
setlocal enabledelayedexpansion

REM Go to origami repository root
cd /D "%~dp0..\..\..\..\.."

echo ============================================
echo Origami Dashboard - Desktop Build
echo ============================================
echo.

set GCC=
if exist "C:\mingw64\bin\gcc.exe" set GCC=C:\mingw64\bin\gcc.exe
if exist "C:\tools\mingw64\bin\gcc.exe" set GCC=C:\tools\mingw64\bin\gcc.exe
if exist "C:\msys64\mingw64\bin\gcc.exe" set GCC=C:\msys64\mingw64\bin\gcc.exe
if exist "C:\msys64\ucrt64\bin\gcc.exe" set GCC=C:\msys64\ucrt64\bin\gcc.exe
if "%GCC%"=="" (
    echo ERROR: gcc not found - required for Fyne desktop
    echo Install: choco install mingw
    echo Or:     https://www.msys2.org/
    pause
    exit /b 1
)
echo gcc: %GCC%
for %%F in ("%GCC%") do set GCCDIR=%%~dpF
set PATH=%GCCDIR%;%PATH%
set CGO_ENABLED=1

set APP_DIR=extensions\fyne\examples\dashboard-app
set BUILD_DIR=%APP_DIR%\build\desktop
set GEN_DIR=%APP_DIR%\build\gen
set OUT=%BUILD_DIR%\dashboard.exe

REM Step 1: Generate icon + Windows resource
echo [1/4] Generating icon...
go run %APP_DIR%\scripts\gen_icon.go
if exist "%APP_DIR%\Icon.ico" (
    echo [2/4] Embedding Windows icon...
    rsrc -ico %APP_DIR%\Icon.ico -arch amd64 -o %APP_DIR%\rsrc_windows_amd64.syso
)

REM Step 2: Compile PHP to Go
echo [3/4] Compiling PHP to Go...
if exist "%GEN_DIR%" rmdir /S /Q "%GEN_DIR%"
go run zy.go compile %APP_DIR% --output="%GEN_DIR%" --pkg=dashboard --entry="%APP_DIR%\app.php"
if %ERRORLEVEL% NEQ 0 (
    echo === COMPILE FAILED ===
    pause
    exit /b 1
)
del /Q "%GEN_DIR%\go.mod" >nul 2>&1

REM Step 3: Build
echo [4/4] Building...
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
go build -ldflags "-H windowsgui" -o %OUT% ./%APP_DIR%/
if %ERRORLEVEL% NEQ 0 (
    echo === BUILD FAILED ===
    pause
    exit /b 1
)

echo.
echo === BUILD SUCCESS ===
echo exe: %cd%\%OUT%
echo.
pause
