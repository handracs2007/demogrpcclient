package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "github.com/handracs2007/demogrpcclient/rpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "log"
)

func main() {
    // Load the client certificate and its key
    clientCert, err := tls.LoadX509KeyPair("client.pem", "client.key")
    if err != nil {
        log.Fatalf("Failed to load client certificate and key. %s.", err)
    }

    // Load the CA certificate
    trustedCert, err := ioutil.ReadFile("cacert.pem")
    if err != nil {
        log.Fatalf("Failed to load trusted certificate. %s.", err)
    }

    // Put the CA certificate to certificate pool
    certPool := x509.NewCertPool()
    if !certPool.AppendCertsFromPEM(trustedCert) {
        log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
    }

    // Create the TLS configuration
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{clientCert},
        RootCAs:      certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
    }

    // Create a new TLS credentials based on the TLS configuration
    cred := credentials.NewTLS(tlsConfig)

    // Dial the gRPC server with the given credentials
    conn, err := grpc.Dial("localhost:8443", grpc.WithTransportCredentials(cred))
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        err = conn.Close()
        if err != nil {
            log.Printf("Unable to close gRPC channel. %s.", err)
        }
    }()

    // Create the request data
    request := &rpc.HelloRequest{
        Name: "Handra",
        Age:  35,
    }

    // Create the gRPC client
    client := rpc.NewDemoServiceClient(conn)
    response, err := client.SayHello(context.Background(), request)
    if err != nil {
        log.Fatalf("Failed to receive response. %s.", err)
    }

    // Print out response from server
    fmt.Println(response.Response)
}
