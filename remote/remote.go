package remote

import (
	"fmt"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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
