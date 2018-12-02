package remote

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jiro4989/sion/util"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	groupFile = "/etc/group"
	userFile  = "/etc/passwd"
)

type SSHConfig struct {
	Host            string
	Port            int
	PrivateKey      []byte
	SSHClientConfig *ssh.ClientConfig
}

func CreateConnection(config *SSHConfig) (*ssh.Client, error) {
	var (
		host         = config.Host
		port         = config.Port
		privateKey   = config.PrivateKey
		clientConfig = config.SSHClientConfig
	)

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	clientConfig.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}

	hostport := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", hostport, clientConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func WithOpenFile(conn *ssh.Client, targetPath string, fn func(*sftp.File) (interface{}, error)) (interface{}, error) {
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	defer sftp.Close()

	f, err := sftp.OpenFile(targetPath, os.O_RDWR|os.O_CREATE)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return fn(f)
}

func fetchHash(conn *ssh.Client, fp string) (map[string]string, error) {
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	defer sftp.Close()

	f, err := sftp.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return convertColonTableToHash(f)
}

func FetchUsers(conn *ssh.Client) (map[string]string, error) {
	return fetchHash(conn, userFile)
}

func FetchGroups(conn *ssh.Client) (map[string]string, error) {
	return fetchHash(conn, groupFile)
}

func convertColonTableToHash(r io.Reader) (map[string]string, error) {
	sc := bufio.NewScanner(r)
	m := make(map[string]string)
	for sc.Scan() {
		line := sc.Text()
		cols := strings.Split(line, ":")
		if len(cols) == 0 || strings.HasPrefix(strings.TrimSpace(cols[0]), "#") {
			continue
		}
		v, id := cols[0], cols[2]
		m[id] = v
	}
	return m, sc.Err()
}

func HasDiff(conn *ssh.Client, srcFilePath, dstFilePath, owner, group, mode string, users, groups map[string]string) (bool, error) {
	hasDiff, err := WithOpenFile(conn, dstFilePath, func(f *sftp.File) (interface{}, error) {
		srcFile, err := os.Open(srcFilePath)
		if err != nil {
			return false, err
		}
		defer srcFile.Close()

		// ファイルサイズで比較し、
		// 一致しないなら後続の判定をスキップしてコピーを実行
		stat, err := f.Stat()
		if err != nil {
			return false, err
		}
		srcStat, err := srcFile.Stat()
		if err != nil {
			return false, err
		}
		if stat.Size() != srcStat.Size() {
			return true, nil
		}

		// ファイル内容で比較し、
		// 一致しないなら後続の判定をスキップしてコピーを実行
		fb, err := GetFileBytes(f)
		if err != nil {
			return false, err
		}
		srcBytes, err := util.GetFileBytes(srcFile)
		if err != nil {
			return false, err
		}
		fmt.Println("seq:", string(fb) == string(srcBytes))
		if !util.EqualBytes(fb, srcBytes) {
			return true, nil
		}

		// 権限を比較し、
		// 一致しないなら後続の判定をスキップしてコピーを実行
		if m := fmt.Sprintf("%04o", stat.Mode()); m != mode {
			return true, nil
		}

		// 所有者を判定し、
		// 一致しないなら後続の判定をスキップしてコピーを実行
		uid := stat.Sys().(*sftp.FileStat).UID
		uname := users[fmt.Sprintf("%d", uid)]
		if uname != owner {
			return true, nil
		}

		// 所有グループを判定し、
		// 一致しないなら後続の判定をスキップしてコピーを実行
		gid := stat.Sys().(*sftp.FileStat).GID
		gname := groups[fmt.Sprintf("%d", gid)]
		if gname != group {
			return true, nil
		}

		return false, nil
	})
	if b, ok := hasDiff.(bool); ok {
		return b, err
	}
	return false, err
}

func GetFileBytes(f *sftp.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return util.ReadByte(f, stat.Size())
}
