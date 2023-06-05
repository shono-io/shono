package commons

type KeySection struct {
	Kind string
	Code string
}

type Key []KeySection

func NewKey(kind, code string) Key {
	return Key{KeySection{Kind: kind, Code: code}}
}

func (k Key) Parent() Key {
	if len(k) == 0 {
		return nil
	}

	return (k)[:len(k)-1]
}

func (k Key) Child(kind, code string) Key {
	return append(k, KeySection{Kind: kind, Code: code})
}

func (k Key) Code() string {
	if len(k) == 0 {
		return ""
	}

	return k[len(k)-1].Code
}

func (k Key) Kind() string {
	if len(k) == 0 {
		return ""
	}

	return k[len(k)-1].Kind
}

func (k Key) String() string {
	if len(k) == 0 {
		return ""
	}

	var s string
	for _, section := range k {
		s += section.Kind + "_" + section.Code + ":"
	}
	return s[:len(s)-1]
}

func (k Key) CodeString() string {
	if len(k) == 0 {
		return ""
	}

	var s string
	for _, section := range k {
		s += section.Code + "__"
	}
	return s[:len(s)-2]
}
