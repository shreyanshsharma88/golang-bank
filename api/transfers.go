package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shreyanshsharma88/golang-bank/auth"
	db "github.com/shreyanshsharma88/golang-bank/db/sqlc"
)

type createTransferStruct struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) createTransfers(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*auth.Payload)
	var req createTransferStruct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(c, req.FromAccountID, req.Currency)

	if !valid {
		return
	}
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if fromAccount.Balance < req.Amount {
		err := errors.New("insufficient balance")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	_, valid = server.validAccount(c, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := server.store.CreateTransfer(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transfer)

}
func (server *Server) validAccount(c *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			err := fmt.Errorf("account not found: %d", accountId)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account currency mismatch: %s != %s", account.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
