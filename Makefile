DEFAULT_SOCKS=127.0.0.1:1080
#DEFAULT_SOCKS=none 	# short for none
#DEFAULT_SOCKS=tor 		# short for 127.0.0.1:9050
DEFAULT_RADIO=anonradio.net:8000
DEFAULT_MOUNT=anonradio
NAME=rstream
CHUNK="1400"
CUSTOM=TRUE
LDFLAGS=-X main.chunk=${CHUNK} -X main.defaultsocks=${DEFAULT_SOCKS}
LDFLAGS2=-X main.defaultradio=${DEFAULT_RADIO} -X main.defaultmount=${DEFAULT_MOUNT} -X main.custom=${CUSTOM}
all: build

build:
	go build -o "${NAME}" -ldflags "${LDFLAGS}" -ldflags "${LDFLAGS2}"
	