all: media

clean:
	@rm -f monitor_servce-linux-*

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v . && mv monitor_server monitor_server-linux-86-64