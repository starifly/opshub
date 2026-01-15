package sshclient

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client SSH客户端
type Client struct {
	client *ssh.Client
}

// NewClient 创建SSH客户端
func NewClient(host string, port int, username, password string, privateKey []byte, passphrase string) (*Client, error) {
	var authMethods []ssh.AuthMethod

	// 优先使用私钥认证
	if len(privateKey) > 0 {
		var signer ssh.Signer
		var err error

		if passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(privateKey)
		}

		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 密码认证
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("至少需要一种认证方式")
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应该验证主机密钥
		Timeout:         10 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %w", err)
	}

	return &Client{client: client}, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Execute 执行命令
func (c *Client) Execute(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建session失败: %w", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("命令执行失败: %s, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ExecuteWithTimeout 执行命令（带超时）
func (c *Client) ExecuteWithTimeout(command string, timeout time.Duration) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建session失败: %w", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// 使用通道实现超时
	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-time.After(timeout):
		// 超时，尝试关闭session
		session.Signal(ssh.SIGTERM)
		session.Close()
		return "", fmt.Errorf("命令执行超时")
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("命令执行失败: %s, stderr: %s", err, stderr.String())
		}
	}

	return stdout.String(), nil
}

// TestConnection 测试连接
func (c *Client) TestConnection() error {
	_, err := c.Execute("echo ok")
	return err
}
