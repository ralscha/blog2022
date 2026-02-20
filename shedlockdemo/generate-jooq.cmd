@echo off
setlocal

set COMPOSE_FILE=docker-compose-codegen.yml

echo [1/2] Running jOOQ code generation in Docker...
docker compose -f %COMPOSE_FILE% run --rm codegen
set EXIT_CODE=%ERRORLEVEL%

echo [2/2] Cleaning up containers...
docker compose -f %COMPOSE_FILE% down -v

if %EXIT_CODE% neq 0 (
    echo.
    echo ERROR: Code generation failed with exit code %EXIT_CODE%.
    exit /b %EXIT_CODE%
)

echo.
echo Done! Generated sources: target\generated-sources\jooq
endlocal
