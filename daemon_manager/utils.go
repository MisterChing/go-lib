package daemon_manager

import "errors"

func FetchMessage(ch <-chan interface{}, batchSize int) []interface{} {
	var (
		ret = make([]interface{}, 0, batchSize)
	)
	for i := 0; i < batchSize; i++ {
		if msg := readChanNoBlock(ch); msg != nil {
			ret = append(ret, msg)
		}
	}
	return ret
}

func readChanNoBlock(ch <-chan interface{}) interface{} {
	select {
	case v, ok := <-ch:
		if ok {
			return v
		}
		return nil
	default:
		return nil
	}
}

func writeChanNoBlock(ch chan<- interface{}, data interface{}) error {
	select {
	case ch <- data:
		return nil
	default:
		return errors.New("write error")
	}
}
