package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/short-loop/shortloop-go/shortloopgin"
	"net/http"
)

func TestEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	obj := map[string]string{
		"message": "test",
	}
	return obj, nil
}

type DecodeRequestFunc func(context.Context, *gin.Context) (interface{}, error)
type EncodeResponseFunc func(context.Context, *gin.Context, interface{}) error

func DecodeRequest(_ context.Context, g *gin.Context) (interface{}, error) {
	return g.Request.Body, nil
}

func EncodeJSONResponse(_ context.Context, c *gin.Context, response interface{}) error {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := http.StatusOK
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(response)
}

func NewHTTPHandler(ep endpoint.Endpoint, dec DecodeRequestFunc, enc EncodeResponseFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req, err := dec(ctx, c)
		if err != nil {
			fmt.Println("Error decoding request: ", err)
			return
		}
		resp, err := ep(ctx, req)
		if err != nil {
			fmt.Println("Error calling endpoint: ", err)
			return
		}
		// ctx not available here
		if err := enc(ctx, c, resp); err != nil {
			fmt.Println("Error encoding response: ", err)
			return
		}
	}
}

func main() {
	router := gin.Default()
	shortloopSdk, err := shortloopgin.Init(shortloopgin.Options{
		ShortloopEndpoint: "http://localhost:8080", // the shortloop url for your org. (Provided by ShortLoop team.)
		ApplicationName:   "go-kit-9",              // your application name here.
		Environment:       "stage",                 // for e.g stage or prod
		LoggingEnabled:    true,                    // enable logging
		LogLevel:          "INFO",
		// no auth key added since its critical
	})
	if err != nil {
		fmt.Println("Error initializing shortloopgin: ", err)
	} else {
		router.Use(shortloopSdk.Filter())
	}
	v1 := router.Group("/v1")
	v1.GET("/test", NewHTTPHandler(TestEndpoint, DecodeRequest, EncodeJSONResponse))

	v1.GET("/hello2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello2"})
	})

	router.Run(":8300")
}

//helloEndpoint := func(c context.Context, request interface{}) (interface{}, error) {
//
//	obj := map[string]string{
//		"message": "hello",
//	}
//	return obj, nil
//}

// Create the HTTP transport handler for the endpoint.
//helloHandler := httptransport.NewServer(
//	makeEndpoint(helloEndpoint),
//	decodeRequest,
//	encodeResponse,
//)

//router.GET("/hello", func(c *gin.Context) {
//	helloHandler.ServeHTTP(c.Writer, c.Request)
//})

// makeEndpoint creates an endpoint from a function that takes a context and a request,
// and returns a response and an error.
//func makeEndpoint(fn endpoint.Endpoint) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//		return fn(ctx, request)
//	}
//}

// decodeRequest decodes an HTTP request into a request object.
//func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	return r, nil
//}

// encodeResponse encodes a response object into an HTTP response.
//func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	return json.NewEncoder(w).Encode(response)
//}
