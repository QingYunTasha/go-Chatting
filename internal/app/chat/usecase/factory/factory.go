package factory

import (
	orm "github.com/QingYunTasha/go-Chatting/internal/infra/orm/factory"
)

type ChatUsecaseRepository struct {
	OrmRepository *orm.OrmRepository
}
