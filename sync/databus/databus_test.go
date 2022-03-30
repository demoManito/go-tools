package databus

// TODO: 测试
// func TestName(t *testing.T) {
// 	mm, _ := proto.Marshal(&Header{Metadata: map[string]string{"d": strconv.FormatInt(time.Now().Add(1*time.Minute).UnixNano()/1e6, 10)}})
// 	fmt.Println(string(mm))
//
// 	msg := &MessagePB{}
// 	err := proto.Unmarshal(mm, msg)
// 	if err != nil {
// 		fmt.Println("err", err)
// 	}
// 	fmt.Println("test", msg.Key)
// }
