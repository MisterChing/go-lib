package daemon_manager

type Queue struct {
	buf chan interface{}
}

func NewQueue(bufSize int) *Queue {
	obj := &Queue{
		buf: make(chan interface{}, bufSize),
	}
	return obj
}

func (q *Queue) Size() int {
	return len(q.buf)
}

func (q *Queue) Buf() <-chan interface{} {
	return q.buf
}

func (q *Queue) Close() {
	close(q.buf)
}

func (q *Queue) EnQueue(message interface{}) {
	q.buf <- message
}

func (q *Queue) DeQueue() interface{} {
	messages := FetchMessage(q.buf, 1)
	if len(messages) > 0 {
		return messages[0]
	}
	return nil
}

func (q *Queue) DeQueueBatch(size int) []interface{} {
	return FetchMessage(q.buf, size)
}
