package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jiro4989/sion/remote"
	"github.com/jiro4989/sion/util"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var cpCommand = &cobra.Command{
	Use:   "cp",
	Short: "cp copies file to remote server",
	Long:  "cp copies file to remote server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			// TODO output help
			fmt.Println("argsたりない")
			return
		}
		srcFilePath, dstFilePath := args[0], args[1]
		fmt.Println(srcFilePath)

		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			panic(err)
		}
		group, err := cmd.Flags().GetString("group")
		if err != nil {
			panic(err)
		}
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			panic(err)
		}

		fmt.Println("owner:", owner)
		fmt.Println("group:", group)
		fmt.Println("mode:", mode)

		user := "ec2-user"

		hostByte, err := ioutil.ReadFile("/home/jiro/host.txt")
		// TODO 一時的な対応
		if err != nil {
			panic(err)
		}
		host := string(hostByte)
		host = strings.TrimSpace(host)

		pk, err := ioutil.ReadFile("/home/jiro/.ssh/sandbox.pem")
		// TODO 一時的な対応
		if err != nil {
			panic(err)
		}

		config := &remote.SSHConfig{
			Host:       host,
			Port:       22,
			PrivateKey: pk,
			SSHClientConfig: &ssh.ClientConfig{
				User:            user,
				Timeout:         5 * time.Second,
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}
		conn, err := remote.CreateConnection(config)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		if err := remote.WithOpenFile(conn, dstFilePath, func(f *sftp.File) error {
			stat, err := f.Stat()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println("Stat:", stat)
			fmt.Printf("Mode: %04o\n", stat.Mode())
			fmt.Println("Size:", stat.Size())

			uid := stat.Sys().(*sftp.FileStat).UID
			fmt.Println("Sys uid:", uid)

			gid := stat.Sys().(*sftp.FileStat).GID
			fmt.Println("Sys gid:", gid)

			uname, err := remote.FindUserName(conn, fmt.Sprintf("%d", uid))
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println("Sys username:", uname)

			gname, err := remote.FindGroupName(conn, fmt.Sprintf("%d", gid))
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println("Sys groupname:", gname)

			var b = make([]byte, stat.Size())
			n, err := f.Read(b)
			if n == 0 {
				return errors.New("ファイル読み込みに失敗")
			}
			if err != nil {
				return err
			}

			of, err := os.Open("hello.txt")
			if err != nil {
				return err
			}
			defer of.Close()

			oStat, err := of.Stat()
			if err != nil {
				return err
			}

			var bb = make([]byte, oStat.Size())
			n, err = of.Read(bb)
			if n == 0 {
				return errors.New("ファイル読み込みに失敗")
			}
			if err != nil {
				return err
			}

			fmt.Println("remote byte :", b)
			fmt.Println("local byte :", bb)
			fmt.Println("sameBytes?: ", util.EqualBytes(b, bb))
			fmt.Println("samePerm?: ", stat.Mode() == oStat.Mode())

			// Create the destination file
			dstFile, err := os.Create("tmpfile.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer dstFile.Close()

			// Copy the file
			f.WriteTo(dstFile)

			return nil
		}); err != nil {
			panic(err)
		}
	},
}

func init() {
	CommandCommand.AddCommand(cpCommand)
	cpCommand.Flags().StringP("owner", "o", "", "owner of remote file")
	cpCommand.Flags().StringP("group", "g", "", "group of remote file")
	cpCommand.Flags().StringP("mode", "m", "", "mode of remote file")
}
