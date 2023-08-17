package utils

// BytesEquals check whether the given bytes is equal
func BytesEquals(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	size := len(a)
	for i := 0; i < size; i++ {
		av := a[i]
		bv := b[i]
		if av != bv {
			return false
		}
	}
	return true
}

// ToCmdLine  字符串转切片
func ToCmdLine(cmd ...string) [][]byte {
	args := make([][]byte, len(cmd))
	for i, s := range cmd {
		args[i] = []byte(s)
	}
	return args
}

// ToCmdLine2 将 set key val 的 commandName(set)和后面的参数 key val 组装成一个 [][]byte
// ...[]byte表示接收不定量的 []byte
func ToCmdLine2(commandName string, args ...[]byte) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(commandName)
	for i, arg := range args {
		result[i+1] = arg
	}
	return result
}
