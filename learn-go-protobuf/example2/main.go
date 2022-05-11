package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lxxxxxxxx/protobuf-go/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	// pm := NewPersonMessage()

	// _ = writeToFile("./person.bin", pm)

	// pm1 := &pb.PersonMessage{}
	// loadFromFile("./person.bin", pm1)

	// fmt.Println(pm1)

	// json := toJSON(pm)
	// fmt.Println(json)

	// pm2 := &pb.PersonMessage{}
	// _ = fromJSON(json, pm2)
	// fmt.Println(pm2)

	em := NewEnumMessage()
	fmt.Println(pb.Gender_name[int32(em.Gender)])
	fmt.Println(em.Gender.String())

	dm := NewDepartmentMessage()
	fmt.Println(dm)

}

func NewDepartmentMessage() *pb.Department {
	dm := &pb.Department{
		Id:   5544,
		Name: "销售部门",
		Employees: []*pb.EmployeeMessage{
			&pb.EmployeeMessage{
				Id:   11,
				Name: "Dave",
			},
			&pb.EmployeeMessage{
				Id:   12,
				Name: "Mike",
			},
		},
		ParentDepartment: &pb.Department{
			Id:   1,
			Name: "总公司",
		},
		ChildrenDepartment: nil,
	}
	return dm
}

func NewEnumMessage() *pb.EnumMessage {
	em := pb.EnumMessage{
		Id:     345,
		Gender: pb.Gender_FEMALE,
	}

	em.Gender = pb.Gender_MALE
	return &em
}

// func toJSON(pb proto.Message) string {
// 	marshaler := jsonpb.Marshaler{Indent: "  "}
// 	str, err := marshaler.MarshalToString(pb)
// 	if err != nil {
// 		log.Fatalln("转为json时错误")
// 		return ""
// 	}
// 	return str
// }

// func fromJSON(in string, pb proto.Message) error {
// 	err := jsonpb.UnmarshalString(in, pb)
// 	if err != nil {
// 		log.Fatalln("读取json时发生错误")
// 		return err
// 	}
// 	return nil
// }

func writeToFile(filename string, pb proto.Message) error {
	bytes, err := proto.Marshal(pb)
	if err != nil {
		log.Fatalln("无法序列化到bytes", err.Error())
		return err
	}

	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		log.Fatalln("无法写入到文件", err.Error())
		return err
	}
	log.Println("写入文件成功")
	return nil
}

func loadFromFile(filename string, pb proto.Message) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("无法读取文件")
		return err
	}
	if err := proto.Unmarshal(bytes, pb); err != nil {
		log.Fatalln("反序列化失败", err.Error())
		return err
	}
	log.Println("读取成功")
	return nil
}

func NewPersonMessage() *pb.Person {
	pm := pb.Person{
		Id:           1234,
		IsAdult:      true,
		Name:         "dave",
		LuckyNumbers: []int32{1, 2, 3, 4, 5},
	}

	fmt.Println(pm)
	pm.Name = "Nick"
	fmt.Println(pm)
	fmt.Printf("pm id:%d\n", pm.GetId())
	return &pm
}
