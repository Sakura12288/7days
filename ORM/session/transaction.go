package session

import "7days/ORM/logg"

func (s *Session) Begin() (err error) {
	logg.Info("事务开启")
	if s.tx, err = s.db.Begin(); err != nil {
		logg.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	logg.Info("事务提交")
	if err = s.tx.Commit(); err != nil {
		logg.Error(err)
		return
	}
	return
}

func (s *Session) RollBack() (err error) {
	logg.Info("事务回滚")
	if err = s.tx.Rollback(); err != nil {
		logg.Error(err)
		return
	}
	return
}
