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
		fmt.Println(owner, group, mode)

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
			// gotoでスキップされる間に初めて宣言される変数が存在するとコンパイ
			// ルエラーになるため、不本意ながらも先頭にまとめて変数宣言
			var (
				uid      uint32
				gid      uint32
				uname    string
				gname    string
				fb       []byte
				srcBytes []byte
				n        int
			)

			srcFile, err := os.Open(srcFilePath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			stat, err := f.Stat()
			if err != nil {
				return err
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
			fmt.Println("seq:", string(fb) == string(srcBytes))
			if !util.EqualBytes(fb, srcBytes) {
				goto execopy
			}

			// 権限を比較し、
			// 一致しないなら後続の判定をスキップしてコピーを実行
			if m := fmt.Sprintf("%04o", stat.Mode()); m != mode {
				goto execopy
			}

			/*
				INFO: ここからはSFTPでpasswd, groupファイルを取得するための通信が発生する
				よってパフォーマンスに大きく影響を与える
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

			fmt.Println("Skipping...")
			goto skipcopy

		execopy:
			fmt.Println("copying...")

			err = f.Truncate(0)
			if err != nil {
				return err
			}

			n, err = f.Write(srcBytes)
			if err != nil {
				return err
			}
			if n == 0 {
				fmt.Println(n)
				return errors.New("ファイル書き込みに失敗しました")
			}

			if err := f.Chmod(stat.Mode()); err != nil {
				return err
			}
			if err := f.Chown(1000, 1000); err != nil {
				return err
			}

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
