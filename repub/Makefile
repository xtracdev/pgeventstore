containerbin:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o repub
	docker build -t xtracdev/repub:latest .

clean:
	rm -f repub
