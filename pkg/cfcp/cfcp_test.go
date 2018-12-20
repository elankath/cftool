package cfcp

import (
	"flag"
)

var passCode = flag.String("code", "", "SSH Pass Code")

// IsPalindrome reports whether s reads the same forward and backward.
// // (Our first attempt.)
// func TestDoCopy(t *testing.T) {
// 	copier := &Copier{
// 		AppName: "itco",
// 		AppGUID: "9775c4e8-1ea6-4080-ab45-059e7e640310",
// 		InstanceIndex: 0,
// 		Simulated: false,
// 		sshHost: "ssh.cf.sap.hana.ondemand.com",
// 		sshPort: 2222,
// 	}
// 	fmt.Printf("Value of code: %s", *passCode)
// 	err := copier.doCopy("README.md", "/", *passCode)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// }
