package common

import (
	"gokit/pkg/common/assert"
	"testing"
)

func TestStartWith(t *testing.T) {
	assert.IsTrue(StartWith("abc", "abc"), "不通过")
	assert.IsTrue(!StartWith("a", "ab"), "不通过")
	assert.IsTrue(StartWith("abcdefg", "abcd"), "不通过")
	assert.IsTrue(StartWith("我是谁", "我"), "不通过")
}
