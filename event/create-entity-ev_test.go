package event

import (
	"context"
	"testing"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/ajainc/chain/ctx/docdb"
	"github.com/ajainc/chain/ctx/docdb/dbtestutil"
	"github.com/ajainc/chain/ctx/logger"
	"github.com/ajainc/chain/event/dao"
	"github.com/ajainc/chain/grpc/gen"
	"github.com/stretchr/testify/assert"
)

type crateEntityTestData struct {
	input         []CreateEntityEv
	expected      []gen.Entity
	optionalCheck func(string, int, *db.DB, CreateEntityEv, gen.Entity, gen.Entity)
	errChech      func(string, error)
	redoCheck     func(string, int, *db.DB, CreateEntityEv)
}

func TestCreateEntity(t *testing.T) {

	var datas []crateEntityTestData
	{

		input := gen.Entity{
			ObjectID: "test_id",
			Coord: &gen.Coord{
				X: 110.5,
				Y: 223.2,
			},
			WidthHeight: &gen.WidthHeight{
				W: 1000.1,
				H: 800.2,
			},
			Columns: []*gen.Column{},
		}
		expected := input
		expected.Name = "NewEntity"

		d := crateEntityTestData{}
		d.input = []CreateEntityEv{{Entity: input}}
		d.expected = []gen.Entity{expected}

		datas = append(datas, d)
	}

	for _, data := range datas {
		doErrCheck := data.errChech != nil
		didErrCheck := false
		doRedoCheck := data.redoCheck != nil

		c := context.Background()
		c, _ = logger.WithContext(c, logger.DiscardLogger)
		c, closeDB, err := dbtestutil.WithTestDBContext(c)

		errCheckIfNeed := func(data crateEntityTestData, err error) {
			if err != nil {
				// test error case if func for check is set
				if doErrCheck {
					didErrCheck = true
					data.errChech("errCheck", err)
				} else {
					t.Error(err)
				}
			}
		}

		func() {
			defer closeDB()
			if err != nil {
				panic(err)
			}
			d := docdb.FromContext(c)

			for idx, ev := range data.input {
				err = ev.Exec(c)

				if err != nil {
					// test error case if func for check is set
					errCheckIfNeed(data, err)
					continue
				}

				_, result, err := docdb.GetByObjectID(d, docdb.COLL_ENTITY, ev.Entity.ObjectID)
				errCheckIfNeed(data, err)

				var actual gen.Entity
				err = dao.UnmarshalEntity(result, &actual)
				errCheckIfNeed(data, err)

				assert.Equal(t, data.expected[idx], actual)

				if doRedoCheck {
					data.redoCheck("redocheck", idx, d, ev)
				}
			}

			if doErrCheck && !didErrCheck {
				t.Error("error check has never ran")
			}
		}()
	}
}
