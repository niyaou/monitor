package gnet

// 环形buffer,专为TcpConnection定制,在收发包时,可以减少内存分配和拷贝
// NOTE:不支持多线程,不具备通用性
type RingBuffer struct {
	// 数据
	buffer []byte
	// 写位置
	w int
	// 读位置
	r int
}

// 指定大小的RingBuffer,不支持动态扩容
func NewRingBuffer(size int) *RingBuffer {
	if size <= 0 {
		return nil
	}
	return &RingBuffer{
		buffer: make([]byte, size),
	}
}

// Buffer
func (_self *RingBuffer) GetBuffer() []byte {
	return _self.buffer
}

func (_self *RingBuffer) Size() int {
	return len(_self.buffer)
}

// 未被读取的长度
func (_self *RingBuffer) UnReadLength() int {
	// logger.Info("---->UnRead-----.w %v  .r %v  ---UnRead-----,%v", _self.w, _self.r, _self.w-_self.r)
	return _self.w - _self.r
}

// 读取指定长度的数据
func (_self *RingBuffer) ReadFull(readLen int) []byte {
	// logger.Info("---ReadFull-----readLen,%v", readLen)
	if _self.UnReadLength() < readLen {
		return nil
	}
	readBuffer := _self.ReadBuffer()
	if len(readBuffer) >= readLen {
		// 数据连续,不产生copy
		// logger.Info("数据连续,不产生copy %v", readLen)
		_self.SetReaded(readLen)
		return readBuffer[0:readLen]
	} else {
		// 数据非连续,需要重新分配数组,并进行2次拷贝
		// logger.Info(" 数据非连续,需要重新分配数组,并进行2次拷贝")
		data := make([]byte, readLen)
		// 先拷贝RingBuffer的尾部
		n := copy(data, readBuffer)
		// 再拷贝RingBuffer的头部
		copy(data[n:], _self.buffer)
		// logger.Info("非连续,需要重新分配数组,并进行2次拷贝 %v", readLen)
		_self.SetReaded(readLen)
		return data
	}
}

//// 读位置
//func (_self *RingBuffer) ReadIndex() int {
//	return _self.r%len(_self.buffer)
//}

// 设置已读取长度
func (_self *RingBuffer) SetReaded(readedLength int) {
	_self.r += readedLength
	// logger.Info("-->SetReaded--- %v --.w %v  .r %v  ---UnRead-----,%v", readedLength, _self.w, _self.r, _self.w-_self.r)

}

// 返回可读取的连续buffer(不产生copy)
// NOTE:调用ReadBuffer之前,需要先确保UnReadLength()>0
func (_self *RingBuffer) ReadBuffer() []byte {
	writeIndex := _self.w % len(_self.buffer)
	readIndex := _self.r % len(_self.buffer)
	if readIndex < writeIndex {
		// [_______r.....w_____]
		//         <- n ->
		// 可读部分是连续的
		return _self.buffer[readIndex : readIndex+_self.w-_self.r]
	} else {
		// [........w_______r........]
		//                  <-  n1  ->
		// <-  n2  ->
		// 可读部分被分割成尾部和头部两部分,先返回尾部那部分
		return _self.buffer[readIndex:]
	}
}

// 返回可写入的连续buffer
func (_self *RingBuffer) WriteBuffer() []byte {
	writeIndex := _self.w % len(_self.buffer)
	readIndex := _self.r % len(_self.buffer)
	if readIndex < writeIndex {
		// [_______r.....w_____]
		// 可写部分被成尾部和头部两部分,先返回尾部那部分
		return _self.buffer[writeIndex:]
	} else if readIndex > writeIndex {
		// [........w_______r........]
		// 可写部分是连续的
		return _self.buffer[writeIndex:readIndex]
	} else {
		if _self.r == _self.w {
			return _self.buffer[writeIndex:]
		}
		return nil
	}
}

// 设置已写入长度
func (_self *RingBuffer) SetWrited(writedLength int) {
	_self.w += writedLength
	// logger.Info("---->SetWrited----%v --.w %v  .r %v  --UnRead-----,%v", writedLength, _self.w, _self.r, _self.w-_self.r)
}

// 写入数据
func (_self *RingBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	bufferSize := len(_self.buffer)
	canWriteSize := bufferSize + _self.r - _self.w
	if canWriteSize <= 0 {
		return 0, ErrBufferFull
	}
	writeIndex := _self.w % bufferSize
	// 有足够的空间可以把p写完
	if canWriteSize >= len(p) {
		n = copy(_self.buffer[writeIndex:], p)
		// 如果没能一次写完,说明写在尾部了,剩下的直接写入头部
		if n < len(p) {
			n += copy(_self.buffer[0:], p[n:])
		}
		_self.w += n
		// logger.Info("---->有足够的空间可以把p写完--canWriteSize%v---.w %v  .r %v  --UnRead-----,%v", canWriteSize, _self.w, _self.r, _self.w-_self.r)
		return
	} else {
		n = copy(_self.buffer[writeIndex:], p[0:canWriteSize])
		// 如果没能一次写完,说明写在尾部了,剩下的直接写入头部
		if n < canWriteSize {
			n += copy(_self.buffer[0:], p[n:canWriteSize])
		}
		_self.w += n
		// logger.Info("---->如果没能一次写完,说明写在尾部了,剩下的直接写入头部-----.w %v  .r %v  --UnRead-----,%v", _self.w, _self.r, _self.w-_self.r)
	}
	return
}
