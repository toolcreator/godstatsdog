# go docker stats daemon (godstatsdog)

A rewrite and enhancement of [dstatsd](https://github.com/toolcreator/dstatsd) in [Go](https://golang.org/).

godstatsdog collects the metrics also provided by the
[`docker stats`](https://docs.docker.com/engine/reference/commandline/stats/) command (see below) from your running
containers and exposes them via HTTP in a format scrapable by
[Prometheus](https://github.com/prometheus/prometheus).
A multiarch (amd64, arm, arm64) docker image is available at
[dockerhub](https://hub.docker.com/r/toolcreator/godstatsdog).

Supported metrics:

- [x] CPU %
- [x] MEM %
- [x] MEM USAGE
- [x] MEM LIMIT
- [x] NET I/O
- [x] BLOCK I/O
- [x] PIDS

## Usage

```shell
$ godstatsdog -help
Usage of godstatsdog:
  -interval int
        The update interval in milliseconds (default 1000)
  -port int
        The port godstatsdog listens on (default 8080)
```

When running, godstatsdog listens to the specified port and provides the metrics at the `/metrics` endpoint.

### via Docker

```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock:ro toolcreator/godstatsdog
```

The arguments shown above can also be passed, e.g.:

```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock:ro toolcreator/godstatsdog -interval 10000 -port 12345
```

Using the `-port` option may be particularly useful when running the container with `--network=host`
(i.e., when port mapping is not available).
When available, port mapping can of course be used as well, e.g.:

```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock:ro -p 12345:8080 toolcreator/godstatsdog
```

#### docker-compose

```yml
godstatsdog:
  image: toolcreator/godstatsdog
  command:
    - "-interval=10000"
  ports:
    - '12345:8080'
  volumes:
    - '/var/run/docker.sock:/var/run/docker.sock:ro'
```

Or, with `network_mode: "host"`:

```yml
godstatsdog:
  image: toolcreator/godstatsdog
  network_mode: "host"
  command:
    - "-interval=10000"
    - "-port=12345"
  volumes:
    - '/var/run/docker.sock:/var/run/docker.sock:ro'
```

### Without Docker

1. Download the source code: `git clone https://github.com/toolcreator/godstatsdog.git`
2. Enter the root directory: `cd godstatsdog`
3. Install dependencies: `go get ./...`
4. Compile: `go build`
5. Install: `go install`
6. Run: `godstatsdog`

You may also skip step 5 and use `./godstatsdog` to run the program instead.

## Metrics

All metrics are of type [gauge](https://prometheus.io/docs/concepts/metric_types/#gauge)
and are labeled with `container_id` and `container_name`.

| Name                                  | Description                                                                 |
| ------------------------------------- | --------------------------------------------------------------------------- |
| godstatsdog_cpu_percent               | The percentage of the host’s CPU the container is using                     |
| godstatsdog_memory_usage_bytes        | The total amount of memory the container is using                           |
| godstatsdog_memory_limit_bytes        | The total amount of memory the container is allowed to use                  |
| godstatsdog_memory_percent            | The percentage of the host’s memory the container is using                  |
| godstatsdog_network_received_bytes    | The amount of data the container has received over its network interface    |
| godstatsdog_network_transmitted_bytes | The amount of data the container has transmitted over its network interface |
| godstatsdog_block_read_bytes          | The amount of data the container has read from block devices on the host    |
| godstatsdog_block_written_bytes       | The amount of data the container has written to block devices on the host   |
| godstatsdog_process_ids               | The number of processes or threads the container has created                |
