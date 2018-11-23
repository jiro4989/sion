package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	user = flag.String("u", "", "user")
	port = flag.Int("P", 22, "port")
	// password = flag.String("p", "", "password")
)

func run() int {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return 2
	}

	privateKey, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(*password),
			ssh.PublicKeys(signer),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostport := fmt.Sprintf("%s:%d", flag.Arg(1), *port)
	conn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect %v: %v", hostport, err)
		return 1
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open new session: %v", err)
		return 1
	}
	defer session.Close()

	go func() {
		time.Sleep(5 * time.Second)
		conn.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	err = session.Run(strings.Join(flag.Args()[2:], " "))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		if ee, ok := err.(*ssh.ExitError); ok {
			return ee.ExitStatus()
		}
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}

func CopyFile() {
	// ローカルのファイルとリモート先のファイルを比較
	// 差分があった場合だけコピーを実行
}

// func Diff() bool {
// 	// タイムスタンプのdiff
// 	if !equalsTimestamp(inf, outf) {
// 		return true
// 	}
// 	// 権限のdiff
// 	if !equalsPermission(inf, outf) {
// 		return true
// 	}
// 	// 所有者のdiff
// 	if !equalsOwner(inf, outf) {
// 		return true
// 	}
// 	// 所有グループのdiff
// 	if !equalsGroup(inf, outf) {
// 		return true
// 	}
// 	// ファイルのバイトサイズのdiff
// 	if !equalsBytes(inf, outf) {
// 		return true
// 	}
// 	// ファイル内容のdiff
// 	if !equalsContents(inf, outf) {
// 		return true
// 	}
// 	return false
// }
