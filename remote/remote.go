package remote

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
