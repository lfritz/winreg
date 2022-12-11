.PHONY: all clean

all: shares.exe

clean:
	rm -f shares.exe

shares.exe: main.go
	GOOS=windows GOARCH=amd64 go build -o shares.exe main.go
