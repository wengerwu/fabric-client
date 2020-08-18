package controllers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fabric-client/models"
	"fabric-client/sdkInit"
	"fabric-client/service"
	"fabric-client/util"
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/kataras/iris/v12/middleware/i18n"
)

type FabricSDKController struct {
	Ctx       iris.Context
	ClientMap map[string]*sdkInit.Client
}

type ChannelRequest struct {
	ChannelID string // 通道ID
	OrgName   string // 组织名

	Timestamp int64  //时间戳
	Sign      string //签名
}

type ChaincodeRequest struct {
	ChannelID        string
	OrgName          string
	UserName         string
	ChaincodeID      string
	Fcn              string
	Args             []string
	EventFilter      string //查询链码不用传
	EventCallbackUrl string //查询链码不用传
	Timestamp        int64
	Sign             string
}

type BlcockInfo struct {
	Number       uint64
	PreviousHash string
}

// 创建通道
func (controller *FabricSDKController) PostChannelCreate() Result {
	channelRequest := &ChannelRequest{}
	if result := controller.parseJson(channelRequest); result.Code != OK {
		return result
	}

	src := "orgName=" + channelRequest.OrgName + "&channelID=" + channelRequest.ChannelID + "&timestamp=" + strconv.FormatInt(channelRequest.Timestamp, 10)
	if result := controller.checkSign(channelRequest.Timestamp, channelRequest.Sign, src); result.Code != OK {
		return result
	}

	client, result := controller.getAndCheckClient(channelRequest.OrgName)
	if result.Code != OK {
		return result
	}

	err := client.CreateChannel(channelRequest.ChannelID)
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(CreateChannelError, i18n.Translate(controller.Ctx, "create_channel_fail"), err.Error())
	}

	return Result{Code: OK, Message: i18n.Translate(controller.Ctx, "create_channel_success")}
}

// 加入通道
func (controller *FabricSDKController) PostChaincodeJoin() Result {
	channelRequest := &ChannelRequest{}
	if result := controller.parseJson(channelRequest); result.Code != OK {
		return result
	}

	src := "orgName=" + channelRequest.OrgName + "&channelID=" + channelRequest.ChannelID + "&timestamp=" + strconv.FormatInt(channelRequest.Timestamp, 10)
	if result := controller.checkSign(channelRequest.Timestamp, channelRequest.Sign, src); result.Code != OK {
		return result
	}

	client, result := controller.getAndCheckClient(channelRequest.OrgName)
	if result.Code != OK {
		return result
	}

	err := client.JoinChannel(channelRequest.ChannelID)
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(JoinChannelError, i18n.Translate(controller.Ctx, "join_channel_fail"), err.Error())
	}

	return Result{Code: OK, Message: i18n.Translate(controller.Ctx, "join_channel_success")}
}

// 安装链码
func (controller *FabricSDKController) PostChaincodeInstall() Result {
	ccRequest := &sdkInit.CCRequest{}
	if result := controller.parseJson(ccRequest); result.Code != OK {
		return result
	}

	src := "orgName=" + ccRequest.OrgName + "&chaincodeID=" + ccRequest.ChaincodeID + "&chaincodeVersion=" + ccRequest.ChaincodeVersion + "&chaincodePath=" + ccRequest.ChaincodePath + "&timestamp=" + strconv.FormatInt(ccRequest.Timestamp, 10)
	if result := controller.checkSign(ccRequest.Timestamp, ccRequest.Sign, src); result.Code != OK {
		return result
	}

	client, result := controller.getAndCheckClient(ccRequest.OrgName)
	if result.Code != OK {
		return result
	}

	err := client.InstallCC(ccRequest)
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(InstallCCError, i18n.Translate(controller.Ctx, "install_cc_fail"), err.Error())
	}

	return Result{Code: OK, Message: i18n.Translate(controller.Ctx, "install_cc_success")}
}

