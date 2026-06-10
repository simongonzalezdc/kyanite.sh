@echo off
echo Installing focus.exe to %PROGRAMFILES%\focus ...
if not exist "%PROGRAMFILES%\focus" mkdir "%PROGRAMFILES%\focus"
copy /Y focus.exe "%PROGRAMFILES%\focus\focus.exe"
echo Done.
