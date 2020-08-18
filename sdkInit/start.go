package sdkInit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"gopkg.in/yaml.v2"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

var goPath = os.Getenv("GOPATH")

func InitClientMap() (map[string]*Client, error) {
	clientConfig, err := readClientConfig("config/client-config.yaml")
	if err != nil {
		return nil, err
	}

	clientMap := make(map[string]*Client, len(clientConfig.Clients))
	for _, client := range clientConfig.Clients {
		err = initClient(client)
		if err != nil {
			return nil, err
		}

		clientMap[client.Org.OrgName] = client
	}
	return clientMap, nil
}

func initClient(client *Client) error {
	sdk, err := fabsdk.New(config.FromFile(client.SDKConfigPath))
	if err != nil {
		return fmt.Errorf("初始化【%s】组织的FabricSDK失败:%v", client.Org.OrgName, err)
	}

	clientProvider := sdk.Context(fabsdk.WithOrg(client.Org.OrgName), fabsdk.WithUser(client.Org.OrgAdmin))
	if clientProvider == nil {
		return fmt.Errorf("创建【%s】组织的资源管理客户端Context失败", client.Org.OrgName)
	}

	resmgmtClient, err := resmgmt.New(clientProvider)
	if err != nil {
		return fmt.Errorf("创建【%s】组织的通道管理客户端失败: %v", client.Org.OrgName, err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(client.Org.OrgName))
	if err != nil {
		return fmt.Errorf("创建【%s】组织的OrgMSP客户端实例失败: %v", client.Org.OrgName, err)
	}

	client.SDK = sdk
	client.ResmgmtClient = resmgmtClient
	client.MSPClient = mspClient
	client.ChannelClients = make(map[string]*channel.Client)
	client.LedgerClients = make(map[string]*ledger.Client)
	return nil
}

func CloseClientMap(clientMap map[string]*Client) {
	for _, client := range clientMap {
		client.SDK.Close()
	}
}

func (client *Client) CreateChannel(channelID string) error {
	var err error
	adminIdentity, err := client.MSPClient.GetSigningIdentity(client.Org.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取【%s】签名标识失败: %v", client.Org.OrgAdmin, err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: client.ChannelConfigPath,
		SigningIdentities: []mspctx.SigningIdentity{adminIdentity}}
	_, err = client.ResmgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(client.Org.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败: %v", err)
	}

	fmt.Println("通道已成功创建")
	return nil
}

func (client *Client) JoinChannel(channelID string) error {
	err := client.ResmgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(client.Org.OrdererOrgName))
	if err != nil {
		fmt.Errorf("Peers 加入通道失败: %v", err)
		return err
	}

	fmt.Println("Peers 已成功加入通道")
	return nil
}

func (client *Client) InstallCC(ccRequest *CCRequest) error {
	fmt.Println("开始安装链码......")
	ccPkg, err := gopackager.NewCCPackage(ccRequest.ChaincodePath, goPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v", err)
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    ccRequest.ChaincodeID,
		Path:    ccRequest.ChaincodePath,
		Version: ccRequest.ChaincodeVersion,
		Package: ccPkg,
	}
	_, err = client.ResmgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("安装链码失败: %v", err)
	}

	fmt.Println("指定的链码安装成功")
	return nil
}

func (client *Client) InstantiateCC(ccRequest *CCRequest) error {
	fmt.Println("开始实例化链码......")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{client.Org.OrgMspID})
	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    ccRequest.ChaincodeID,
		Path:    ccRequest.ChaincodePath,
		Version: ccRequest.ChaincodeVersion,
		Args:    ToBytesArgs(ccRequest.Args),
		Policy:  ccPolicy,
	}
	_, err := client.ResmgmtClient.InstantiateCC(ccRequest.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("实例化链码失败: %v", err)
	}

	fmt.Println("链码实例化成功")
	return nil
}

func (client *Client) UpgradeCC(ccRequest *CCRequest) error {
	fmt.Println("开始升级链码......")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{client.Org.OrgMspID})
	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    ccRequest.ChaincodeID,
		Path:    ccRequest.ChaincodePath,
		Version: ccRequest.ChaincodeVersion,
		Args:    ToBytesArgs(ccRequest.Args),
		Policy:  ccPolicy,
	}
	_, err := client.ResmgmtClient.UpgradeCC(ccRequest.ChannelID, upgradeCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("升级链码失败: %v", err)
	}

	fmt.Println("链码升级成功")
	return nil
}

func (client *Client) NewChannelClient(channelClientRequest *ChannelClientRequest) (*channel.Client, error) {
	clientChannelContext := client.SDK.ChannelContext(channelClientRequest.ChannelID, fabsdk.WithUser(channelClientRequest.UserName), fabsdk.WithOrg(channelClientRequest.OrgName))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")
	return channelClient, nil
}

func (client *Client) NewLedgerClient(channelClientRequest *ChannelClientRequest) (*ledger.Client, error) {
	clientLedgerContext := client.SDK.ChannelContext(channelClientRequest.ChannelID, fabsdk.WithUser(channelClientRequest.UserName), fabsdk.WithOrg(channelClientRequest.OrgName))
	ledgerClient, err := ledger.New(clientLedgerContext)
	if err != nil {
		return nil, fmt.Errorf("创建账本客户端失败: %v", err)
	}

	fmt.Println("账本客户端创建成功.")
	return ledgerClient, nil
}

func ToBytesArgs(args []string) [][]byte {
	len := len(args)
	bytesArgs := make([][]byte, len)
	for i := 0; i < len; i++ {
		bytesArgs[i] = []byte(args[i])
	}
	return bytesArgs
}

func readClientConfig(path string) (*ClientConfig, error) {
	conf := &ClientConfig{}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	yaml.NewDecoder(f).Decode(conf)
	return conf, nil
}
