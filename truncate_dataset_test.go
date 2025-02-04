package goqu

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v8/exec"
	"github.com/doug-martin/goqu/v8/exp"
	"github.com/doug-martin/goqu/v8/internal/errors"
	"github.com/doug-martin/goqu/v8/internal/sb"
	"github.com/doug-martin/goqu/v8/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	truncateTestCase struct {
		ds      *TruncateDataset
		clauses exp.TruncateClauses
	}
	truncateDatasetSuite struct {
		suite.Suite
	}
)

func (tds *truncateDatasetSuite) assertCases(cases ...truncateTestCase) {
	for _, s := range cases {
		tds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (tds *truncateDatasetSuite) TestClone() {
	ds := Truncate("test")
	tds.Equal(ds, ds.Clone())
}

func (tds *truncateDatasetSuite) TestExpression() {
	ds := Truncate("test")
	tds.Equal(ds, ds.Expression())
}

func (tds *truncateDatasetSuite) TestDialect() {
	ds := Truncate("test")
	tds.NotNil(ds.Dialect())
}

func (tds *truncateDatasetSuite) TestWithDialect() {
	ds := Truncate("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := GetDialect("default")
	dialectDs := ds.WithDialect("default")
	tds.Equal(md, ds.Dialect())
	tds.Equal(dialect, dialectDs.Dialect())
}

func (tds *truncateDatasetSuite) TestPrepared() {
	ds := Truncate("test")
	preparedDs := ds.Prepared(true)
	tds.True(preparedDs.IsPrepared())
	tds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	tds.True(preparedDs.Restrict().IsPrepared())
}

func (tds *truncateDatasetSuite) TestGetClauses() {
	ds := Truncate("test")
	ce := exp.NewTruncateClauses().SetTable(exp.NewColumnListExpression(I("test")))
	tds.Equal(ce, ds.GetClauses())
}

func (tds *truncateDatasetSuite) TestTable() {
	bd := Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Table("test2"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test2")),
		},
		truncateTestCase{
			ds: bd.Table("test1", "test2"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test1", "test2")),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestCascade() {
	bd := Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Cascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true}),
		},
		truncateTestCase{
			ds: bd.Restrict().Cascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestNoCascade() {
	bd := Truncate("test").Cascade()
	tds.assertCases(
		truncateTestCase{
			ds: bd.NoCascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{}),
		},
		truncateTestCase{
			ds: bd.Restrict().NoCascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: false, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true}),
		},
	)
}

func (tds *truncateDatasetSuite) TestRestrict() {
	bd := Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Restrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Restrict: true}),
		},
		truncateTestCase{
			ds: bd.Cascade().Restrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestNoRestrict() {
	bd := Truncate("test").Restrict()
	tds.assertCases(
		truncateTestCase{
			ds: bd.NoRestrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{}),
		},
		truncateTestCase{
			ds: bd.Cascade().NoRestrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: false}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Restrict: true}),
		},
	)
}

func (tds *truncateDatasetSuite) TestIdentity() {
	bd := Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Identity("RESTART"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Identity: "RESTART"}),
		},
		truncateTestCase{
			ds: bd.Identity("CONTINUE"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Identity: "CONTINUE"}),
		},
		truncateTestCase{
			ds: bd.Cascade().Restrict().Identity("CONTINUE"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true, Identity: "CONTINUE"}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestToSQL() {
	md := new(mocks.SQLDialect)
	ds := Truncate("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToTruncateSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	tds.NoError(err)
	tds.Empty(sql)
	tds.Empty(args)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestToSQL__withPrepared() {
	md := new(mocks.SQLDialect)
	ds := Truncate("test").Prepared(true).SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(true)
	md.On("ToTruncateSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Nil(err)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestToSQL_withError() {
	md := new(mocks.SQLDialect)
	ds := Truncate("test").SetDialect(md)
	c := ds.GetClauses()
	ee := errors.New("expected error")
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToTruncateSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(sb.SQLBuilder).SetError(ee)
	}).Once()

	sql, args, err := ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(ee, err)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestExecutor() {
	mDb, _, err := sqlmock.New()
	tds.NoError(err)

	ds := newTruncateDataset("mock", exec.NewQueryFactory(mDb)).
		Table("table1", "table2")

	tsql, args, err := ds.Executor().ToSQL()
	tds.NoError(err)
	tds.Empty(args)
	tds.Equal(`TRUNCATE "table1", "table2"`, tsql)

	tsql, args, err = ds.Prepared(true).Executor().ToSQL()
	tds.NoError(err)
	tds.Empty(args)
	tds.Equal(`TRUNCATE "table1", "table2"`, tsql)
}

func TestTruncateDataset(t *testing.T) {
	suite.Run(t, new(truncateDatasetSuite))
}
