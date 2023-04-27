package geecache

// A ByteView holds an immutable view of bytes(read-only)
// 只读数据结构: 表示缓存值
type ByteView struct {
	// b 存储真实的缓存值:支持任意的数据类型的存储,例如 字符串,图片...
	b []byte
}

// Len returns the view's length
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice
func (v ByteView) ByteSlice() []byte {
	// b 为只读,返回一个拷贝防止缓存值被外部程序修改
	// bytes 为切片,是引用类型,不会深拷贝(增加一个指针并申请一块新的内存)
	return cloneBytes(v.b)
}

// String return the data as a string, making a copy if necessary
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}