package service

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-platform/config"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type terminal struct{}

var Terminal terminal

type TerminalMessage struct {
	// 操作类型
	Operation string `json:"operation"`
	// 具体数据
	Data string `json:"data"`

	//终端的行数
	Rows uint16 `json:"rows"`

	//终端的列数
	Cols uint16 `json:"cols"`
}

type TerminalSession struct {
	// ws
	wsConn *websocket.Conn
	// 终端的大小
	sizeChan chan remotecommand.TerminalSize
	// flag的chan
	doneChan chan struct{}
}

// websocket的handler 跟gin不是同一个接口
func (t *terminal) WsHander(w http.ResponseWriter, r *http.Request) {

	// 初始化RESTClient，这里不使用用ClientSet，需要重新初始化
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("加载kubeconfig失败，" + err.Error())
	}

	// 解析入参
	if err := r.ParseForm(); err != nil {
		logger.Error("参数解析失败，" + err.Error())
	}

	namespace := r.Form.Get("namespace")
	podname := r.Form.Get("podname")
	containername := r.Form.Get("container_name")

	// 创建websocket
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("创建pty失败，" + err.Error())
		return
	}
	defer func() {
		logger.Info("关闭pty")
		pty.Close()
	}()

	//拼接POST请求
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").Name(podname).Namespace(namespace).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containername,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	// 创建SPDY
	exector, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		logger.Error("建立SPDY连接失败, " + err.Error())
		return
	}

	// 通过websocket来处理Stream数据
	err = exector.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})

	if err != nil {
		logger.Error("执行 pod 命令失败, " + err.Error())
		pty.Write([]byte("执行 pod 命令失败, " + err.Error()))
		// 关闭session
		pty.Done()
	}

}

// 初始化upgrader
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader

}()

// http升级协议到websocket
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

// 关闭TerminalSession
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// 关闭websocket
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// 调整窗口大小
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}

}

// 接收从web端的输入内容
// 返回值表示读取了数据个数
func (t *TerminalSession) Read(p []byte) (int, error) {

	// 从ws读取数据
	_, readmsg, err := t.wsConn.ReadMessage()
	if err != nil {
		logger.Error("read message error,", err)
		// const END_OF_TRANSMISSION = "\u0004"
		// return copy(p, END_OF_TRANSMISSION), err
		return 0, err
	}

	// 解析数据
	var msg TerminalMessage
	if err := json.Unmarshal(readmsg, &msg); err != nil {
		logger.Error("unmarshal msg failed,", err)
		return 0, err
	}

	// 判断操作类型
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		logger.Error("unknown message type," + msg.Operation)
		return 0, fmt.Errorf("unknown message type: %s\n" + msg.Operation)

	}

}

// 输出内容到web端
func (t *TerminalSession) Write(p []byte) (int, error) {

	//  序列化数据
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})

	if err != nil {
		logger.Error("marshal message failed,", err)
		return 0, err
	}

	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {

		logger.Error("write message failed,", err)
		return 0, err
	}

	return len(p), nil
}
