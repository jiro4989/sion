package command

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
			var (
				uid      uint32
				gid      uint32
				uname    string
				gname    string
				fb       []byte
				srcBytes []byte
			)

			srcFile, err := os.Open(srcFilePath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			// ファイルサイズで比較し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			stat, err := f.Stat()
			if err != nil {
				return err
			}
			srcStat, err := srcFile.Stat()
			if err != nil {
				return err
			}
			if stat.Size() != srcStat.Size() {
				goto execopy
			}

			// ファイル内容で比較し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			fb, err = getRemoteFileBytes(f)
			if err != nil {
				return err
			}
			srcBytes, err = getFileBytes(srcFile)
			if err != nil {
				return err
			}
			if !util.EqualBytes(fb, srcBytes) {
				goto execopy
			}

			// 権限を比較し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			if m := fmt.Sprintf("%04o", stat.Mode()); m != mode {
				goto execopy
			}

			/*
				INFO: ここからはSFTPでユーザ、グループファイルを取得することにな
				るので、ネットワーク遅延が速度に響く
			*/

			// 所有者を判定し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			uid = stat.Sys().(*sftp.FileStat).UID
			uname, err = remote.FindUserName(conn, fmt.Sprintf("%d", uid))
			if err != nil {
				return err
			}
			if uname != owner {
				goto execopy
			}

			// 所有グループを判定し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			gid = stat.Sys().(*sftp.FileStat).GID
			gname, err = remote.FindGroupName(conn, fmt.Sprintf("%d", gid))
			if err != nil {
				return err
			}
			if gname != group {
				goto execopy
			}

			goto skipcopy

		execopy:
			fmt.Println("copying...")

		skipcopy:
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

// TODO 抽象化
func getRemoteFileBytes(f *sftp.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return readByte(f, stat.Size())
}

// TODO 抽象化
func getFileBytes(f *os.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return readByte(f, stat.Size())
}

func readByte(f io.Reader, size int64) ([]byte, error) {
	var b = make([]byte, size)
	n, err := f.Read(b)
	if n == 0 {
		return nil, errors.New("ファイル読み込みに失敗")
	}
	if err != nil {
		return nil, err
	}

	return b, nil
}
