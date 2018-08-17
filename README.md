bao
===============

![bao](https://user-images.githubusercontent.com/760637/44001953-85521f68-9e3b-11e8-8fb6-6a4ddbb5d45d.png)


bao is a KISS ssh tunnel built in go. It automatically reconnects when the connection drops and it also it comes in with a handful of nice features, like default ed25519 ssh crypto keys, and server key pinning so you know who you're speaking too.

  * server side: use the script to spin a new user + key pair with *no rights but port-forwarding what you specified*
  * for the clients: a nice and easy systray (or cli) app to reliably port-forward to your host


### why ?
bao makes it simple to share a part of a host you run somewhere for others to use. Nice for running apps unexposed to the internet like a file webserver, your favourite rss reader, etc, and sharing data with the other end whether computer or human.

### server config
Just run the `newUser.sh` script to spin a new user on your server with only access to the ports you specify.

```
$ sudo bash newUser.sh self-hosted-service 8000 1234
all done! conf file is called bao.conf
```

### client setup

Either download release for [linux](https://github.com/pldubouilh/bao/releases/download/0.0.1/Linux.release) or [mac](https://github.com/pldubouilh/bao/releases/download/0.0.1/Mac.release.zip). Once started, bao will look for config files in `~/.ssh/bao/`, or `$PWD`. It's also possible to hardcode a key at build time in `utils.go`.

### client build
Default build comes with the cross platform systray tool. Building the cli tool is just a matter of commenting-out the ui bit in `main.go`.

Build with the following commands. The config file can either be hardcoded at build time - and if nothing's hardcoded, it'll look up in `~/.ssh/bao`, or `$PWD`.

```
dep ensure
go build main.go
```

 _note linux_ : linux build would need `libgtk-3-dev` and `libappindicator3-dev`

 _note mac_: mac build should be executed from a mac. Build output should be set into the `.app` skeleton in `builds/bao.app/Contents/MacOS/bao`

