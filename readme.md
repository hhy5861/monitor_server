## Monitor_server build and deploy
To get a development environment you will need:

* Go 1.8+

Run the following:

* install [Godep](https://github.com/tools/godep): `go get github.com/tools/godep`
* run `godep restore` install rely packge 
* run `make build` to build the binary
* run `make media` to build the media
* copy config in to path `/etc/monitor_servce` [Config](http://gitlab.pnlyy.com/lib/monitor_servce/blob/master/config/config.yaml)
* run `./monitor_servce --help` for options