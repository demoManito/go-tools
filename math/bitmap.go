package math

// BitMap bit map interface
type BitMap interface {
	Set(int)
	UnSet(int)
	IsSet(int) bool
}

var _ BitMap = new(BitMap64)
var _ BitMap = new(BitMap16)
var _ BitMap = new(BitMap32)
var _ BitMap = new(BitMap64)

// BitMap bit 位
type (
	BitMap8  int8  // 最大支持 8 位, 最大支持 8 位,   0000 0000
	BitMap16 int16 // 最大支持 16 位, 最大支持 16 位, 0000 0000 0000 0000
	BitMap32 int32 // 最大支持 32 位, 最大支持 32 位, 0000 0000 0000 0000 0000 0000 0000 0000
	BitMap64 int64 // 最大支持 64 位, 最大支持 64 位, 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000
)

// Set 设置第 i 位
func (b8 *BitMap8) Set(i int) {
	*b8 |= 1 << (i - 1)
}

// UnSet 取消第 i 位
// 0&^0 = 0; 1&^0=1; 0&^1=0; 1&^1=0
func (b8 *BitMap8) UnSet(i int) {
	*b8 &^= 1 << (i - 1)
}

// IsSet 第 i 位是否存在, 存在返回 true, 不存在 false
func (b8 *BitMap8) IsSet(i int) bool {
	return *b8&(1<<(i-1)) != 0
}

func (b16 *BitMap16) Set(i int) {
	*b16 |= 1 << (i - 1)
}

func (b16 *BitMap16) UnSet(i int) {
	*b16 &^= 1 << (i - 1)
}

func (b16 *BitMap16) IsSet(i int) bool {
	return *b16&(1<<(i-1)) != 0
}

func (b32 *BitMap32) Set(i int) {
	*b32 |= 1 << (i - 1)
}

func (b32 *BitMap32) UnSet(i int) {
	*b32 &^= 1 << (i - 1)
}

func (b32 *BitMap32) IsSet(i int) bool {
	return *b32&(1<<(i-1)) != 0
}

func (b64 *BitMap64) Set(i int) {
	*b64 |= 1 << (i - 1)
}

func (b64 *BitMap64) UnSet(i int) {
	*b64 &^= 1 << (i - 1)
}

func (b64 *BitMap64) IsSet(i int) bool {
	return *b64&(1<<(i-1)) != 0
}
