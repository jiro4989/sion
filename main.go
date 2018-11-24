package main

import (
	"github.com/jiro4989/sion/cli/command"
)

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
