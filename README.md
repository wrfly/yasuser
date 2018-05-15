# Yasuser

> Yet another self-hosted URL shortener.

*Short* or *Restore* your URL, like <https://git.io/> or <https://goo.gl/>
but under **your** control.

Master: [![Master Build Status](https://travis-ci.org/wrfly/yasuser.svg?branch=master)](https://travis-ci.org/wrfly/yasuser)
Develop: [![Develop Build Status](https://travis-ci.org/wrfly/yasuser.svg?branch=develop)](https://travis-ci.org/wrfly/yasuser)

## Run

```sh
docker run --name yasuser -dti \
    -p 8084:8080 -e PREFIX=https://your.domain.com \
    wrfly/yasuser
```

Or use the [docker-compose.yml](./docker-compose.yml).

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

## Features

- [x] it works
- [x] blob database
- [x] length and validate
- [ ] memory cache
- [ ] redis database
- [ ] customization
- [ ] TTL of URL
- [ ] rate limit
- [ ] management(auth)
  - [ ] remove(domain or keywords)
  - [ ] blacklist(domain or keywords)
  - [ ] whitelist(domain or keywords)
- [ ] statistic
  - [ ] URL status
  - [ ] runtime metrics
- [ ] UI index
