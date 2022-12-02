package transport

type RequestOption func (request *Request)




func NewRequest(service, method string) *Request {
	return &Request{
		ServiceName: service,
		MethodName: method,
	}
}


type SharedHeader struct {
	HeadLength uint32
	BodyLength uint32

	MessageID uint32

	Version     byte
	CompressionFormat byte
	SerializationFormat byte
}


type requestHeader struct {
	SharedHeader
	meta map[string]string
}

type Request struct {
	SharedHeader

	ServiceName string
	MethodName  string
	// todo: ctx is ignored
	Arg any
	Meta map[string]string
	Payload []byte
}

type Response struct {
	SharedHeader

	Error   []byte
	// TODO : 业务Error BizError []byte
	Payload []byte
}
