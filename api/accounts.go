package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/shreyanshsharma88/golang-bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD INR"`
}

type getAccount struct {
	ID int64 `uri:"id"`
}

type listAccounts struct {
	PageID int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

type updateAccountRequest struct {
	Currency string `json:"currency" binding:"oneof=USD INR"`
	Balance int64 `json:"balance" binding:"required"` 
	Owner string `json:"owner" binding:"required"`

}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(c, args)
	if err != nil {

		if pqErr, ok := err.(*pq.Error) ; ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation" , "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)

}

func (server *Server) getAccount(c *gin.Context) {

	var req getAccount
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	c.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(c *gin.Context) {
	var req listAccounts

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	accounts , err := server.store.ListAccounts(c, db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	c.JSON(http.StatusOK, accounts)
	
}

func ( server *Server) deleteAccount (c *gin.Context) {
	var req getAccount

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(c, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, "account deleted")
	
}

func (server *Server) updateAccount(c *gin.Context) {
	var req updateAccountRequest
	var params getAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.UpdateAccount(c, db.UpdateAccountParams{
		Balance: req.Balance,
		Owner: req.Owner,

	})

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)
}