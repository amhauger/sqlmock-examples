package basics

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testDB        *sql.DB
	dbMocker      sqlmock.Sqlmock
	mockDB        BasicsDB
	valuesColumns *sqlmock.Rows
)

func TestBasics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Basics Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testDB, dbMocker, err = sqlmock.New()
	Expect(err).Should(BeNil(), "error creating mock sql db in BeforeSuite")

	gdb, err := gorm.Open("postgres", testDB)
	Expect(err).Should(BeNil(), "error using mock sql db as gormDB in BeforeSuite")
	mockDB = BasicsDB{
		Connection: gdb,
	}
})

var _ = Describe("SQL Mock Basics", func() {
	BeforeEach(func() {
		valuesColumns = sqlmock.NewRows([]string{"val"})
	})

	Context("reseting the mock db each test", func() {
		BeforeEach(func() {
			var err error
			testDB, dbMocker, err = sqlmock.New()
			Expect(err).Should(BeNil(), "error creating mock sql db in BeforeSuite")

			gdb, err := gorm.Open("postgres", testDB)
			Expect(err).Should(BeNil(), "error using mock sql db as gormDB in BeforeSuite")
			mockDB = BasicsDB{
				Connection: gdb,
			}
		})

		It("should return no values", func() {
			dbMocker.
				ExpectQuery("SELECT \\* FROM \"values\" WHERE \\(val=\\$1\\)").
				WithArgs(1).WillReturnRows(valuesColumns)

			rows, err := mockDB.GetValFromSql(1)
			Expect(err).Should(BeNil())
			Expect(rows).Should(BeNil())
			Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
		})

		It("should return a value", func() {
			dbMocker.
				ExpectQuery("SELECT \\* FROM \"values\" WHERE .*").
				WithArgs(1).
				WillReturnRows(valuesColumns.AddRow(1))

			rows, err := mockDB.GetValFromSql(1)
			Expect(err).Should(BeNil())
			Expect(rows).ShouldNot(BeNil())
			Expect(dbMocker.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})

		It("should return an error", func() {
			dbMocker.
				ExpectQuery("^SELECT").
				WithArgs(1).WillReturnError(errors.New("an error"))

			rows, err := mockDB.GetValFromSql(1)
			Expect(err).ShouldNot(BeNil())
			Expect(rows).Should(BeNil())
			Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
		})
	})

	Context("not reseting the mock db each test", func() {
		Context("testing GetEverything()", func() {
			It("will return everything and nothing at all", func() {
				dbMocker.
					ExpectQuery("SELECT \\* ").
					WillReturnRows(sqlmock.NewRows([]string{"1"}))

				returnedRows, returnedError := mockDB.GetEverything()
				Expect(returnedError).Should(BeNil())
				Expect(returnedRows).ShouldNot(BeNil())

				dbMocker.ExpectQuery("SELECT \\* FROM .*")

				mockDB.ReturnNothing()
				Expect(dbMocker.ExpectationsWereMet()).ShouldNot(BeNil())
			})

			It("will break", func() {
				Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
			})
		})
	})
})
