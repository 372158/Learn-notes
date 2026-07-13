# DAY 2

# 1.泛型

链接：[Go语言泛型 - GoLang文档](https://www.golangdev.cn/zh-cn/advance/generic.html)

结构体泛型写法：

~~~go
type Company[T int | string, S int | string] struct {
    Name	 string
    Id 		 T
    Stuff []S
}
~~~



SayAble 是一个泛型接口， Person 实现了该接口。

~~~go
type SayAble[T int | string] interface {
    Say() T
}

type Person[T int | string] struct {
	msg T
}

func (p Person[T]) Say() T {
    return p.msg
}

func main() {
	var s SayAble[string]
    s = Person[string]{"hello world"}
    fmt.Println(s.Say())
}
~~~



# 2.迭代器

链接：[Go语言迭代器 - GoLang文档](https://www.golangdev.cn/zh-cn/advance/iterator.html)

闭包求解斐波那契数列的例子：

~~~go
func Fibonacci(n int) func() (int, bool){
    a, b, c := 1, 1, 2
    i := 0
    return func() (int, bool) {
        if i >= n {
		return 0, false
        } else if i < 2 {
            f := i
            i++
            return f, true
        }
        
        a, b = b, c
        c = a + b
        i++
        return a, true
}
}
~~~

可以改造成迭代器,如下：

~~~go
func Fibonacci(n int) func(yield func(int) bool) {
    a ,b ,c := 0, 1, 1
    return func(yield func(int) bool) {
        for range n {
            if !yield(a) {
			return 
            }
            a, b = b , c
            c = a + b
        }
    }
}

func mian() {
    n := 8
    for f := range Fibonacci(n) {
        fmt.Println(f)
    }
}
~~~

##推送式迭代器

~~~go
for f := range Fibonacci(n) {
    fmt.Println(f)
}
~~~

这个例子其实就相当于，如下：

~~~go
Fibonacco(n) (func(f int) bool {
    fmt.Println(f)
    return true
}) 
~~~

循环体的 body 就是迭代器的回调函数`yiled`，当函数返回`true`迭代器会继续迭代，否则就会停止

~~~go
for index, value := range iterator() {
    fmt.Println(index, value)}
~~~

在Go中的表现形式由`range`返回被迭代的元素。

## 拉取式迭代器

推送式迭代器（pushing iterator）是由迭代器来控制迭代的逻辑，用户被动获取元素，相反的拉取式迭代器（pulling iterator）就是由用户来控制迭代逻辑，主动的去获取序列元素。一般而言，拉取式迭代器都会有特定的函数如`next()`，`stop()`来控制迭代的开始或结束，它可以是一个闭包或者结构体。

~~~go
scanner := bufio.NewScanner(file)
for sacnner.Scan() {
    line, err := scanner.Text(), scanner.Err()
    if err != nil {
        fmt.Println(err)
        return 
    }
    fmt.Println(line)
}
~~~

如上所示，Scanner 通过方法`Text()`来获取文件中的下一行文本，通过方法`Scan()`来表示迭代是否结束，这也是拉取式迭代器的一种模式。Scanner 采用结构体来记录状态，而在`iter`库定义的拉取式迭代器采用闭包来记录状态，我们通过`iter.Pull`或`iter.Pull2`函数就可以将一个标准的推送式迭代器转换为拉取式迭代器，`iter.Pull`与`iter.Pull2`的区别就是后者的返回值有两个



把之前的斐波那契迭代器改造成拉取式迭代器，如下:

~~~go
func main() {
    n := 10
    next , stop := iter.Pull(Fibonacci(n))
    defer stop()
    for {
        fibn, ok := next()
        if !ok {
            break
        }
        fmt.Println(fibn)
    }
}
~~~

## 标准库

有很多标准库也支持了迭代器，最常用的就是`slices`和`maps`标准库，下面介绍几个比较实用的功能。

### slices.All

~~~go
func ALL[Slice ~[]E, E any](s Slice) iter.Seq2[int, E]
~~~

`slices.All`会将切片转换成一个切片迭代器

~~~
func main() {
    s := []int {1, 2, 3, 4, 5}
    for i, n := range slices.All(s) {
        fmt.Println(i, n)
    }
}
~~~

输出

~~~
0 1
1 2
2 3
3 4
4 5

~~~

### slices.Values

~~~go
func Values[Slice ~[]E, E any](s Slice) iter.Seq[E]
~~~

`slices.Values`会将切片转换成一个切片迭代器，但是不带索引

~~~go
 func main() {
	s := []int {1, 2, 3, 4, 5}
	for n := range  slices.Values(s) {
        fmt.Println
	}
}
~~~

输出

~~~go
1
2
3
4
5
~~~

### slices.Chunk

~~~go
func Chunk[Slice ~[]E, E any](s Slice,n int) iter.Seq[Slice]
~~~

`slices.Chunk`函数会返回一个迭代器，该迭代器会以 n 个元素为切片推送给调用者

~~~go
func main() {
    s := []int {1, 2, 3, 4, 5}
    for chunk := range slices.Chunk(s, 2) {
        fmt.Println(chunk)
    }
}
~~~

输出

~~~go
[1 2]
[3 4]
[5]
~~~

###slices.Collect

```
func Collect[E any](seq iter.Seq[E]) []E
```

`slices.Collect`函数会将切片迭代器收集成一个切片

```
func main() {
  s := []int{1, 2, 3, 4, 5}
  s2 := slices.Collect(slices.Values(s))
  fmt.Println(s2)
}
```

输出

```
[1 2 3 4 5]
```

**maps.Keys**

```
func Keys[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[K]
```

`maps.Keys`会返回一个迭代 map 所有键的迭代器，配合`slices.Collect`可以直接收集成一个切片。

```
func main() {
  m := map[string]int{"one": 1, "two": 2, "three": 3}
  keys := slices.Collect(maps.Keys(m))
  fmt.Println(keys)
}
```

输出

```
[three one two]
```

由于 map 是无序的，所以输出也不固定

###maps.Values

```
func Values[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[V]
```

`maps.Values`会返回一个迭代 map 所有值的迭代器，配合`slices.Collect`可以直接收集成一个切片。

```
func main() {
  m := map[string]int{"one": 1, "two": 2, "three": 3}
  keys := slices.Collect(maps.Values(m))
  fmt.Println(keys)
}
```

输出

```
[3 1 2]
```

由于 map 是无序的，所以输出也不固定

###maps.All

```
func All[Map ~map[K]V, K comparable, V any](m Map) iter.Seq2[K, V]
```

`maps.All`可以将一个 map 转换为成一个 map 迭代器

```
func main() {
  m := map[string]int{"one": 1, "two": 2, "three": 3}
  for k, v := range maps.All(m) {
    fmt.Println(k, v)
  }
}
```

一般不会这么直接用，都是拿来配合其他数据流处理函数的。

###maps.Collect

```
func Collect[K comparable, V any](seq iter.Seq2[K, V]) map[K]V
```

`maps.Collect`可以将一个 map 迭代器收集成一个 map

```
func main() {
  m := map[string]int{"one": 1, "two": 2, "three": 3}
  m2 := maps.Collect(maps.All(m))
  fmt.Println(m2)
}
```

collect 函数一般作为数据流处理的终结函数来使用。

## 链式调用

通过上面标准库提供的函数，我们可以将其组合来处理数据流，比如对数据流进行排序，如下

```
sortedSlices := slices.Sorted(slices.Values(s))
```

go 的迭代器采用的是闭包，只能像这样嵌套函数调用，本身没法链式调用，调用链长了以后可读性会很差，但我们可以自己通过结构体来记录迭代器，就能够实现链式调用。

### demo

一个简单的链式调用 demo 如下所示，它包含了`Filter`，`Map`，`Find`，`Some`等常用的功能。

```
package iterx

import (
  "iter"
  "slices"
)

type SliceSeq[E any] struct {
  seq iter.Seq2[int, E]
}

func (s SliceSeq[E]) All() iter.Seq2[int, E] {
  return s.seq
}

func (s SliceSeq[E]) Filter(filter func(int, E) bool) SliceSeq[E] {
  return SliceSeq[E]{
    seq: func(yield func(int, E) bool) {
      // 重新组织索引
      i := 0
      for k, v := range s.seq {
        if filter(k, v) {
          if !yield(i, v) {
            return
          }
          i++
        }
      }
    },
  }
}

func (s SliceSeq[E]) Map(mapFn func(E) E) SliceSeq[E] {
  return SliceSeq[E]{
    seq: func(yield func(int, E) bool) {
      for k, v := range s.seq {
        if !yield(k, mapFn(v)) {
          return
        }
      }
    },
  }
}

func (s SliceSeq[E]) Fill(fill E) SliceSeq[E] {
  return SliceSeq[E]{
    seq: func(yield func(int, E) bool) {
      for i, _ := range s.seq {
        if !yield(i, fill) {
          return
        }
      }
    },
  }
}

func (s SliceSeq[E]) Find(equal func(int, E) bool) (_ E) {
  for i, v := range s.seq {
    if equal(i, v) {
      return v
    }
  }
  return
}

func (s SliceSeq[E]) Some(match func(int, E) bool) bool {
  for i, v := range s.seq {
    if match(i, v) {
      return true
    }
  }
  return false
}

func (s SliceSeq[E]) Every(match func(int, E) bool) bool {
  for i, v := range s.seq {
    if !match(i, v) {
      return false
    }
  }
  return true
}

func (s SliceSeq[E]) Collect() []E {
  var res []E
  for _, v := range s.seq {
    res = append(res, v)
  }
  return res
}

func (s SliceSeq[E]) Sort(cmp func(x, y E) int) []E {
  collect := s.Collect()
  slices.SortFunc(collect, cmp)
  return collect
}

func (s SliceSeq[E]) SortStable(cmp func(x, y E) int) []E {
  collect := s.Collect()
  slices.SortStableFunc(collect, cmp)
  return collect
}

func Slice[S ~[]E, E any](s S) SliceSeq[E] {
  return SliceSeq[E]{seq: slices.All(s)}
}
```

然后我们就可以通过链式调用来处理了，看几个使用案例。

**处理元素值**

```
func main() {
  s := []string{"apple", "banana", "cherry"}
  all := iterx.Slice(s).Map(strings.ToUpper).All()
  for i, v := range all {
    fmt.Printf("index: %d, value: %s\n", i, v)
  }
}
```

输出

```
index: 0, value: APPLE
index: 1, value: BANANA
index: 2, value: CHERRY
```

**寻找某一个指定值**

```
func main() {
  s := []int{1, 2, 3, 4, 5}
  result := iterx.Slice(s).Find(func(i int, e int) bool {
    return e == 3
  })
  fmt.Println(result)
}
```

输出

```
3
```

**填充切片**

```
func main() {
  s := []int{1, 2, 3, 4, 5}
  result := iterx.Slice(s).Fill(6).Collect()
  fmt.Println(result)
}
```

输出

```
[6 6 6 6 6]
```

**过滤元素**

```
func main() {
  s := []int{1, 2, 3, 4, 5}
  filter := iterx.Slice(s).Filter(func(i int, e int) bool {
    return e%2 == 0
  }).All()
  for i, v := range filter {
    fmt.Printf("Index: %d, Value: %d\n", i, v)
  }
}
```

输出

```
Index: 0, Value: 2
Index: 1, Value: 4
```

比较可惜的是 Go 目前还不支持简写匿名函数，就像 js，rust，java 中的箭头函数一样，否则链式调用还可以更加简洁和优雅一些。