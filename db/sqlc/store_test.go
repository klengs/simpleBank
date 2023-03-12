package db

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	accountFrom := CreateRandomAccount(t)
	accountTo := CreateRandomAccount(t)

	fmt.Println(">> before:", accountFrom.Balance, accountTo.Balance)

	n := 100
	amount := int64(1)

	arg := TransferTxParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        amount,
	}

	errsCh := make(chan error)
	resultsCh := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100000)) * time.Nanosecond)

			result, err := testStore.TransferTx(context.Background(), arg)

			errsCh <- err
			resultsCh <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errsCh

		require.NoError(t, err)

		result := <-resultsCh

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, accountFrom.ID)
		require.Equal(t, transfer.ToAccountID, accountTo.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check FromEntry
		fromEntry := result.FromEntry
		require.Equal(t, fromEntry.AccountID, accountFrom.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// check ToEntry
		toEntry := result.ToEntry
		require.Equal(t, toEntry.AccountID, accountTo.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check Accounts
		fromAccountResult := result.FromAccount
		require.NotEmpty(t, fromAccountResult)
		require.Equal(t, fromAccountResult.ID, accountFrom.ID)

		toAccountResult := result.ToAccount
		require.NotEmpty(t, toAccountResult)
		require.Equal(t, toAccountResult.ID, accountTo.ID)

		fmt.Println(">> tx:", fromAccountResult.Balance, toAccountResult.Balance)
		// check Balance

		diff1 := accountFrom.Balance - fromAccountResult.Balance
		diff2 := toAccountResult.Balance - accountTo.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k > 0 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccountByID(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccountByID(context.Background(), accountTo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, updatedAccount1.Balance, accountFrom.Balance-int64(n*int(amount)))
	require.Equal(t, updatedAccount2.Balance, accountTo.Balance+int64(n*int(amount)))

}
