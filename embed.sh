#!/bin/bash

echo 'package main' > embedded.go
echo >> embedded.go
echo -n 'var monhtml = []byte(`' >> embedded.go
grep -B 1000 -F script static/mon.html | head -n -1 >> embedded.go
echo '  <script>' >> embedded.go
cat static/mon.js >> embedded.go
echo '  </script>' >> embedded.go
grep -A 1000 -F script static/mon.html | tail -n +2 >> embedded.go
echo '`)' >> embedded.go
