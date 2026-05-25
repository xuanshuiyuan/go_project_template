package web

type WebService struct {
	Utils *Utils
}

func NewWeb() *WebService {
	web := &WebService{
		Utils: NewUtils(),
	}
	return web
}
