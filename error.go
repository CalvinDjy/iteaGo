package itea

type Error interface {
	Code() int
	Error() string
}


type BusinessError struct {
	C int
	E string
}

func (be *BusinessError) Code() int {
	return be.C
}

func (be *BusinessError) Error() string {
	return be.E
}

type DatabaseError struct {
	C int
	E string
}

func (de *DatabaseError) Code() int {
	return de.C
}

func (de *DatabaseError) Error() string {
	return de.E
}

type ServerError struct {
	C int
	E string
}

func (se *ServerError) Code() int {
	return se.C
}

func (se *ServerError) Error() string {
	return se.E
}

type ParamsError struct {
	C int
	E string
}

func (pe *ParamsError) Code() int {
	return pe.C
}

func (pe *ParamsError) Error() string {
	return pe.E
}
