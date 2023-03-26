package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simpleBank/db/sqlc"
	"simpleBank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccount int64  `json:"from_account" binding:"required,min=1"`
	ToAccount   int64  `json:"to_account" binding:"required,min=1"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccount, req.Currency)

	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Username != fromAccount.Owner {
		err := errors.New("aucount does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccount, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccount,
		ToAccountID:   req.ToAccount,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccountByID(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid currency")))
		return account, false
	}

	return account, true
}
