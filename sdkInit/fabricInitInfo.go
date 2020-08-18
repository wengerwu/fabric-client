package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type ClientConfig struct {
	Clients []*Client `yaml:"clients"`
}

type Client struct {
	Org               Org `yaml:"org"`
	SDK               *fabsdk.FabricSDK
	ResmgmtClient     *resmgmt.Client
	MSPClient         *mspclient.Client
	ChannelClients    map[string]*channel.Client
	LedgerClients    map[string]*ledger.Client
	SDKConfigPath     string `yaml:"sdkConfigPath"`
	ChannelConfigPath string `yaml:"channelConfigPath"` // 通道配置路径
}

type Org struct {
	OrgName        string `yaml:"orgName"`        // 组织名称
	OrgAdmin       string `yaml:"orgAdmin"`       // 组织管理员
	OrdererOrgName string `yaml:"ordererOrgName"` // 排序组织名称
	OrgMspID       string `yaml:"orgMspID"`       // 组织成员关系服务提供者标识
}

type CCRequest struct {
	ChannelID string // 通道ID
	OrgName   string // 组织名称

	ChaincodeID      string   //链码ID
	ChaincodeVersion string   //链码版本
	ChaincodePath    string   //链码路径
	Args             []string //链码参数

	Timestamp int64  //时间戳
	Sign      string //签名
}

type ChannelClientRequest struct {
	ChannelID string // 通道ID
	OrgName   string // 组织名称
	UserName  string //用户名称
}
