#!/bin/sh
# Windows용 exe 빌드: ../dashboard.html 을 넣어서 dist/ 에 만든다.
set -e
cd "$(dirname "$0")"
cp ../dashboard.html dashboard.html
mkdir -p ../dist
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-H=windowsgui -s -w" -o ../dist/MyeongdeokDashboard.exe .
rm dashboard.html
echo "built: dist/MyeongdeokDashboard.exe"
