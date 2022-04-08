package intermediates

import (
	"database/sql"
	"testing"
	"time"

	driver "database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testDB        *sql.DB
	dbMocker      sqlmock.Sqlmock
	mockDB        IntermediateDB
	valuesColumns *sqlmock.Rows
)

func TestIntermediates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Intermediates Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testDB, dbMocker, err = sqlmock.New()
	Expect(err).Should(BeNil(), "error creating mock sql db in BeforeSuite")

	gdb, err := gorm.Open("postgres", testDB)
	Expect(err).Should(BeNil(), "error using mock sql db as gormDB in BeforeSuite")
	mockDB = IntermediateDB{
		Connection: gdb,
	}
})

var _ = Describe("SQL Mock Intermediates", func() {
	BeforeEach(func() {
		valuesColumns = sqlmock.NewRows([]string{"id", "val", "effective_at"})

		var err error
		testDB, dbMocker, err = sqlmock.New()
		Expect(err).Should(BeNil(), "error creating mock sql db in BeforeSuite")

		gdb, err := gorm.Open("postgres", testDB)
		Expect(err).Should(BeNil(), "error using mock sql db as gormDB in BeforeSuite")
		mockDB = IntermediateDB{
			Connection: gdb,
		}
	})

	Context("without using AnyArgs", func() {
		It("will give us an unexpected error", func() {
			dbMocker.ExpectBegin()
			dbMocker.ExpectExec("UPDATE \"values\" .*").
				WithArgs(1, time.Now(), 1).
				WillReturnResult(sqlmock.NewResult(1, 1)) // inserted row num, num rows effected
			dbMocker.ExpectCommit()

			err := mockDB.InsertVal()
			Expect(err).Should(BeNil())
			Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
		})
	})

	Context("using AnyArgs", func() {
		var valuesRow *sqlmock.Rows
		BeforeEach(func() {
			valuesRow = valuesColumns.AddRow("11111111-1111-1111-1111-111111111111", "1", time.Now())
		})

		Context("and only AnyArgs", func() {
			It("will return expected rows", func() {
				dbMocker.ExpectQuery("SELECT").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(valuesRow)

				val, err := mockDB.GetValFromSqlByID()
				Expect(err).Should(BeNil())
				Expect(val.Val).Should(BeNumerically("==", 1))
				Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
			})
		})

		Context("using our generated Args", func() {
			It("should evaluate a UUID correctly", func() {
				dbMocker.ExpectQuery("SELECT").
					WithArgs(AnyUUID{}).
					WillReturnRows(valuesRow)

				val, err := mockDB.GetValFromSqlByID()
				Expect(err).Should(BeNil())
				Expect(val.Val).Should(BeNumerically("==", 1))
				Expect(dbMocker.ExpectationsWereMet()).Should(BeNil())
			})
		})
	})
})

type AnyUUID struct{}

func (a AnyUUID) Match(v driver.Value) bool {
	val, ok := v.(string)
	if len(val) != 36 {
		ok = false
	}
	return ok
}
