package process

type Process struct {
	Name 			string
	Class 			string
	ExecuteMethod 	string
	Params 			map[string]interface{}
}