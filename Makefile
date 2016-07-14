soundofdata: main.go embedded.go embed.sh static/mon.html static/mon.js
	go generate
	go build


clean:
	rm -f soundofdata
