package session

//存放实现了的钩子
// Hooks constants

type BeforeQuery interface {
	BeforeQuery(s *Session) error
}
type AfterQuery interface {
	AfterQuery(s *Session) error
}
type BeforeUpdate interface {
	BeforeUpdate(s *Session) error
}
type AfterUpdate interface {
	AfterUpdate(s *Session) error
}
type AfterInsert interface {
	AfterInsert(s *Session) error
}
type BeforeInsert interface {
	BeforeInsert(s *Session) error
}
type BeforeDelete interface {
	BeforeDelete(s *Session) error
}
type AfterDelete interface {
	AfterDelete(s *Session) error
}
