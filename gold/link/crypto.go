package link

import (
	"net"
	"sync"
	"unsafe"

	tea "github.com/fumiama/gofastTEA"
)

// Me 是本机的抽象
type Me struct {
	// 本机私钥
	// 利用 Curve25519 生成
	// https://pkg.go.dev/golang.org/x/crypto/curve25519
	// https://www.zhihu.com/question/266758647
	privKey [32]byte
	// 本机虚拟 ip
	me net.IP
	// 本机 endpoint
	myend *net.UDPAddr
	// 本机活跃的所有连接
	connections map[string]*Link
	// 读写同步锁
	connmapmu sync.RWMutex
	// 本机监听的 endpoint
	myconn *net.UDPConn
	// 本机路由表
	router *Router
}

// NewMe 设置本机参数
func NewMe(privateKey *[32]byte, myIP string, myEndpoint string) (m Me) {
	m.privKey = *privateKey
	var err error
	m.myend, err = net.ResolveUDPAddr("udp", myEndpoint)
	if err != nil {
		panic(err)
	}
	m.me = net.ParseIP(myIP)
	m.myconn, err = m.listen()
	if err != nil {
		panic(err)
	}
	m.connections = make(map[string]*Link)
	m.router = &Router{
		list:  make([]*net.IPNet, 1, 16),
		table: make(map[string]*Link, 16),
	}
	m.router.SetDefault(nil)
	return
}

// Encode 使用 TEA 加密
func (l *Link) Encode(b []byte) (eb []byte) {
	if b == nil {
		return
	}
	if l.key == nil {
		eb = b
	} else {
		// 在此处填写加密逻辑，密钥是l.key，输入是b，输出是eb
		// 不用写return，直接赋值给eb即可
		eb = (*tea.TEA)(unsafe.Pointer(l.key)).Encrypt(b)
	}
	return
}

// Decode 使用 TEA 解密
func (l *Link) Decode(b []byte) (db []byte) {
	if b == nil {
		return
	}
	if l.key == nil {
		db = b
	} else {
		// 在此处填写解密逻辑，密钥是l.key，输入是b，输出是db
		// 不用写return，直接赋值给db即可
		db = (*tea.TEA)(unsafe.Pointer(l.key)).Decrypt(b)
	}
	return
}
