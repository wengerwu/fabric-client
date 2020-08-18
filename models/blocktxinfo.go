package models

import (
	"fabric-client/db"
	"fabric-client/util"
)

type BlockTXInfo struct {
	Id int `json:"id" xorm:"pk autoincr INT(10) notnull"`
	Number uint64 `json:"number" xorm:"varchar(255) notnull"`
	PreviousHash string `json:"previous_hash" xorm:"varchar(255) notnull"`
	TxId string `json:"tx_id" xorm:"varchar(255) notnull"`
	Timestamp int64 `json:"timestamp" xorm:"bigInt notnull"`
	ChannelId string `json:"channel_id" xorm:"varchar(255) notnull"`
}

//加入信息
func CreateBlockInfo(blocktxinfo ...*BlockTXInfo) (int64,error) {
	e:=db.MasterEngine()
	return	e.Insert(blocktxinfo)
}

//获取分页区块数据
func GetPaginationBlock(page *util.Pagination) ([]*BlockTXInfo, int64, error) {
	e := db.MasterEngine()
	blocktxinfo := make([]*BlockTXInfo, 0)
	s := e.Limit(page.Limit, page.Start)
	if page.SortName != "" {
		switch page.SortOrder {
		case "asc":
			s.Asc(page.SortName)
		case "desc":
			s.Desc(page.SortName)
		}
	}
	count, err := s.FindAndCount(&blocktxinfo)
	return blocktxinfo, count, err
}

//获取所有区块数据
func GetAllBlock() ([]*BlockTXInfo, int64, error) {
	e := db.MasterEngine()
	blocktxinfo := make([]*BlockTXInfo, 0)
	count, err := e.FindAndCount(&blocktxinfo)
	return blocktxinfo, count, err
}

//根据TxId获取blocktxinfo
func GetBlockByTxId(blocktxinfo *BlockTXInfo) (*BlockTXInfo, error) {
	e := db.MasterEngine()
	_, err := e.Where("tx_id=?", blocktxinfo.TxId).Get(blocktxinfo)
	return blocktxinfo, err
}