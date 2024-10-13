build.cli:
	go build -ldflags "-s -w" -o tmp/fbrowse main.go
	cp tmp/fbrowse ~/.local/bin/
