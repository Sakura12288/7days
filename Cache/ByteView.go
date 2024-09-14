package Cache

//可以理解成表面的cache，对底层Cache的内容进行封装，防止被修改

type ByteView struct {
	b []byte
}

//定义Len()方法，实现Value接口，以作为我们的Cache的值

func (v ByteView) Len() int {
	return len(v.b)
}

//实现String方法，方便直接输出b所对应的内容

func (v ByteView) String() string {
	return string((v.b))
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte { //深拷贝，防止底层数据被修改
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
