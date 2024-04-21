package core

func RegisterType(name string, id any) error {
	return globalProvider.RegisterType(name, id)
}
