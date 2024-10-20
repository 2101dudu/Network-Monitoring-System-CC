package main

import (
	u "CC/utils"
	"fmt"
)

func main() {
	ack_test := u.New_ack_builder().
		Set_request_id(3).
		Has_ackowledged()
	fmt.Println(ack_test)
}
