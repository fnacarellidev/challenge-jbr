package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/fnacarellidev/challenge-jbr/testutil"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/stretchr/testify/suite"
)

type GraphQLApiTestSuite struct {
	suite.Suite
	pgContainer *types.PostgresContainer
	ctx context.Context
}

func (suite *GraphQLApiTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutil.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	os.Setenv("DATABASE_URL", suite.pgContainer.ConnectionString)
}

func (suite *GraphQLApiTestSuite) TearDownTestSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestGraphQLApiSuite(t *testing.T) {
    suite.Run(t, new(GraphQLApiTestSuite))
}
