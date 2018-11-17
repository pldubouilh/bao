bao
===============

![bao](https://user-images.githubusercontent.com/760637/44001953-85521f68-9e3b-11e8-8fb6-6a4ddbb5d45d.png)

bao is a KISS ssh tunnel built in go

### features
  * simple systray tool to connect connect to the remote machine
  * secure ; default ed25519 ssh crypto keys, and server key pinning
  * a simple run-once server script to spin new SSH users with *no other* rights than port-forwarding

### why ?
bao makes it simple to share a part of a host you run somewhere for others to use. It was initially designed to provide transport for [gossa](https://github.com/pldubouilh/gossa), but it can also serve for other purposes.

### server config
Just run the `newUser.sh` script to spin a new user on your server with only access to the ports you specify. The generated config will be at `./bao.conf`

```
$ sudo bash newUser.sh self-hosted-service 8000 1234
all done!
```

### client setup
Either download release for [linux](https://github.com/pldubouilh/bao/releases/download/0.0.1/Linux.release) or [mac](https://github.com/pldubouilh/bao/releases/download/0.0.1/Mac.release.zip). Once started, bao will look for a config files in `~/.ssh/bao/`, or `$PWD` (on the mac release it's `bao.app/Contents/MacOS/`).

### client build
```
dep ensure
go build main.go
```

 _note linux_ : linux build would need `libgtk-3-dev` and `libappindicator3-dev`

 _note mac_: mac build should be executed from a mac. The built blob should be copied into the `.app` skeleton in `builds/bao.app/Contents/MacOS/bao`

