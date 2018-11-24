package remote

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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

func WithOpenFile(conn *ssh.Client, targetPath string, fn func(*sftp.File) error) error {
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftp.Close()

	f, err := sftp.Open(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(f)
}

func FindUserName(conn *ssh.Client, uid string) (string, error) {
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return "", err
	}
	defer sftp.Close()

	f, err := sftp.Open(userFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return lookupId(f, uid)
}

// TODO ここ完全に使いまわしになっていてダサイ
func FindGroupName(conn *ssh.Client, gid string) (string, error) {
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return "", err
	}
	defer sftp.Close()

	f, err := sftp.Open(groupFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return lookupId(f, gid)
}

func lookupId(r io.Reader, id string) (string, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		cols := strings.Split(line, ":")
		if len(cols) == 0 || strings.HasPrefix(strings.TrimSpace(cols[0]), "#") {
			continue
		}
		userName, userId := cols[0], cols[2]
		if userId == id {
			return userName, nil
		}
	}
	return "", sc.Err()
}
