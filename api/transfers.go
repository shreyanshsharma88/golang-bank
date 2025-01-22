package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/shreyanshsharma88/golang-bank/db/sqlc"
)

type createTransferStruct struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID int64 `json:"to_account_id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}
func (server *Server) createTransfers(c *gin.Context) {
	var req createTransferStruct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(c, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(c, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}
	transfer, err := server.store.CreateTransfer(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transfer)
	
}
func (server *Server) validAccount (c *gin.Context , accountId int64 , currency string) bool {
	account, err := server.store.GetAccount(c, accountId)
	if err != nil {
		return false
	}
	if account.Currency != currency {
		return false
	}
	return true
}