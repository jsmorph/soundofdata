package main

var monhtml = []byte(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Sound of Data</title>
    <meta name="viewport" content="initial-scale=1.0">
    <meta charset="utf-8">
    <style>
      html, body {
      height: 100%;
      margin: 1em;
      padding: 0;
      }
    </style>
  </head>
  <body>
    <input id="volume" type="range" min="0" max="1" step="0.1" value="0.3"/>
    <div id="charts"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/smoothie/1.27.0/smoothie.min.js"></script>
  <script>
var wsURL = "ws://" + location.host + "/ws";

var ctx = new(window.AudioContext || window.webkitAudioContext)();
var initialGain = 0.3;
var maxGraphs = 5;
var graphCount = 0;

function init() {

    var ws = new WebSocket(wsURL);

    ws.onmessage = function(event) {
	var m = JSON.parse(event.data)
	var line = m.line.trim();
	var data = line.split(" ");
	var id = data[0];
	var freq = parseFloat(data[1]);
	var duration = parseFloat(data[2]);
	var value = freq;
	if (4 <= data.length) {
	    value = parseFloat(data[3]);
	}
	tone(id, freq, duration, value);
    }

    ws.onopen = function(event) {
    }

    ws.onclose = function(event) {
	alert("connection closed");
    }

    ws.onerror = function(event) {
	alert("connection error: " + event);
    }

    document.getElementById('volume').addEventListener('change', function() {
	initialGain = parseFloat(this.value);
    });
}

window.onload = init;

var streams = {}; // name -> gain, osc, ticker function

function gesture() {
}

function tone(name, freq, dur, x) {
    var t = streams[name];

    var now = ctx.currentTime;

    if (!t) {
	var g = ctx.createGain();
	g.gain.setValueAtTime(initialGain, now);
	g.connect(ctx.destination);

	var o = ctx.createOscillator();
	o.type = "sine";
	o.frequency.value = freq;
	o.detune.value = 0;
	o.connect(g);
	o.start(0);

	t = {
	    gain: g, 
	    osc: o,
        };

	if (graphCount < maxGraphs) {
	    t.ticker = createTimeline(name);
	}

	streams[name] = t;
    } else {
	// t.osc.frequency.setTargetAtTime(freq, ctx.currentTime+ 0.01, 1);
	t.gain.gain.cancelScheduledValues(now);
	t.osc.frequency.setValueAtTime(freq, now);
	t.gain.gain.setTargetAtTime(initialGain, now, 0.5);
    }
    t.gain.gain.setTargetAtTime(0.0, ctx.currentTime+dur, 0.5);
    if (x) {
	t.ticker(x);
    }
}

function createTimeline(name) {
    graphCount++;
    var container = document.getElementById("charts");
    var node = document.createElement("div");
    container.appendChild(node);
    var label = document.createElement("div");
    node.appendChild(label);
    label.innerHTML = name;
    var canvas = document.createElement("canvas");
    node.appendChild(canvas);

    // http://smoothiecharts.org/builder/

    var chart = new SmoothieChart({
	millisPerPixel:69,
	interpolation:'step',
	grid:{fillStyle:'#ffffff',
	      strokeStyle:'#e8e8ff',
	      sharpLines:true,
	      millisPerLine:5000,
	      verticalSections:0},
	labels:{fillStyle:'#004080'},
	timestampFormatter:SmoothieChart.timeFormatter
	// minValue:0
    });

    var points = new TimeSeries();

    chart.addTimeSeries(points, { strokeStyle: '#ff8040', lineWidth: 2 });
    chart.streamTo(canvas, 1000);

    return function tick(x) {
	points.append(new Date().getTime(), x);
    }
}
  </script>
    <script src="mon.js"></script>
    <button onclick="gesture()">Go</button>
  </body>
</html>
`)
