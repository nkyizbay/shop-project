package item

type Service interface {
	CreateItem(item *Item) error
	CancelItem(id int) error
}

type defaultService struct {
	itemRepo Repository
}

func NewItemService(itemRepo Repository) Service {
	return &defaultService{itemRepo: itemRepo}
}

func (s *defaultService) CreateItem(i *Item) error {
	if err := s.itemRepo.Create(i); err != nil {
		return err
	}

	return nil
}

func (s *defaultService) CancelItem(id int) error {
	if err := s.itemRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
