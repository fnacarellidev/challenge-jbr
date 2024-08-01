package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/fnacarellidev/challenge-jbr/backend/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/challenge-jbr/backend/tests/testhelpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BackendUnitTestSuite struct {
    suite.Suite
	pgContainer *testhelpers.PostgresContainer
	sqlcQueries *pgquery.Queries
	ctx context.Context
	conn *pgx.Conn
}

func (suite *BackendUnitTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.conn, err = pgx.Connect(suite.ctx, suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.sqlcQueries = pgquery.New(suite.conn)
}

func (suite *BackendUnitTestSuite) TearDownTestSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	suite.conn.Close(suite.ctx)
}

func (suite *BackendUnitTestSuite) TestGetCourtCase() {
	t := suite.T()

	cnj := "5001682-88.2024.8.13.0672"
	courtCase, err := suite.sqlcQueries.GetCourtCase(suite.ctx, cnj)
	require.NoError(t, err, "should have no error")
	assert.Equal(t, courtCase.Cnj, cnj)
	assert.Equal(t, courtCase.CourtOfOrigin, "TJSP")
	assert.Equal(t, courtCase.Plaintiff, "Alice Johnson")
	assert.Equal(t, courtCase.Defendant, "Bob Smith")
}

func (suite *BackendUnitTestSuite) TestInsertCourtCase() {
	t := suite.T()

	cnj := "5001680-88.2024.8.13.0672"
	plaintiff := "John Doe"
	defendant := "Foo Bar"
	courtOfOrigin := "TJNY"
	startDate := pgtype.Timestamptz{
		Time: time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
	_, err := suite.sqlcQueries.InsertCourtCase(suite.ctx, pgquery.InsertCourtCaseParams{
		Cnj: cnj,
		Plaintiff: plaintiff,
		Defendant: defendant,
		CourtOfOrigin: courtOfOrigin,
		StartDate: startDate,
	})
	require.NoError(t, err, "should have no error")

	courtCase, err := suite.sqlcQueries.GetCourtCase(suite.ctx, cnj)
	require.NoError(t, err, "should have no error")
	assert.Equal(t, courtCase.Cnj, cnj)
	assert.Equal(t, courtCase.Plaintiff, plaintiff)
	assert.Equal(t, courtCase.Defendant, defendant)
	assert.Equal(t, courtCase.CourtOfOrigin, courtOfOrigin)
	assert.Equal(t, courtCase.StartDate, startDate)
}

func TestBackendUnitTestSuite(t *testing.T) {
    suite.Run(t, new(BackendUnitTestSuite))
}
