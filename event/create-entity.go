package event

import (
	"context"
	"fmt"

	"github.com/ajainc/chain/ctx/docdb"
	"github.com/ajainc/chain/ctx/docdb/dblogic"
	"github.com/ajainc/chain/grpc/gen"
	"github.com/fatih/structs"
)

type CreateEntity struct {
	Entity gen.Entity
}

func (ev CreateEntity) Description() string {
	return "Create Entity"
}

func (ev CreateEntity) Do(c context.Context) error {
	coll := docdb.COLL_ENTITY

	db := docdb.FromContext(c)

	if dblogic.ObjectIDExists(db, coll, ev.Entity.ObjectID) {
		return fmt.Errorf("entity  [%s] already exists", ev.Entity.ObjectID)
	}

	entities := db.Use(coll)
	_, err := entities.Insert(structs.Map(ev.Entity))

	return err
}

func (ev CreateEntity) Undo(c context.Context) error {
	coll := docdb.COLL_ENTITY

	db := docdb.FromContext(c)

	id, _, err := dblogic.GetByObjectID(db, coll, ev.Entity.ObjectID)
	if err != nil {
		return err
	}
	entities := db.Use(coll)
	return entities.Delete(id)

}
