build_for_windows:
	GOOS=windows GOARCH=amd64 go build -o saveWatcher.exe main.go