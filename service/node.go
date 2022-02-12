package service

import (
	"os"

	"google.golang.org/grpc"

	"github.com/binacsgo/inject"
)

// NodeService the node service
type NodeService interface{}

// NodeServiceImpl Web crypto service implement
type NodeServiceImpl struct {
	Conn     *grpc.ClientConn `inject-name:"Conn"`
	Cos      CosClient        `inject-name:"Cos"`
	Crypto   CryptoClient     `inject-name:"Crypto"`
	Pastebin PastebinClient   `inject-name:"Pastebin"`
	TinyURL  TinyURLClient    `inject-name:"TinyURL"`
	User     UserClient       `inject-name:"User"`
}

// AfterInject do inject
func (impl *NodeServiceImpl) AfterInject() error {
	return nil
}

func InitService(conn *grpc.ClientConn) *NodeServiceImpl {
	oriStdout := os.Stdout
	f, _ := os.OpenFile("/dev/null", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	os.Stdout = f
	defer func() {
		os.Stdout = oriStdout
	}()

	nodeSvc := NodeServiceImpl{}

	// Services
	inject.Regist("Conn", conn)
	inject.Regist("Cos", &CosClientImpl{})
	inject.Regist("Crypto", &CryptoClientImpl{})
	inject.Regist("Pastebin", &PastebinClientImpl{})
	inject.Regist("TinyURL", &TinyURLClientImpl{})
	inject.Regist("User", &UserClientImpl{})
	inject.Regist("Node", &nodeSvc)

	inject.DoInject()
	return &nodeSvc
}
