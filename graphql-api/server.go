package main

import (
    "log"
    "net/http"

    "github.com/graphql-go/handler"
)

func main() {
    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
    })

    http.Handle("/graphql", h)
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

