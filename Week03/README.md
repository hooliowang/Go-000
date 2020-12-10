Execute the `go run main.go` command and start three servers as below.

Server1: test server 1, listen 8080 port\n
Server2: test server 2, listen 8081 port\n
Server3: server for shutdown the assigned server, listen 8082 port

Execute `curl http://127.0.0.1:8082/close?servername=s1` and it will only shutdown Server1.
Execute `curl http://127.0.0.1:8082/close?servername=s2` and it will only shutdown Server2.
Execute `curl http://127.0.0.1:8082/close?servername=all` and it will shutdown all the servers including itself.

Execute "CTRL+C" command in the terminal and it will shutdown all the three servers.