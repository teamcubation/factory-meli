@echo off
echo ============================
echo  Run Docker Container
echo ============================

docker compose --project-name twitter up -d

echo.
echo  Container esta em execucao.
pause