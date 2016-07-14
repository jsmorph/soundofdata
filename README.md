A little [Go](https://golang.org/) program that listens for data on
`stdin`, forwards that data via a Web socket, and serves up a little
HTML page that uses the
[Web Audio API](https://www.w3.org/TR/webaudio/) for rendering the
data as sound.

## Usage

Pipe data with lines of the form `seriesid frequency duration value`.
_Seriesid_ is an opaque string, which can't contain a space, to
identify a time series.  _Frequency_ is in Hertz, _duration_ is in
seconds (a float).  _Value_ is an optional float that will be
displayed in a graph on the HTML page.

When a line arrives at listeners' Web pages, the series frequency is
set to _frequency_ for _duration_ seconds.  If a _value_ is given,
it'll be displayed on the series graph.  Otherwise the _frequency_ is
added to the graph.

Then look at [`http://localhost:9080/`](http://localhost:9080/).

## Example

See `demo.sh`.

```Shell
make # (or go generate && go build)
(while true; do 
	FREQ="$((80 + RANDOM % 400))"
	echo a $FREQ 3
	sleep 2
	FREQ="$((200 + RANDOM % 300))"
	echo b $FREQ 0.5
	sleep 0.5
	FREQ="$((180 + RANDOM % 400))"
	Y="$((FREQ/10))"
	echo a $FREQ 1 $Y
	sleep $(echo "scale=3;1 + $((RANDOM % 4000))/1000.0" | bc -l)
done) | ./soundofdata
# http://localhost:9080/
```

## Notes

1. With Firefox, the audio pops some.  Doesn't on Chrome.  I don't
   what's going on.
2. [Named pipes](http://linux.die.net/man/3/mkfifo) make good input
   rendezvous points.
3. Watch out for
   [`stdio` buffering](http://www.pixelbeat.org/programming/stdio_buffering/)
   with your pipes.