// 实例化链码
func (controller *FabricSDKController) PostChaincodeInstantiate() Result {
	ccRequest := &sdkInit.CCRequest{}
	if result := controller.parseJson(ccRequest); result.Code != OK {
		return result
	}

	src := "channelID=" + ccRequest.ChannelID + "&orgName=" + ccRequest.OrgName + "&chaincodeID=" + ccRequest.ChaincodeID + "&chaincodeVersion=" + ccRequest.ChaincodeVersion + "&chaincodePath=" + ccRequest.ChaincodePath + "&timestamp=" + strconv.FormatInt(ccRequest.Timestamp, 10)
	if result := controller.checkSign(ccRequest.Timestamp, ccRequest.Sign, src); result.Code != OK {
		return result
	}

	client, result := controller.getAndCheckClient(ccRequest.OrgName)
	if result.Code != OK {
		return result
	}

	err := client.InstantiateCC(ccRequest)
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(InstantiateCCError, i18n.Translate(controller.Ctx, "instantiate_cc_fail"), err.Error())
	}

	return Result{Code: OK, Message: i18n.Translate(controller.Ctx, "instantiate_cc_success")}
}

// 升级链码
func (controller *FabricSDKController) PostChaincodeUpgrade() Result {
	ccRequest := &sdkInit.CCRequest{}
	if result := controller.parseJson(ccRequest); result.Code != OK {
		return result
	}

	src := "channelID=" + ccRequest.ChannelID + "&orgName=" + ccRequest.OrgName + "&chaincodeID=" + ccRequest.ChaincodeID + "&chaincodeVersion=" + ccRequest.ChaincodeVersion + "&chaincodePath=" + ccRequest.ChaincodePath + "&timestamp=" + strconv.FormatInt(ccRequest.Timestamp, 10)
	if result := controller.checkSign(ccRequest.Timestamp, ccRequest.Sign, src); result.Code != OK {
		return result
	}

	client, result := controller.getAndCheckClient(ccRequest.OrgName)
	if result.Code != OK {
		return result
	}

	err := client.UpgradeCC(ccRequest)
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(UpgradeCCError, i18n.Translate(controller.Ctx, "upgrade_cc_fail"), err.Error())
	}

	return Result{Code: OK, Message: i18n.Translate(controller.Ctx, "upgrade_cc_success")}
}

type TxInfo struct {
	Channel   string
	TxID      string
	Timestamp []byte
}

