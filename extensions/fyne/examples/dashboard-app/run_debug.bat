@echo off
cd /D "%~dp0build\desktop"
dashboard.exe 2>&1
echo.
echo Exit code: %ERRORLEVEL%
pause
