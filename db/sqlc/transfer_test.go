package db

import (
	"context"
	"database/sql"
	"go-bank/util"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)


func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	log.Println(arg.FromAccountID)
	log.Println(transfer.FromAccountID)

	// Cast the FromAccountID and the ToAccountID to sql.NullInt64 
	accountID1 := sql.NullInt64{Int64: arg.FromAccountID, Valid: true}
	accountID2 := sql.NullInt64{Int64: arg.ToAccountID, Valid: true}

	require.Equal(t, accountID1, transfer.FromAccountID)
	require.Equal(t, accountID2, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(),transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt.Time, transfer2.CreatedAt.Time, time.Second)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	// Cast the FromAccountID and the ToAccountID to sql.NullInt64 
	accountID1 := sql.NullInt64{Int64: account1.ID, Valid: true}

	arg := ListTransfersParams {
		FromAccountID: accountID1,
		ToAccountID: accountID1,
		Offset: 5,
		Limit: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	log.Println(transfers)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID.Valid == accountID1.Valid || transfer.ToAccountID.Valid == accountID1.Valid)
	}
}