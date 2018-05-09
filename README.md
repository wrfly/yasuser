# Short-URL

**Short** or **Restore** your URL, like <https://git.io/> or <https://goo.gl/> but for all domains and under your control.

Master: [![Master Build Status](https://travis-ci.org/wrfly/short-url.svg?branch=master)](https://travis-ci.org/wrfly/short-url)
Develop: [![Develop Build Status](https://travis-ci.org/wrfly/short-url.svg?branch=develop)](https://travis-ci.org/wrfly/short-url)

## Run

```sh
docker run --name short-url -dti \
    -p 8084:8080 -v `pwd`:/data \
    -e DB_PATH=/data/short-url.db \
    -e PREFIX=https://u.kfd.me \
    wrfly/short-url
```

Or use the [docker-compose.yml](./docker-compose.yml).

### Help

```bash
NAME:
   short-url - short your url

USAGE:
   short-url [global options] command [command options] [arguments...]

VERSION:
   Version: 0.1.0  Commit: 5d7172a  Date: 2018-05-07

AUTHOR:
   wrfly <mr.wrfly@gmail.com>

COMMANDS:
    help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value     port number (default: 8080) [$PORT]
   --prefix value   short URL prefix (default: "https://u.kfd.me") [$PREFIX]
   --db-path value  database path (default: "short-url.db") [$DB_PATH]
   --db-type value  database type: redis or file (default: "file") [$DB_TYPE]
   --redis value    database path (default: "localhost:6379/0") [$REDIS]
   --debug, -d      log level: debug (default: false) [$DEBUG]
   --help, -h       show help (default: false)
   --version, -v    print the version (default: false)
```

## Usage

```bash
# short your URL
➜  ~ curl https://u.kfd.me -d "https://kfd.me"
https://u.kfd.me/1
➜  ~

# restore it
➜  ~ curl http://u.kfd.me/1
<a href="https://kfd.me">Found</a>.

```

## ToDo

- [x] it works
- [x] blob database
- [x] length and validate
- [ ] redis database
- [ ] customization
- [ ] TTL of URL
- [ ] rate limit
- [ ] manage(auth)
- [ ] statistic
  - [ ] URL status(with auth)
  - [ ] runtime metrics(with auth)
