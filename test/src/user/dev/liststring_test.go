package liststring

import (
	"testing"
	"fmt"
	"github.com/bmizerany/assert"
)
/////////////////////////
func TestListString(t *testing.T){
	// 值
	var lst ListString
	lst.Append("hello")
	fmt.Println("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)
	assert.Equal(t, 1, lst.Len())

	lst.Append("hello")
	fmt.Println("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)
	assert.Equal(t, 2, lst.Len())
	// 指针
	plst := new(ListString)
	plst.Append("word")
	fmt.Println("%v (len: %d)", plst, plst.Len()) // &[2] (len: 1)
}
