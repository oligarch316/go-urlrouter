package graph

import "fmt"

type Key interface {
	fmt.Stringer
	sealedKey()
}

type (
	KeyConstant  string
	KeyParameter string
	KeyWildcard  struct{}
)

func (KeyConstant) sealedKey()  {}
func (KeyParameter) sealedKey() {}
func (KeyWildcard) sealedKey()  {}

func (KeyWildcard) String() string     { return "wild" }
func (kc KeyConstant) String() string  { return fmt.Sprintf("const(%s)", string(kc)) }
func (kp KeyParameter) String() string { return fmt.Sprintf("param(%s)", string(kp)) }
