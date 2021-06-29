release:
	go build -a -installsuffix cgo -o status

local:
	go run *.go --bind-address=127.0.0.1:8080 --cluster-name=local
