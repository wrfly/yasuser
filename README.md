# short-url

短网址服务

## server

```bash
./short-url -h
```

## client

```bash
# create a short URL
➜  ~ curl https://u.kfd.me -d "https://kfd.me"
https://u.kfd.me/000000
➜  ~ 
# get a short URL
➜  ~ curl https://u.kfd.me/000000   
<a href="https://kfd.me">Found</a>.

```

# design

1. 监听端口，处理两种请求（GET/POST）
1. POST： 创建一个短域名
1. GET: 还原一个短域名
1. short -> handler -> query in DB -> {get long} | {none} -> return
1. long -> handler -> md5(index) -> query in DB -> {get short} | {none -> generate -> cache in DB(short->long)} -> return

