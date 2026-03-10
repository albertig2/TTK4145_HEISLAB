@echo off
for /L %%i in (1,1,6) do (
    start cmd /K "cd /d .."
)

@REM for å kjøre: .\startSimTerminal.bat