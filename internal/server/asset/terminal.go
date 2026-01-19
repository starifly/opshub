package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	assetbiz "github.com/ydcloud-dy/opshub/internal/biz/asset"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HostInfo 主机信息
type HostInfo struct {
	ID           uint
	IP           string
	Port         int
	SSHUser      string
	CredentialID uint
}

// CredentialInfo 凭证信息
type CredentialInfo struct {
	ID         uint
	Type       string
	Username   string
	Password   string
	PrivateKey string
}

// TerminalSession 终端会话
type TerminalSession struct {
	ID         string
	HostID     uint
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
	StdinPipe  io.WriteCloser
	StdoutPipe  io.Reader
	StderrPipe  io.Reader
	CreatedAt   time.Time
}

// TerminalManager 终端管理器
type TerminalManager struct {
	sessions   map[string]*TerminalSession
	mu         sync.RWMutex
	hostUseCase *assetbiz.HostUseCase
}

// NewTerminalManager 创建终端管理器
func NewTerminalManager(hostUseCase *assetbiz.HostUseCase) *TerminalManager {
	return &TerminalManager{
		sessions:   make(map[string]*TerminalSession),
		hostUseCase: hostUseCase,
	}
}

// CreateSession 创建SSH会话
func (tm *TerminalManager) CreateSession(ctx context.Context, hostID uint, cols, rows uint16) (*TerminalSession, error) {
	// 获取主机信息
	hostVO, err := tm.hostUseCase.GetByID(ctx, hostID)
	if err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %w", err)
	}

	// 获取凭证（需要解密后的凭证）
	var credential *assetbiz.Credential
	if hostVO.CredentialID > 0 {
		credentialRepo := tm.hostUseCase.GetCredentialRepo()
		credential, err = credentialRepo.GetByIDDecrypted(ctx, hostVO.CredentialID)
		if err != nil {
			return nil, fmt.Errorf("获取凭证信息失败: %w", err)
		}
	} else {
		return nil, fmt.Errorf("主机未配置凭证")
	}

	// 解析私钥
	var signer ssh.Signer
	var authMethod ssh.AuthMethod

	if credential.Type == "key" {
		// 使用私钥认证
		if credential.Password != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(credential.PrivateKey), []byte(credential.Password))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(credential.PrivateKey))
		}
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		// 使用密码认证
		authMethod = ssh.Password(credential.Password)
	}

	// SSH配置
	config := &ssh.ClientConfig{
		User:            hostVO.SSHUser,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	address := fmt.Sprintf("%s:%d", hostVO.IP, hostVO.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %w", err)
	}

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("创建SSH会话失败: %w", err)
	}

	// 设置终端模式 - 简单可靠配置
	modes := ssh.TerminalModes{
		ssh.ECHO:          1, // SSH 服务器处理回显
		ssh.ICRNL:         1, // 输入 CR 转 NL
		ssh.ONLCR:         1, // 输出 NL 转 CR-NL
		ssh.ISIG:          1, // 信号处理
		ssh.IUTF8:         1, // UTF-8
	}

	if err := session.RequestPty("xterm-256color", int(cols), int(rows), modes); err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("请求伪终端失败: %w", err)
	}

	// 获取管道
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取stdin管道失败: %w", err)
	}

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取stdout管道失败: %w", err)
	}

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取stderr管道失败: %w", err)
	}

	// 启动shell
	if err := session.Shell(); err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("启动shell失败: %w", err)
	}

	// 创建会话对象
	terminalSession := &TerminalSession{
		ID:         fmt.Sprintf("%d-%d", hostID, time.Now().Unix()),
		HostID:     hostID,
		SSHClient:  client,
		SSHSession: session,
		StdinPipe:  stdinPipe,
		StdoutPipe: stdoutPipe,
		StderrPipe: stderrPipe,
		CreatedAt:  time.Now(),
	}

	// 保存会话
	tm.mu.Lock()
	tm.sessions[terminalSession.ID] = terminalSession
	tm.mu.Unlock()

	return terminalSession, nil
}

// GetSession 获取会话
func (tm *TerminalManager) GetSession(sessionID string) (*TerminalSession, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	session, ok := tm.sessions[sessionID]
	return session, ok
}

// CloseSession 关闭会话
func (tm *TerminalManager) CloseSession(sessionID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	session, ok := tm.sessions[sessionID]
	if !ok {
		return fmt.Errorf("会话不存在")
	}

	if session.SSHSession != nil {
		session.SSHSession.Close()
	}
	if session.SSHClient != nil {
		session.SSHClient.Close()
	}

	delete(tm.sessions, sessionID)
	return nil
}

