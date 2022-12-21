# receipt-processor
Implementation of Receipt Processor web service. The service runs on port 8080.

The command to run the receipt-processor app (without container; assumes Go 1.17 installed) is  
`$ go run server.go`

To run the unit tests  
`$ go test --cover ./...`

A dockerfile with go-1.17 on alpine is also provided for convenience.

The commands to build and run the receipt-processor image and container are as follows  
`$ docker build --tag receipt-processor .`  
`$ docker run -d -p 8080:8080 receipt-processor`



