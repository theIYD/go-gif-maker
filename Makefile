build:
	@go build -o bin/gif-maker

compile: 
	@echo "Compiling for linux, mac, windows"
	@GOOS=linux GOARCH=386 go build -o bin/linux32/gif-maker main.go
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/gif-maker main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/windows/gif-maker.exe main.go
	@GOOS=windows GOARCH=386 go build -o bin/windows32/gif-maker.exe main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/macos/gif-maker main.go
	@GOOS=freebsd GOARCH=386 go build -o bin/freebsd386/gif-maker main.go