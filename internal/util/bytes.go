package util

import "bytes"

// TODO: RmBSpace 去除 p 中全部空白字符
func RmBSpace(p []byte) []byte {
	if bytes.IndexByte(p, '\n') == -1 {
		return p
	}

	return p
}
