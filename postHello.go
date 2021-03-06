package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
)

type postHelloInput struct {
    Name string
}

type postHelloOutput struct {
    Result string
}


func main() {
    handler := rest.ResourceHandler{}
    handler.SetRoutes(
        &rest.Route{"POST", "/hello", func(w rest.ResponseWriter, req *rest.Request) {
            input := postHelloInput{}
            err := req.DecodeJsonPayload(&input)

            if err != nil {
                rest.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if input.Name == "" {
                rest.Error(w, "Name is required", 400)
                return
            }

            log.Printf("%#v", input)

            w.WriteJson(&postHelloOutput{
            "Hello, " + input.Name,
         })
        }},
    )
    log.Printf("Server started")
    http.ListenAndServe(":9999", &handler)
}

/*
    Terminal

    $http -v POST localhost:9999/hello "Content-Type:application/json; charset=UTF-8" Name="Your Name"
    
*/
