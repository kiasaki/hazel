package data

import (
	"time"

	"github.com/kiasaki/gorp"
)

type ModelWithDates struct {
	Created int64 `created_at`
	Updated int64 `updated_at`
}

func (i *Invoice) PreInsert(s gorp.SqlExecutor) error {
	i.Created = time.Now().UnixNano()
	i.Updated = i.Created
	return nil
}

func (i *Invoice) PreUpdate(s gorp.SqlExecutor) error {
	i.Updated = time.Now().UnixNano()
	return nil
}
