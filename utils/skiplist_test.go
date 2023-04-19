package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
该函数的输入参数是一个整数len，表示要生成的随机字符串的长度。函数返回一个字符串类型的随机字符串。
函数的实现过程如下：
创建一个指定长度的字节数组bytes，用于存储生成的随机字符。
使用for循环生成指定长度的随机字符。在每次循环中，生成一个0到25之间的随机整数，并加上65，得到一个ASCII码值。这个值表示大写字母A到Z中的一个字符。然后将这个字符转换成一个字节，并存储到bytes数组中。
将bytes数组转换成字符串类型，并返回该字符串。
需要注意的是，RandString函数中使用了一个全局变量r，它是一个随机数生成器。在该代码中，没有详细说明r是如何初始化的。通常情况下，为了保证随机数的质量和随机性，应该使用一个种子来初始化随机数生成器。
*/

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

/**
函数名：TestSkipListBasicCRUD

输入参数：t *testing.T，是Go语言中的测试框架参数，用于进行测试断言和错误报告。

输出结果：该函数没有明确的返回值，但是它将执行一系列的CRUD操作，并使用断言（assert）方法判断每次操作的结果是否符合预期，以便开发人员进行调试和验证。

函数实现：
该函数首先创建了一个容量为1000的SkipList对象，用于执行测试操作。

接下来，函数分别执行以下操作：

第1个操作：创建一个随机的Entry对象entry1，包括一个10位的随机字符串作为key和"Val1"作为value，然后将其插入到SkipList中，并使用assert.Equal方法判断SkipList中搜索到的entry1的value是否等于"Val1"。
第2个操作：创建一个随机的Entry对象entry2，包括一个10位的随机字符串作为key和"Val2"作为value，然后将其插入到SkipList中，并使用assert.Equal方法判断SkipList中搜索到的entry2的value是否等于"Val2"。
第3个操作：使用一个不存在的key（随机字符串）进行搜索，使用assert.Nil方法判断返回值是否为nil。
第4个操作：创建一个新的Entry对象entry2_new，包括entry1的key和"Val1+1"作为value，然后将其插入到SkipList中，并使用assert.Equal方法判断SkipList中搜索到的entry2_new的value是否等于"Val1+1"。
*/

func TestSkipListBasicCRUD(t *testing.T) {
	list := NewSkiplist(1000)

	//Put & Get
	entry1 := NewEntry([]byte(RandString(10)), []byte("Val1"))
	list.Add(entry1)
	vs := list.Search(entry1.Key)
	assert.Equal(t, entry1.Value, vs.Value)

	entry2 := NewEntry([]byte(RandString(10)), []byte("Val2"))
	list.Add(entry2)
	vs = list.Search(entry2.Key)
	assert.Equal(t, entry2.Value, vs.Value)

	//Get a not exist entry
	assert.Nil(t, list.Search([]byte(RandString(10))).Value)

	//Update a entry
	entry2_new := NewEntry(entry1.Key, []byte("Val1+1"))
	list.Add(entry2_new)
	assert.Equal(t, entry2_new.Value, list.Search(entry2_new.Key).Value)
}

func Benchmark_SkipListBasicCRUD(b *testing.B) {
	list := NewSkiplist(100000000)
	key, val := "", ""
	maxTime := 1000
	for i := 0; i < maxTime; i++ {
		//number := rand.Intn(10000)
		key, val = RandString(10), fmt.Sprintf("Val%d", i)
		entry := NewEntry([]byte(key), []byte(val))
		list.Add(entry)
		searchVal := list.Search([]byte(key))
		assert.Equal(b, searchVal.Value, []byte(val))
	}
}

func TestDrawList(t *testing.T) {
	list := NewSkiplist(1000)
	n := 12
	for i := 0; i < n; i++ {
		index := strconv.Itoa(r.Intn(90) + 10)
		key := index + RandString(8)
		entryRand := NewEntry([]byte(key), []byte(index))
		list.Add(entryRand)
	}
	list.Draw(true)
	fmt.Println(strings.Repeat("*", 30) + "分割线" + strings.Repeat("*", 30))
	list.Draw(false)
}

func TestConcurrentBasic(t *testing.T) {
	const n = 1000
	l := NewSkiplist(100000000)
	var wg sync.WaitGroup
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("Keykeykey%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Add(NewEntry(key(i), key(i)))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			require.EqualValues(t, key(i), v.Value)
			return

			require.Nil(t, v)
		}(i)
	}
	wg.Wait()
}

func Benchmark_ConcurrentBasic(b *testing.B) {
	const n = 1000
	l := NewSkiplist(100000000)
	var wg sync.WaitGroup
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("keykeykey%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Add(NewEntry(key(i), key(i)))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			require.EqualValues(b, key(i), v.Value)
			require.NotNil(b, v)
		}(i)
	}
	wg.Wait()
}

func TestSkipListIterator(t *testing.T) {
	list := NewSkiplist(100000)

	//Put & Get
	entry1 := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry1)
	assert.Equal(t, entry1.Value, list.Search(entry1.Key).Value)

	entry2 := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry2)
	assert.Equal(t, entry2.Value, list.Search(entry2.Key).Value)

	//Update a entry
	entry2_new := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry2_new)
	assert.Equal(t, entry2_new.Value, list.Search(entry2_new.Key).Value)

	iter := list.NewSkipListIterator()
	for iter.Rewind(); iter.Valid(); iter.Next() {
		fmt.Printf("iter key %s, value %s", iter.Item().Entry().Key, iter.Item().Entry().Value)
	}
}
