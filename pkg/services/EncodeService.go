package services

type encodeService struct {
}

func NewEncodeService() EncodeService {
	return &encodeService{}
}

type EncodeService interface{}

func (c *encodeService) EncodeSha256(content []byte) {

}