// 链码执行
func (controller *FabricSDKController) PostChaincodeExec() Result {
	chaincodeRequest := &ChaincodeRequest{}
	if result := controller.parseJson(chaincodeRequest); result.Code != OK {
		return result
	}

	if len(chaincodeRequest.Args) < 1 {
		return controller.getInternalServerError(ArgsError, i18n.Translate(controller.Ctx, "cc_args_len_error", 1), nil)
	}

	src := "args[0]=" + chaincodeRequest.Args[0] + "&channelID=" + chaincodeRequest.ChannelID + "&orgName=" + chaincodeRequest.OrgName + "&userName=" + chaincodeRequest.UserName
	if len(chaincodeRequest.Args) > 1 {
		for i := 1; i < len(chaincodeRequest.Args); i++ {
			src += "&args[" + strconv.Itoa(i) + "]=" + chaincodeRequest.Args[i]
		}
	}
	src += "&chaincodeID=" + chaincodeRequest.ChaincodeID + "&fcn=" + chaincodeRequest.Fcn + "&eventCallbackUrl=" + chaincodeRequest.EventCallbackUrl + "&timestamp=" + strconv.FormatInt(chaincodeRequest.Timestamp, 10)
	if result := controller.checkSign(chaincodeRequest.Timestamp, chaincodeRequest.Sign, src); result.Code != OK {
		return result
	}

	serviceSetup, result := controller.getServiceSetup(chaincodeRequest)
	if result.Code != OK {
		return result
	}

	if chaincodeRequest.EventFilter != "" {
		go serviceSetup.SetEvent(chaincodeRequest.EventFilter, chaincodeRequest.EventCallbackUrl, func(eventFilter string, eventCallbackUrl string, event *fab.CCEvent) {
			if event != nil {
				fmt.Printf("事件的TxID是%s\n", event.TxID)
				if eventCallbackUrl != "" {
					//用http发送event对象到callbackUrl
					data, err := json.Marshal(event)
					if err != nil {
						fmt.Println(err)
						return
					}

					resp, err := http.Post(eventCallbackUrl, "application/json", bytes.NewReader(data))
					defer resp.Body.Close()
					if err != nil {
						fmt.Println(err)
						return
					}

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err)
						return
					}

					fmt.Println(string(body))
				}
			} else {
				fmt.Printf("不能根据指定的事件ID接收到相应的链码事件(%s)\n", eventFilter)
			}
		})
	}

	response, err := serviceSetup.Execute(chaincodeRequest.Fcn, sdkInit.ToBytesArgs(chaincodeRequest.Args))
	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(ExecCCError, i18n.Translate(controller.Ctx, "exec_cc_fail"), err.Error())
	}

	txInfoPayload := response.Payload
	txInfo := &TxInfo{}
	json.Unmarshal(txInfoPayload, txInfo)
	timestamp := &timestamp.Timestamp{}
	json.Unmarshal(txInfo.Timestamp, timestamp)
	fmt.Println(txInfo)
	fmt.Println(timestamp.Seconds)

	block, err :=serviceSetup.QueryBlockByTxID(response.TransactionID)
	if err != nil {
		fmt.Printf("根据id获取交易信息失败: %s\n", err)
		return controller.getInternalServerError(QueryBlockByIdError, i18n.Translate(controller.Ctx, "query_block_by_id"), err.Error())
	}

	previousHash := hex.EncodeToString(block.Header.PreviousHash)
	fmt.Printf("区块上一个hash: %s\n", previousHash)

	blockTXInfo := new(models.BlockTXInfo)
	blockTXInfo.Number = block.Header.Number
	blockTXInfo.PreviousHash = previousHash
	blockTXInfo.TxId = txInfo.TxID
	blockTXInfo.Timestamp = timestamp.Seconds
	blockTXInfo.ChannelId = txInfo.Channel

	_, err = models.CreateBlockInfo(blockTXInfo)
	if err != nil {
		return controller.getInternalServerError(QueryBlockError, i18n.Translate(controller.Ctx, "insert_block_database_fail"), err.Error())
	}

	fmt.Printf("执行链码成功，交易hash:%s\n", response.TransactionID)
	return Result{OK, i18n.Translate(controller.Ctx, "exec_cc_success"), response}
}

//测试用http发送event对象到callbackUrl
func (controller *FabricSDKController) PostCallback() Result {
	event := &fab.CCEvent{}
	if result := controller.parseJson(event); result.Code != OK {
		return result
	}

	return Result{OK, "callbackUrl的event对象", event}
}

// 链码查询
func (controller *FabricSDKController) PostChaincodeQuery() Result {
	chaincodeRequest := &ChaincodeRequest{}
	if result := controller.parseJson(chaincodeRequest); result.Code != OK {
		return result
	}

	if len(chaincodeRequest.Args) < 1 {
		return controller.getInternalServerError(QueryCCError, i18n.Translate(controller.Ctx, "cc_args_len_error", 1), nil)
	}

	src := "args[0]=" + chaincodeRequest.Args[0] + "&channelID=" + chaincodeRequest.ChannelID + "&orgName=" + chaincodeRequest.OrgName + "&userName=" + chaincodeRequest.UserName
	if len(chaincodeRequest.Args) > 1 {
		for i := 1; i < len(chaincodeRequest.Args); i++ {
			src += "&args[" + strconv.Itoa(i) + "]=" + chaincodeRequest.Args[i]
		}
	}
	src += "&chaincodeID=" + chaincodeRequest.ChaincodeID + "&fcn=" + chaincodeRequest.Fcn + "&timestamp=" + strconv.FormatInt(chaincodeRequest.Timestamp, 10)
	if result := controller.checkSign(chaincodeRequest.Timestamp, chaincodeRequest.Sign, src); result.Code != OK {
		return result
	}

	serviceSetup, result := controller.getServiceSetup(chaincodeRequest)
	if result.Code != OK {
		return result
	}

	response, err := serviceSetup.Query(chaincodeRequest.Fcn, sdkInit.ToBytesArgs(chaincodeRequest.Args))

	if err != nil {
		fmt.Println(err.Error())
		return controller.getInternalServerError(QueryCCError, i18n.Translate(controller.Ctx, "query_cc_fail"), err.Error())
	}

	fmt.Println(response.Responses[0].Timestamp)
	return Result{OK, i18n.Translate(controller.Ctx, "query_cc_success"), response}
}

