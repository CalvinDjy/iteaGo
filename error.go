package itea

type Error interface {
	GetCode() 	int
	GetError() 	string
}

type BusinessError struct {
	Code 		int
	Error 		string
}

func (be *BusinessError) GetCode() int {
	return be.Code
}

func (be *BusinessError) GetError() string {
	return be.Error
}

type DatabaseError struct {
	Code 		int
	Error 		string
}

func (be *DatabaseError) GetCode() int {
	return be.Code
}

func (be *DatabaseError) GetError() string {
	return be.Error
}

type ServerError struct {
	Code 		int
	Error 		string
}

func (be *ServerError) GetCode() int {
	return be.Code
}

func (be *ServerError) GetError() string {
	return be.Error
}

type ParamsError struct {
	Code 		int
	Error 		string
}

func (be *ParamsError) GetCode() int {
	return be.Code
}

func (be *ParamsError) GetError() string {
	return be.Error
}
