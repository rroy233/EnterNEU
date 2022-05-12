package tgService

import (
	"strings"
	"testing"
)

func TestGetPartIndex1(t *testing.T) {
	offset := 6
	length := 3
	if o, l := getPartIndex("糊啦啦啦户啦户啦啦啦", "户啦啦"); o != offset || l != length {
		t.Errorf("错误output(%d,%d)!=(%d,%d)\n", o, l, offset, length)
	}
	t.Log("ok")
}

func TestGetPartIndex2(t *testing.T) {
	offset := 11
	length := 2
	if o, l := getPartIndex("abbdjhbhd\nss2s2ss2", "s2"); o != offset || l != length {
		t.Errorf("错误output(%d,%d)!=(%d,%d)\n", o, l, offset, length)
	}
	t.Log("ok")
}

func TestGetPartIndex3(t *testing.T) {
	offset := -1
	length := 2
	if o, l := getPartIndex("aaaaaaaaaaa", "s2"); o != offset || l != length {
		t.Errorf("错误output(%d,%d)!=(%d,%d)\n", o, l, offset, length)
	}
	t.Log("ok")
}

func TestGetPartIndex4(t *testing.T) {
	offset := 100000 * 3
	length := 3
	if o, l := getPartIndex(strings.Repeat("ss2", 100000)+"s2a", "s2a"); o != offset || l != length {
		t.Errorf("错误output(%d,%d)!=(%d,%d)\n", o, l, offset, length)
	}
	t.Log("ok")
}

func TestGetPartIndex5(t *testing.T) {
	offset := 36
	length := 10
	if o, l := getPartIndex("该功能允许所有已授权的用户使用\n授权原则:仅允许授权给东大学生!!!\n 请务必严格验证其身份。\n\n请按照下面的格式进行发送:\n /add <UID> (如/add 123)", "请务必严格验证其身份"); o != offset || l != length {
		t.Errorf("错误output(%d,%d)!=(%d,%d)\n", o, l, offset, length)
	}
	t.Log("ok")
}
