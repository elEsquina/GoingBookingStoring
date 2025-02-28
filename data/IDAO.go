package data

type EntityType interface {
	Book | Author | Customer | Order
}

type IDAO[T EntityType] interface {
	Create(obj T) (T, error)
	GetById(id int) (T, error)
	Update(id int, obj T) (T, error)
	Delete(id int) error
	GetAll() ([]T, error)
	//Search(criteria SearchCriteria) ([]T, error)
}


