package models

import "fabric-client/db"

type Blocktxinfo struct {
	Id int `json:"id" xorm:"pk autoincr INT(10) notnull"`
	Number uint64 `json:"number" xorm:"varchar(255) notnull"`
	PreviousHash string `json:"previous_hash" xorm:"varchar(255) notnull"`
	TxId string `json:"tx_id" xorm:"varchar(255) notnull"`
	Timestamp int64 `json:"timestamp" xorm:"bigInt notnull"`
	ChannelId string `json:"channel_id" xorm:"varchar(255) notnull"`
}

//加入信息
func CreateBlockinfo(blocktxinfo ...*Blocktxinfo) (int64,error) {
	e:=db.MasterEngine()
	return	e.Insert(blocktxinfo)
}