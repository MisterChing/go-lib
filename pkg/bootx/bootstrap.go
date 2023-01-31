package bootx

type BeforeServerFunc func() error
type AfterServerFunc func()

type Bootloader struct {
	BeforeServerChain []BeforeServerFunc
	AfterServerChain  []AfterServerFunc
}

func NewBootStrap() *Bootloader {
	return new(Bootloader)
}

func (boot *Bootloader) AddBeforeServerFunc(fns ...BeforeServerFunc) {
	for _, fn := range fns {
		boot.BeforeServerChain = append(boot.BeforeServerChain, fn)
	}
}

func (boot *Bootloader) AddAfterServerFunc(fns ...AfterServerFunc) {
	for _, fn := range fns {
		boot.AfterServerChain = append(boot.AfterServerChain, fn)
	}
}

func (boot *Bootloader) SetUp() error {
	for _, fn := range boot.BeforeServerChain {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (boot *Bootloader) Destroy() {
	for _, fn := range boot.AfterServerChain {
		fn()
	}
}
