package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zyjblockchain/sandy_log/log"
	"github.com/zyjblockchain/tt_tac/logics"
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
		log.Infof("发送跨链转账交易前端传入参数：Address：%s, amount: %s, orderType: %d", logic.Address, logic.Amount, logic.OrderType)
		// logic
		txHash, err := logic.SendTacTx()
		if err != nil {
			log.Errorf("send tac tx logic err: %v", err)
			serializer.ErrorResponse(c, utils.SendTacTxErrCode, utils.SendTacTxErrMsg, err.Error())
			return
		} else {
			serializer.SuccessResponse(c, tHash{TxHash: txHash}, "success")
		}
	}
}
