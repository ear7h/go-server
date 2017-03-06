echo "killing other instances, compiling, and starting new server"
sudo pkill go-server
/usr/local/go/bin/go build go-server.go
sudo ./go-server & disown
echo "done"
