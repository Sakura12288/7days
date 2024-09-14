package main

import (
	"7days/Cache/pb"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
)

func main() {
	test := &pb.Student{
		Name:   "geektutu",
		Male:   true,
		Scores: []int32{98, 85, 88},
	}
	data, err := proto.Marshal(test)
	fmt.Println(string(data))
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	newTest := &pb.Student{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if test.GetName() != newTest.GetName() {
		log.Fatalf("data mismatch %q != %q", test.GetName(), newTest.GetName())
	}
	fmt.Printf("%T", newTest)
}
