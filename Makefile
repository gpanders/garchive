LDFLAGS := -s -w

all: garchive

garchive: main.go
	go build -ldflags "$(LDFLAGS)"
