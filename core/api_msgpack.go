package core

func MarshalMsgPack(in any) ([]byte, error) {
	return globalProvider.MarshalMsgPack(in)
}

func UnmarshalMsgPack(out any, b []byte) error {
	return globalProvider.UnmarshalMsgPack(out, b)
}
