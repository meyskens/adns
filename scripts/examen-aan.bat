@echo off

copy adns.exe C:\Temp\
start C:\Temp\adns.exe

set scriptFileName=%~n0
set scriptFolderPath=%~dp0
set powershellScriptFileName=%scriptFileName%.ps1

PowerShell -NoProfile -ExecutionPolicy Bypass -Command "& {Start-Process PowerShell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File ""%scriptFolderPath%%powershellScriptFileName%""' -Verb RunAs}"
