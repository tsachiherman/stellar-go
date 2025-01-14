package history

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stellar/go/services/horizon/internal/test"
)

type transactionParticipant struct {
	TransactionID int64 `db:"history_transaction_id"`
	AccountID     int64 `db:"history_account_id"`
}

func getTransactionParticipants(tt *test.T, q *Q) []transactionParticipant {
	var participants []transactionParticipant
	sql := sq.Select("history_transaction_id", "history_account_id").
		From("history_transaction_participants").
		OrderBy("(history_transaction_id, history_account_id) asc")

	err := q.Select(tt.Ctx, &participants, sql)
	if err != nil {
		tt.T.Fatal(err)
	}

	return participants
}

func TestTransactionParticipantsBatch(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetHorizonDB(t, tt.HorizonDB)
	q := &Q{tt.HorizonSession()}

	batch := q.NewTransactionParticipantsBatchInsertBuilder(0)

	transactionID := int64(1)
	otherTransactionID := int64(2)
	accountID := int64(100)

	for i := int64(0); i < 3; i++ {
		tt.Assert.NoError(batch.Add(tt.Ctx, transactionID, accountID+i))
	}

	tt.Assert.NoError(batch.Add(tt.Ctx, otherTransactionID, accountID))
	tt.Assert.NoError(batch.Exec(tt.Ctx))

	participants := getTransactionParticipants(tt, q)
	tt.Assert.Equal(
		[]transactionParticipant{
			{TransactionID: 1, AccountID: 100},
			{TransactionID: 1, AccountID: 101},
			{TransactionID: 1, AccountID: 102},
			{TransactionID: 2, AccountID: 100},
		},
		participants,
	)
}
