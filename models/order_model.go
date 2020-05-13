package models

import "github.com/jinzhu/gorm"

const (
	EthToTtOrderType = 1 // 以太坊转tt链
	TtToEthOrderType = 2 // tt链转以太坊
)

type Order struct {
	FromAddr      string
	RecipientAddr string
	Amount        string
	OrderType     int
	State         int  // 订单状态, 0: pending; 1. 完成；2. 失败; 3. 超时
	CollectionId  uint // collection 表的外键
	gorm.Model
}

// Create
func (o *Order) Create() error {
	if err := DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

// GetByAddr
func (o *Order) GetByAddr() (*Order, error) {
	oo := Order{}
	err := DB.Where(o).Last(&oo).Error
	return &oo, err
}

func (o *Order) Update(oo Order) error {
	return DB.Model(o).Updates(oo).Error
}
