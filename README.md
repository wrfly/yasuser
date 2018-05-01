# Short-URL

**Short** or **Restore** your URL, like <https://git.io/> or <https://goo.gl/> but for all domains and under your control.

Master: [![Master Build Status](https://travis-ci.org/wrfly/shorturl.svg?branch=master)](https://travis-ci.org/wrfly/shorturl)
Develop: [![Develop Build Status](https://travis-ci.org/wrfly/shorturl.svg?branch=develop)](https://travis-ci.org/wrfly/shorturl)

## Backend

```bash
NAME:
   short-url - short your url

USAGE:
   short-url [global options] command [command options] [arguments...]

VERSION:
   Version: 0.1.0  Commit: 2664d03  Date: 2018-05-01

AUTHOR:
   wrfly <mr.wrfly@gmail.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen value, -l value   listen port number (default: 8082)
   --domain value             short URL prefix(like https://u.kfd.me) (default: "https://u.kfd.me")
   --db-path value, -p value  database path (default: "short-url.db")
   --db-type value, -t value  database type: redis or file (default: "file")
   --redis value, -r value    database path (default: "localhost:6379/0")
   --debug, -d                log level: debug (default: false)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)
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
- [ ] redis database
- [ ] customize URL
- [ ] URL timeout
- [ ] rate limit