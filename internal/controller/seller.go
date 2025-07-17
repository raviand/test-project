package controller

// get

// post

type SellerStruct struct {
	Name string
}

func (s *SellerStruct) CreateSeller(name *string) error {
	s.Name = *name
	// conn -> commit trx
	return nil
}
