package pb

import (
	"context"
)

type ExampleService struct {
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) Say(ctx context.Context, r *SayReq) (*SayRes, error) {
	return &SayRes{Message: []byte(r.Message)}, nil
}
