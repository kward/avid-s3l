# carbonio

This package provides a UI to control the carbonio device of the Avid S3L. The
UI is accessible via the command-line, a browser (HTTP), and via Open Sound
Control (OSC).

This code is under active development. For more background information about the
CarbonIO device (based on what has been discovered), see
[Avid S3L AVB Testing](https://docs.google.com/document/d/15VqiMrEWSea8XkfREXDXc20SR_aW4acPiuXZWqr5f08/edit?usp=sharing)
(a Google Doc).

## Development and testing

To set up your environment for development or testing, see the
[Environment setup](#env_setup) section below.

```shell
$ mkdir ~/wa
$ cd ~/wa
$ ln -s . src
$ mkdir -p github.com/kward
$ git clone --recursive https://github.com/kward/avid-s3l.git github.com/kward
$ cd github.com/kward/avid-s3l/carbonio
$ go test ./...
```

Update existing submodules.

```shell
$ git submodule update --init
```

### Local SPI devices directory

Create a local SPI device structure.

```shell
$ go run carbonio.go --spi_base_dir /tmp/spi internal create_spi
```

### Run the binary

```shell
$ go run carbonio.go --spi_base_dir /tmp/spi status
LED    STATUS
Power  Off
Status Off
Mute   Off
```

### bindata

bindata (https://github.com/go-bindata/go-bindata) is used to bind binary data
(e.g., HTML files and templates) into the binary, enabling it to be fully
self-contained.

```shell
$ go-bindata -debug -o static/bindata.go -pkg static --ignore=\\.go$ -fs -prefix="${PWD}/static" static/
$ go-bindata -debug -o templates/bindata.go -pkg templates -ignore=\\.go$ -prefix="${PWD}/templates" templates/...
```

## Building the binary

### Update helpers/bindata.go

```shell
$ go-bindata -o static/bindata.go -pkg static --ignore=\\.go$ -fs -prefix="${PWD}/static" static/
$ go-bindata -o templates/bindata.go -pkg templates -ignore=\\.go$ -prefix="${PWD}/templates" templates/...
```

### Compile

The `carbonio.arm7` binary can be built as follows. Obviously this is only
usable on ARM7 capable hardware, which the Diamond Platform in the Avid S3L is.

```shell
$ GOOS=linux GOARM=7 GOARCH=arm go build -o carbonio.arm7 carbonio.go
```

It is also possible to build a `carbonio` binary on macOS or Linux for running
locally. (This likely works under Windows as well, but it hasn't yet been tested
by the author.)

```shell
$ go build carbonio.go
```

## <a name="env_setup"></a>Environment setup

Install the Go language from https://golang.org/dl/.

Install supporting software.

```shell
$ go get -u github.com/go-bindata/go-bindata/...
```
