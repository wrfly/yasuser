# Benchmark

Use [wrk](https://github.com/wg/wrk) and this
[script](benchmark.lua) for benchmark.

## Env

```txt
OS:   Ubuntu 18.04
CPU:  Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
Mem:  15.9G
DISK: NVMe SSD Controller SM961/PM961
```

## With boltdb file storage

Commands:

```bash
docker run --name yasuser -ti --rm \
    -p 8084:8084 \
    -e YASUSER_SHORTENER_STORE_DBPATH=/data/yasuser.db \
    -v `pwd`:/data \
    wrfly/yasuser
```

Results:

```txt
Running 10s test @ http://localhost:8084/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.78ms    5.11ms  55.03ms   93.84%
    Req/Sec    11.17k     4.79k   20.20k    58.00%
  222561 requests in 10.02s, 29.50MB read
Requests/sec:  22207.48
Transfer/sec:      2.94MB
```

Docker stats:

```txt
CONTAINER           CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
be146e68cd66        18.75%              352.1MiB / 15.42GiB   2.23%               32.7MB / 45.6MB     0B / 827MB          11
```

## With redis storage

Commands:

```bash
# create network
docker network create yasuser

# start redis
docker run --network yasuser --rm -ti redis

# start yasuser
docker run --name yasuser -ti --rm \
    -p 8084:8084 \
    --network yasuser \
    -e YASUSER_SHORTENER_STORE_DBTYPE=redis \
    -e YASUSER_SHORTENER_STORE_REDIS=redis://redis:6379 \
    wrfly/yasuser
```

Results:

```txt
Running 10s test @ http://localhost:8084/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   670.09us    0.95ms  19.74ms   92.58%
    Req/Sec    10.03k     1.86k   13.73k    61.00%
  199613 requests in 10.01s, 25.31MB read
Requests/sec:  19945.86
Transfer/sec:      2.53MB
```

Docker stats (`e35e5bf4be91` is redis):

```txt
CONTAINER           CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
3aa9af759910        0.00%               38.7MiB / 15.42GiB    0.25%               47.8MB / 68.2MB     0B / 0B             15
e35e5bf4be91        0.14%               11.04MiB / 15.42GiB   0.07%               28.5MB / 18.5MB     0B / 0B             4
```

---

The result tested in container maybe impacted be the docker's network. `yasuser` running outside
the container(or use the host network) can have a **great** improvement with redis database:

Commands:

```bash
# start redis
docker run --network host --rm -ti redis

# start yasuser
docker run --name yasuser -ti --rm \
    --network host \
    -e YASUSER_SHORTENER_STORE_DBTYPE=redis \
    wrfly/yasuser
```

Results:

```txt
Running 10s test @ http://localhost:8084/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   468.49us    0.85ms  15.56ms   92.72%
    Req/Sec    17.33k     3.90k   27.56k    67.00%
  344930 requests in 10.01s, 43.74MB read
Requests/sec:  34462.63
Transfer/sec:      4.37MB
```

Docker stats (`8e903278981c` is redis):

```txt
CONTAINER           CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
2e886ff5e6e3        0.00%               41.77MiB / 15.42GiB   0.26%               0B / 0B             0B / 0B             13
8e903278981c        0.19%               17.76MiB / 15.42GiB   0.11%               0B / 0B             0B / 3.92MB         4
```
