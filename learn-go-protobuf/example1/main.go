package main

import (
	"fmt"

	pb "com.example.test/person"
	"google.golang.org/protobuf/proto"
)

func main() {
	p := pb.Person{
		Name: "jone",
		Age:  18,
		Sex:  pb.Sex_MAN,
	}
	fmt.Println("p : ", p.String())
	fmt.Println("p : ", "name:", p.GetName(), "age:", p.GetAge(), "sex:", p.GetSex())

	// 使用 google.golang.org/protobuf/proto 序列化和反序列化
	bytes, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}
	var p2 pb.Person
	err = proto.Unmarshal(bytes, &p2)
	fmt.Println("p2 : ", "name:", p2.GetName(), "age:", p2.GetAge(), "sex:", p2.GetSex())
}