// HandleSSHConnection 处理SSH WebSocket连接
func (s *HTTPServer) HandleSSHConnection(c *gin.Context) {
	fmt.Printf("=== WebSocket连接请求 ===\n")
	fmt.Printf("URL: %s\n", c.Request.URL.String())
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("Headers: %v\n", c.Request.Header)

	hostIdStr := c.Param("id")
	hostId, err := strconv.Atoi(hostIdStr)
	if err != nil {
		fmt.Printf("无效的主机ID: %s, error: %v\n", hostIdStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的主机ID"})
		return
	}
	fmt.Printf("主机ID: %d\n", hostId)

	// 从URL参数读取终端尺寸，默认80x24
	colsStr := c.DefaultQuery("cols", "80")
	rowsStr := c.DefaultQuery("rows", "24")
	cols, err := strconv.ParseUint(colsStr, 10, 16)
	if err != nil || cols == 0 {
		cols = 80
	}
	rows, err := strconv.ParseUint(rowsStr, 10, 16)
	if err != nil || rows == 0 {
		rows = 24
	}
	fmt.Printf("终端尺寸: %dx%d (cols=%s, rows=%s)\n", cols, rows, colsStr, rowsStr)

	// 升级到WebSocket
	fmt.Printf("开始升级到WebSocket...\n")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("WebSocket升级失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket升级失败"})
		return
	}
	defer conn.Close()
	fmt.Printf("WebSocket升级成功!\n")

	// 创建SSH会话
	fmt.Printf("开始创建SSH会话...\n")
	session, err := s.terminalManager.CreateSession(c.Request.Context(), uint(hostId), uint16(cols), uint16(rows))
	if err != nil {
		fmt.Printf("SSH会话创建失败: %v\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("连接失败: %s\r\n", err.Error())))
		return
	}
	defer s.terminalManager.CloseSession(session.ID)
	fmt.Printf("SSH会话创建成功! Session ID: %s\n", session.ID)

	// 启动goroutine从SSH读取输出并发送到WebSocket
	var wg sync.WaitGroup
	wg.Add(2)

	// 读取stdout
	go func() {
		defer wg.Done()
		fmt.Printf("开始读取SSH stdout...\n")
		buf := make([]byte, 1024)
		for {
			n, err := session.StdoutPipe.Read(buf)
			if n > 0 {
				fmt.Printf("从SSH stdout读取 %d 字节\n", n)
				// 使用二进制消息以保留原始字节（包括CR/LF控制字符）
				conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if err != nil {
				fmt.Printf("SSH stdout读取结束: %v\n", err)
				return
			}
		}
	}()

	// 读取stderr
	go func() {
		defer wg.Done()
		fmt.Printf("开始读取SSH stderr...\n")
		buf := make([]byte, 1024)
		for {
			n, err := session.StderrPipe.Read(buf)
			if n > 0 {
				fmt.Printf("从SSH stderr读取 %d 字节\n", n)
				// 使用二进制消息以保留原始字节（包括CR/LF控制字符）
				conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if err != nil {
				fmt.Printf("SSH stderr读取结束: %v\n", err)
				return
			}
		}
	}()

	fmt.Printf("开始处理WebSocket消息...\n")
	// 处理来自WebSocket的消息并发送到SSH
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WebSocket读取结束: %v\n", err)
			break
		}

		if messageType == websocket.TextMessage {
			// 尝试解析JSON消息（可能是resize命令）
			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err == nil {
				if msgType, ok := msg["type"].(string); ok && msgType == "resize" {
					cols, colsOk := msg["cols"].(float64)
					rows, rowsOk := msg["rows"].(float64)
					if colsOk && rowsOk {
						fmt.Printf("收到resize请求: %dx%d\n", int(cols), int(rows))
						// 调整SSH会话窗口大小
						if err := session.SSHSession.WindowChange(int(rows), int(cols)); err != nil {
							fmt.Printf("调整窗口大小失败: %v\n", err)
						} else {
							fmt.Printf("窗口大小调整成功\n")
						}
						continue
					}
				}
			}
			// 如果不是resize命令，当作普通输入发送到SSH
			fmt.Printf("从WebSocket收到文本 %d 字节，发送到SSH stdin\n", len(data))
			session.StdinPipe.Write(data)
		} else if messageType == websocket.BinaryMessage {
			fmt.Printf("从WebSocket收到二进制 %d 字节，发送到SSH stdin\n", len(data))
			session.StdinPipe.Write(data)
		}
	}

	fmt.Printf("等待所有goroutine结束...\n")
	wg.Wait()
	fmt.Printf("会话结束\n")
}

// ResizeTerminal 调整终端大小
func (s *HTTPServer) ResizeTerminal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "终端大小调整功能待实现"})
}