//区块分页查询
func (controller *FabricSDKController) GetPaginationBlock() Result {
	page, err := util.NewPagination(controller.Ctx)
	if err != nil {
		return controller.getInternalServerError(iris.StatusInternalServerError, i18n.Translate(controller.Ctx, "get_page_data_fail"), err.Error())
	}

	src:="SortOrder="+page.SortOrder+"&PageNumber="+strconv.Itoa(page.PageNumber)+"&SortName="+page.SortName+"&StartDate="+"&Limit="+strconv.Itoa(page.Limit)+ "&timestamp=" + strconv.FormatInt(page.Timestamp, 10)
	if result := checkSign(controller.Ctx,page.Timestamp,page.Sign,src); result.Code != OK {
		return result
	}

	blocks, count, err := models.GetPaginationBlock(page)
	if err != nil {
		return controller.getInternalServerError(iris.StatusInternalServerError, i18n.Translate(controller.Ctx, "get_block_fail"), err.Error())
	}
	response := util.BootstrapTableVO{
		Total: count,
		Rows:  blocks,
	}
	return Result{OK, i18n.Translate(controller.Ctx, "get_block_success"), response}
}

func (controller *FabricSDKController) parseJson(jsonObjectPtr interface{}) Result {
	return parseJson(controller.Ctx, jsonObjectPtr)
}

func (controller *FabricSDKController) checkSign(timestamp int64, sign string, src string) Result {
	return checkSign(controller.Ctx, timestamp, sign, src)
}

func (controller *FabricSDKController) getServiceSetup(chaincodeRequest *ChaincodeRequest) (*service.Setup, Result) {
	client, result := controller.getAndCheckClient(chaincodeRequest.OrgName)
	if result.Code != OK {
		return nil, result
	}

	var err error
	key := chaincodeRequest.ChannelID + chaincodeRequest.OrgName + chaincodeRequest.UserName
	channelClientRequest := &sdkInit.ChannelClientRequest{
		ChannelID: chaincodeRequest.ChannelID,
		OrgName:   chaincodeRequest.OrgName,
		UserName:  chaincodeRequest.UserName,
	}

	channelClient, ok := client.ChannelClients[key]
	if !ok {
		channelClient, err = client.NewChannelClient(channelClientRequest)
		if err != nil {
			fmt.Println(err.Error())
			return nil, controller.getInternalServerError(NewChannelClientError, i18n.Translate(controller.Ctx, "new_channelclient_fail"), err.Error())
		}
		client.ChannelClients[key] = channelClient
	}

	ledgerClient,ok := client.LedgerClients[key]
	if !ok {
		ledgerClient, err = client.NewLedgerClient(channelClientRequest)
		if err != nil {
			fmt.Println(err.Error())
			return nil, controller.getInternalServerError(NewLedgerClientError, i18n.Translate(controller.Ctx, "new_ledgerclient_fail"), err.Error())
		}
		client.LedgerClients[key] = ledgerClient
	}

	serviceSetup := &service.Setup{ChaincodeID: chaincodeRequest.ChaincodeID, Client: channelClient,LClient: ledgerClient}
	return serviceSetup, Result{Code: OK}
}

func (controller *FabricSDKController) getAndCheckClient(orgName string) (*sdkInit.Client, Result) {
	client, ok := controller.ClientMap[orgName]
	if !ok {
		return nil, controller.getInternalServerError(GetAndCheckClientError, i18n.Translate(controller.Ctx, "get_and_check_client", orgName), orgName)
	}
	return client, Result{Code: OK}
}

func (controller *FabricSDKController) getInternalServerError(code int, message string, data interface{}) Result {
	return getInternalServerError(controller.Ctx, code, message, data)
}
