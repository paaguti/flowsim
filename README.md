# flowsim

A flexible IP traffic generator, originally developed in the scope of https://github.com/mami-project/trafic and heavily modified as a pet project afterwards

## Dependencies

A UNIX or UNIX-like OS

### Go

`flowsim` depends on golang-1.13. Follow instructions at https://golang.org/doc/install. We use the modules feature to be able to fix the release of the quic-go fork we want to use.

* *Note*: installing an official binary distribution is recommended

### quic-go

`flowsim` depends on https://github.com/ferrieux/quic-go (v0.7.7), a specfically tuned fork of the `spinvec` branch of https://github.com/ferrieux/quic-go, which, in turn is a fork of https://github.com/lucas-clemente/quic-go. We need this specific version, because it implements the spinbit.


## Build and install

Clone this repository under `$HOME/go/src/github.com/paaguti`:

```
git clone https://github.com/paaguti/flowsim
cd flowsim/
go test
# Ignore the message:
# ?   	github.com/paaguti/flowsim	[no test files]
# The important thing is that all dependencies are installed
#
go build .
go install .
```

# flowsim

iperf3 is a good traffic generator, but it has its limitations. While developing [`trafic`](https://github.com/mami-project/trafic), an [issue](https://github.com/esnet/iperf/issues/768) regarding setting the total bytes transferred on a TCP stream was discovered. In order to accurately simulate web-short and ABR video streams, flowsin was developed. It follows the philosophy of iperf3 (server and client mode in one application).

## flowsim and trafic

I have started this independent github after the project supporting the initial development of trafic (and flowsim) finished. I intend to contribute back to the original trafic github while that is active.

## flowsim modes

`flowsim` can be started as a TCP or QUIC server or client,  or as a UDP source or sink. It supports IPv4 and IPv6 addressing and sets the DSCP field in the IP header of the packets it generates. By default, the server and sink modes use the IPv4 loopback address (`127.0.0.1`) by default. Interface addresses have to be set explicitly.

## flowsim as a TCP or QUIC server

Once started as a server, `flowsim` will basically sit there and wait for the client to request bunches of data over a TCP connection.

```
Usage:
  flowsim server [flags]

Flags:
  -T, --TOS string   Value of the DSCP field in the IP layer (number or DSCP id) (default "CS0")
  -h, --help         help for server
  -I, --ip string    IP address or host name bound to the flowsim server (default "127.0.0.1")
  -1, --one-off      Just accept one connection and quit (default is run until killed)
  -p, --port int     TCP port bound to the flowsim server (default 8081)
  -Q, --quic         Use QUIC (default is TCP)
```

Note in the normal mode, `flowsim` will be executed until killed with a `SIGINT` sinal (i.e. `Control-C` from the keyboard). The `--one-off` option will make `flowsim` quit after a flow has been served.

The size of the TCP PDU served and the moment where a connection is closed are determined by the client.

## flowsim as a TCP or QUIC client

When `flowsim` is started as a client, a number of TCP segments with a fixed size will be requested from the server. All segments will be served over the same TCP connection, which is closed afterwards.

```
Usage:
  flowsim client [flags]

Flags:
  -T, --TOS string     Value of the DSCP field in the IP packets (valid int or DSCP-Id) (default "CS0")
  -N, --burst string   Size of each burst (as x(.xxx)?[kmgtKMGT]?) (default "1M")
  -h, --help           help for client
  -t, --interval int   Interval in secs between bursts (default 10)
  -I, --ip string      IP address or host name of the flowsim server to talk to (default "127.0.0.1")
  -n, --iter int       Number of bursts (default 6)
  -p, --port int       TCP port of the flowsim server (default 8081)
  -Q, --quic           Use QUIC (default is TCP)
```

## flowsim as a UDP source

```
Usage:
  flowsim source [flags]

Flags:  -T, --TOS string      Value of the DSCP field in the IP packets (valid int or DSCP-Id) (default "CS0")
  -h, --help            help for source
  -I, --ip string       IP address or host name of the flowsim UDP sink to talk to (default "127.0.0.1")
  -L, --local string    Outgoing source IP address to use; determins interface (default: empyt-any interface)
  -N, --packet string   Size of each packet (as x(.xxx)?[kmgtKMGT]?) (default "1k")
  -p, --port int        UDP port of the flowsim UDP sink (default 8081)
  -P, --pps int         Packets per second (default 10)
  -t, --time int        Total time sending (default 6)
  -v, --verbose         Print info re. all generated packets```

## flowsim as a UDP sink

```

## flowsim as a UDP sink

In UDP sink mode, following options are available:

```
Usage:
  flowsim sink [flags]

Flags:
  -h, --help        help for sink
  -I, --ip string   IP address or host name to listen on for the flowsim UDP sink (default "127.0.0.1")
  -m, --multi       Stay in the sink forever and print stats for multiple incoming streams
  -p, --port int    UDP port of the flowsim UDP sink (default 8081)
  -v, --verbose     Print per packet info
```

Note that it makes no sense to set the DSCP in this mode and this option is therefore not present and that the default mode is to receive one flow and quit after printing the QoE statistics for that flow.

## Output format

The TCP and QUIC client outputs results in JSON format:

```
{
  "burst" : "1048576",
  "server" : "127.0.0.1:8081",
  "start" : "2019-03-31T10:15:02+02:00",
  "times": [
    { "time" : "10.343846ms", "xferd" : "1048576" , "n" : "1" },
    { "time" : "18.662231ms", "xferd" : "1048576" , "n" : "2" },
    { "time" : "16.11958ms", "xferd" : "1048576" , "n" : "3" },
    { "time" : "16.177647ms", "xferd" : "1048576" , "n" : "4" },
    { "time" : "17.967623ms", "xferd" : "1048576" , "n" : "5" },
    { "time" : "14.527681ms", "xferd" : "1048576" , "n" : "6" }
  ]
}
```

* `burst`: programmed burst size
* `server`: remote IP address
* `start`: start time in RFC format
* `time`: time needed to request and download
* `xferd`: effectively transferred bytes
* `n`: sequence number

The UDP sink also prints result in JSON format

```
{
  "127.0.0.1:53025" : {
    "Delay" :  "5875.82 us",
    "Jitter" : "382.58 us",
    "Loss" : "0",
    "Reorder" : "0",
    "Samples" : "60"
   }
}
```
