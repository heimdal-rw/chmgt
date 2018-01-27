package config

import "strconv"

type boolFlag struct {
	set   bool
	value bool
}

func (b *boolFlag) Set(s string) error {
	v, err := strconv.ParseBool(s)
	b.value = v
	b.set = true
	return err
}

func (b *boolFlag) Get() interface{} { return b.value }

func (b *boolFlag) String() string { return strconv.FormatBool(b.value) }

func (b *boolFlag) IsBoolFlag() bool { return true }

type intFlag struct {
	set   bool
	value int
}

func (i *intFlag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	i.value = int(v)
	return err
}

func (i *intFlag) Get() interface{} { return i.value }

func (i *intFlag) String() string { return strconv.Itoa(int(i.value)) }

type stringFlag struct {
	set   bool
	value string
}

func (s *stringFlag) Set(val string) error {
	s.value = val
	return nil
}

func (s *stringFlag) Get() interface{} { return s.value }

func (s *stringFlag) String() string { return s.value }
