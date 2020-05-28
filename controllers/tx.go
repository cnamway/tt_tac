package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zyjblockchain/sandy_log/log"
	"github.com/zyjblockchain/tt_tac/logics"
	"github.com/zyjblockchain/tt_tac/models"
	"github.com/zyjblockchain/tt_tac/serializer"
	"github.com/zyjblockchain/tt_tac/utils"
)

type tHash struct {
	TxHash string `json:"tx_hash"`
}

// SendTacTx 发送跨链转账
func SendTacTx() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logic logics.SendTacTx
		err := c.ShouldBind(&logic)
		if err != nil {
			log.Errorf("ApplyOrder should binding error: %v", err)
			serializer.ErrorResponse(c, utils.VerifyParamsErrCode, utils.VerifyParamsErrMsg, err.Error())
			return
		}
		// 把传入的amount换算成最小单位
		logic.Amount = utils.FormatTokenAmount(logic.Amount, 8)
		log.Infof("发送跨链转账交易前端传入参数：Address：%s, amount: %s, orderType: %d, tacOrderId: %d", logic.Address, logic.Amount, logic.OrderType, logic.TacOrderId)

		// logic
		// 查看是否存在订单
		tacOrd, err := (&models.TacOrder{Model: gorm.Model{ID: logic.TacOrderId}}).GetOrder()
		if err != nil {
			// 表示不存在ord
			log.Errorf("跨链转账申请者发送交易到中间地址传入的orderId在tacOrder表中不存在。 db error: %v, tacOrderId: %d", err, logic.TacOrderId)
			serializer.ErrorResponse(c, utils.SendTacTxErrCode, utils.SendTacTxErrMsg, errors.New("跨链转账申请者发送交易到中间地址传入的orderId在tacOrder表中不存在").Error())
			return
		}
		txHash, err := logic.SendTacTx()
		if err != nil {
			// 发送跨链转账发送者转账到中转地址交易失败，则把跨链转账订单置位失败状态
			log.Errorf("跨链转账申请者发送跨链转账交易失败, 把申请跨链转账订单置位失败状态；error: %v", err)
			oo := &models.TacOrder{}
			oo.ID = logic.TacOrderId
			_ = oo.Update(models.TacOrder{State: 2}).Error()
			serializer.ErrorResponse(c, utils.SendTacTxErrCode, utils.SendTacTxErrMsg, err.Error())
			return
		} else {
			// sendTxHash更新到tacOrder中
			_ = tacOrd.Update(models.TacOrder{SendTxHash: txHash})
			serializer.SuccessResponse(c, tHash{TxHash: txHash}, "success")
		}
	}
}
