// 명덕초 학급 대시보드 실행기 (Windows)
// - dashboard.html 을 exe 안에 넣어 두고, 이 컴퓨터(127.0.0.1)에서만 열리는 작은 웹서버로 띄운다.
// - Microsoft Edge(없으면 Chrome)를 '앱 창' 모드로 열고, 창을 닫으면 프로그램도 끝난다.
// - exe 옆에 logo.png 를 두면 학교 로고로 쓴다.
package main

import (
	_ "embed"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed dashboard.html
var page []byte

// 포트를 고정해야 브라우저 저장소(설정·체크 기록)가 유지된다.
const addr = "127.0.0.1:47321"
const url = "http://" + addr + "/"

func main() {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// 이미 실행 중이면 창만 하나 더 연다.
		openWindow()
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/", "/index.html", "/dashboard.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.Write(page)
		case "/logo.png":
			http.ServeFile(w, r, filepath.Join(dir, "logo.png"))
		default:
			http.NotFound(w, r)
		}
	})
	go http.Serve(ln, mux)

	if cmd := openWindow(); cmd != nil {
		cmd.Wait() // 앱 창을 닫으면 종료
		return
	}
	select {} // 기본 브라우저로 연 경우: 계속 실행
}

func openWindow() *exec.Cmd {
	local := os.Getenv("LOCALAPPDATA")
	profile := filepath.Join(local, "MyeongdeokDashboard", "browser")
	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(os.Getenv("ProgramFiles"), `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(os.Getenv("ProgramFiles"), `Google\Chrome\Application\chrome.exe`),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), `Google\Chrome\Application\chrome.exe`),
		filepath.Join(local, `Google\Chrome\Application\chrome.exe`),
	}
	for _, b := range candidates {
		if _, err := os.Stat(b); err != nil {
			continue
		}
		cmd := exec.Command(b, "--app="+url, "--user-data-dir="+profile,
			"--start-maximized", "--no-first-run", "--no-default-browser-check",
			"--autoplay-policy=no-user-gesture-required")
		if cmd.Start() == nil {
			return cmd
		}
	}
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	return nil
}
