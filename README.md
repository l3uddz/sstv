# sstv

sstv is a small tool to interface SmoothStreams with various IPTV players and/or Plex DVR.

## Installing sstv

sstv offers [pre-compiled binaries](https://github.com/l3uddz/sstv/releases/latest) for Linux, MacOS and Windows for
each official release. In addition, we also offer a [Docker image](#docker)!

Alternatively, you can build the sstv binary yourself. To build sstv on your system, make sure:

1. Your machine runs Linux, macOS or WSL2
2. You have [Go](https://golang.org/doc/install) installed (1.16 or later)
3. Clone this repository and cd into it from the terminal
4. Run `go build -o sstv ./cmd/sstv` from the terminal

You should now have a binary with the name `sstv` in the root directory of the project. To start sstv, simply
run `./sstv`.

If you need to debug certain sstv behaviour, either add the `-v` flag for debug mode or the `-vv` flag for trace mode to
get even more details about internal behaviour.

## Introduction

sstv assumes some level of competence in regard to configuration.

## Configuration

sstv assumes a `config.yml` will present in the same folder as the binary, or in the user configuration directory:

- Linux - `$HOME/.config/sstv`
- MacOS - `$HOME/Library/Application Support/.config/sstv`
- Windows - `%AppData%\.config\sstv`

#### Sample Configuration

```yaml
public_url: http://localhost:1411
smoothstreams:
  username: username@domain.com
  password: password
  site: viewstvn
  server: deu
```

#### Configuration Options

`public_url`

- This **MUST** be set correctly
- It is used by the playlist generator and various other aspects of the application.
- It defines a location that the IPTV player can use to reach sstv.
- Typically, it would generally be http://system-ip:1411

`smoothstreams:`

- `username` - Your username
- `password` - Your password
- `site` - Site you are a member of (`viewstvn` / `viewss` / `view247`)
- `server` - Server to use for streams

### URL(s)

#### Playlist

Your playlist can be reached at `public_url`/playlist.m3u8

It supports the following query arguments:

- type (`0` = HLS **default** / `1` = MPEG2TS)
- server (Defaults to the server specified in your `config.yml` if not set)

#### EPG

The EPG can be reached at `public_url`/epg.xml

### Docker

sstv has an accompanying docker image which can be found on [Docker Hub](https://hub.docker.com/r/cloudb0x/sstv).

#### Version Tags

sstv's Docker image provides various versions that are available via tags. The `latest` tag usually provides the latest
stable version. Others are considered under development and caution must be exercised when using them.

| Tag | Description |
| :----: | --- |
| latest | Latest stable version from a tagged GitHub release |
| master | Most recent GitHub master commit |

#### Usage

```bash
docker run \
  --name=sstv \
  -e "PUID=1000" \
  -e "PGID=1001" \
  -p 1411:1411 \
  -v "/opt/sstv:/config" \
  --restart=unless-stopped \
  -d cloudb0x/sstv
```

#### Parameters

sstv's Docker image supports the following parameters.

| Parameter | Function |
| :----: | --- |
| `-p 1411:1411` | The port used by sstv's webserver |
| `-e PUID=1000` | The UserID to run the sstv binary as |
| `-e PGID=1000` | The GroupID to run the sstv binary as |
| `-e APP_VERBOSITY=0` | The sstv logging verbosity level to use. (0 = info, 1 = debug, 2 = trace) |
| `-e APP_SSDP=false` | Advertise via SSDP (`true` == default) |
| `-v /config` | sstv's config and log file directory |

#### Cloudbox

The following Docker setup should work for many Cloudbox users.

**WARNING: You still need to configure the `config.yml` file!**

Make sure to replace `DOMAIN.TLD` with your domain and `YOUR_EMAIL` with your email.

```bash
docker run \
  --name=sstv \
  -e "PUID=1000" \
  -e "PGID=1001" \
  -e "VIRTUAL_HOST=sstv.DOMAIN.TLD" \
  -e "VIRTUAL_PORT=1411" \
  -e "LETSENCRYPT_HOST=sstv.DOMAIN.TLD" \
  -e "LETSENCRYPT_EMAIL=YOUR_EMAIL" \
  -v "/opt/sstv:/config" \
  --label="com.github.cloudbox.cloudbox_managed=false" \
  --network=cloudbox \
  --network-alias=sstv  \
  --restart=unless-stopped \
  -d cloudb0x/sstv
```

## Donate

If you find this project helpful, feel free to make a small donation:

- [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

- [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)

- [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

- BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL