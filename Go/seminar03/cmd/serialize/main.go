package main

import (
	"log"
	"os"

	v1 "example.com/seminar03/internal/api/proto/service/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	res := v1.EchoResponse{
		Message: "Hello",
		Time:    "12:00",
		Ip:      "127.0.0.1",
	}

	bytes, err := proto.Marshal(&res)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("out", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
