package main

import (
	"flag"

	"github.com/jiro4989/sion/cli/command"
)

var (
	user = flag.String("u", "", "user")
	port = flag.Int("P", 22, "port")
	// password = flag.String("p", "", "password")
)

// func CreateConnection(pemPath, host string) *ssh.Client {
// 	privateKey, err := ioutil.ReadFile(pemPath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	signer, err := ssh.ParsePrivateKey(privateKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	config := &ssh.ClientConfig{
// 		User: *user,
// 		Auth: []ssh.AuthMethod{
// 			// ssh.Password(*password),
// 			ssh.PublicKeys(signer),
// 		},
// 		Timeout:         5 * time.Second,
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}
//
// 	hostport := fmt.Sprintf("%s:%d", host, *port)
// 	conn, err := ssh.Dial("tcp", hostport, config)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "cannot connect %v: %v", hostport, err)
// 		return nil
// 	}
// 	return conn
// }

// func sshrun(pemPath, host string, cmds []string) int {
// 	conn := CreateConnection(pemPath, host)
// 	defer conn.Close()
//
// 	session, err := conn.NewSession()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "cannot open new session: %v", err)
// 		return 1
// 	}
// 	defer session.Close()
//
// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		conn.Close()
// 	}()
//
// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr
// 	session.Stdin = os.Stdin
// 	err = session.Run(strings.Join(cmds, " "))
// 	if err != nil {
// 		fmt.Fprint(os.Stderr, err)
// 		if ee, ok := err.(*ssh.ExitError); ok {
// 			return ee.ExitStatus()
// 		}
// 		return 1
// 	}
// 	return 0
// }
//
// func scp(pemPath, host string) int {
// 	conn := CreateConnection(pemPath, host)
// 	defer conn.Close()
//
// 	session, err := conn.NewSession()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "cannot open new session: %v", err)
// 		return 1
// 	}
// 	defer session.Close()
//
// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr
// 	// session.Stdin = os.Stdin
// 	go func() {
// 		w, _ := session.StdinPipe()
// 		defer w.Close()
// 		content := "123456789\n"
// 		//fmt.Fprintln(w, "D0755", 0, "testdir") // mkdir
// 		fmt.Fprintln(w, "C0644", len(content), "testfile1")
// 		fmt.Fprint(w, content)
// 		fmt.Fprint(w, "\x00") // transfer end with \x00
// 		fmt.Fprintln(w, "C0644", len(content), "testfile2")
// 		fmt.Fprint(w, content)
// 		fmt.Fprint(w, "\x00")
// 	}()
// 	if err := session.Run("/usr/bin/scp -tr ./"); err != nil {
// 		panic("Failed to run: " + err.Error())
// 	}
//
// 	return 0
// }

func main() {
	if err := command.RootCommand.Execute(); err != nil {
		panic(err)
	}
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
