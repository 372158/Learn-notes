# 49. 字母异位词分组 (Group Anagrams)

> 难度:中等 | 分类:哈希
> 来源:[力扣 HOT 100 - Day 2](https://leetcode.cn/problems/group-anagrams/description/?envType=study-plan-v2\&envId=top-100-liked)
> 状态:✅ 已完成(学生独立解出,2026-09-03)

## 题目描述

给你一个字符串数组，请你将 字母异位词 组合在一起。可以按任意顺序返回结果列表。

**示例 1:**

**输入:** strs = \["eat", "tea", "tan", "ate", "nat", "bat"]

**输出:** \[\["bat"],\["nat","tan"],\["ate","eat","tea"]]

**解释：**

- 在 strs 中没有字符串可以通过重新排列来形成 `"bat"`。

- 字符串 `"nat"` 和 `"tan"` 是字母异位词，因为它们可以重新排列以形成彼此。

- 字符串 `"ate"` ，`"eat"` 和 `"tea"` 是字母异位词，因为它们可以重新排列以形成彼此。

**示例 2:**

**输入:** strs = \[""]

**输出:** \[\[""]]

**示例 3:**

**输入:** strs = \["a"]

**输出:** \[\["a"]]

**提示：**

- `1 <= strs.length <= 10^4`

- `0 <= strs[i].length <= 100`

- `strs[i]` 仅包含小写字母

***

## 我的思路

<!-- 先自己想,把想法写在这里;卡住了向老师要提示 -->

**(2026-09-03,学生独立想出)** 把每个字符串的字母排序(如 `"eat"/"tea"/"ate"` 排完都是 `"aet"`),排序结果相同的就是一组 → **key 存排序后的字符串**;value 存同一组的所有原字符串。

## 我的代码(第一版,Go,2026-09-03)

```go
import "sort"

func groupAnagrams(strs []string) [][]string {
    m := make(map[string][]string) // key=排序后的字符串, value=同一组的原字符串

    for _, str := range strs {
        b := []byte(str)                                          // string 不可变,先转成可排序的字节切片
        sort.Slice(b, func(i, j int) bool { return b[i] < b[j] }) // 按字母升序排

        key := string(b)             // 切片转回字符串作 key
        m[key] = append(m[key], str) // 把原字符串归入这一组
    }

    res := [][]string{}
    for _, group := range m {
        res = append(res, group)
    }
    return res
}
```

**第一版代码批改**:逻辑全对,一次通过。仅 1 处 gofmt 风格:比较器 `bool{` 大括号前要留空格 → `bool {`。另外 `make(map[string][]string)` 与 Day 1 的 `map[int]int{}` 字面量写法都合法,make 写法更常见。

**Go ↔ C++ 对照**:

| Go                                        | C++                                           | 含义             |
| ----------------------------------------- | --------------------------------------------- | -------------- |
| `b := []byte(str)` + sort                 | `sort(s.begin(), s.end())`(C++ string 可变,直接排) | string → 可排序序列 |
| `sort.Slice(b, func(i,j int) bool {...})` | `sort(v.begin(), v.end(), cmp)`               | 排序;Go 靠显式比较器   |
| `string(b)`                               | `string(v.begin(), v.end())`                  | 序列 → string    |
| `m[key] = append(m[key], str)`            | `m[key].push_back(str)`                       | 分组追加           |
| `for _, group := range m`                 | `for (auto& [k, v] : m)`                      | 收集所有组          |

## 复杂度分析

设 n = 字符串个数(≤10⁴),k = 单串最大长度(≤100):

- **时间复杂度:** O(n · k log k) —— 每个字符串排序花 O(k log k),共 n 个;map 存取 O(1)。

- **空间复杂度:** O(n · k) —— map 里最终存了所有字符串,这也是输出的必要空间。

**复杂度第二课**:Day 1 学的"数循环层数"只是起点;总时间 = **循环次数 × 每圈的活儿**,循环里"干重活"(排序)的成本要乘进总账。对比暴力法(两两比较是否互为异位词)O(n²·k):n=10⁴、k=100 时约 10¹⁰ 次操作,远超 Day 1 的"1 秒法则"(≈10⁸);排序法约 10⁴×100×7 ≈ 7×10⁶ 次,瞬间完成。

## 关键知识点(Go 新语法)

1. **string 是不可变的(immutable)**:不能原地排序,必须 `[]byte(str)` 转可变切片,排完 `string(b)` 转回去。C++ 的 `std::string` 可变可直接 sort —— 两种语言的重要差异。
2. **sort.Slice 与比较器**:`func(i, j int) bool { return b[i] < b[j] }` 是匿名函数,作为参数传给排序器,告诉它"什么叫小"。这是第一次接触"函数当值用"。Go 1.21+ 有 `slices.Sort(b)` 一行搞定,但自定义比较器必须会写,后面很多题(如合并 K 个链表)要用。
3. **append 的零值妙用**:`m[key] = append(m[key], str)` —— key 不存在时 `m[key]` 是零值 nil,而 `append(nil, x)` 会自动新建切片。因此无需先判断 key 是否存在,一句话完成"取出旧组 → 追加 → 存回"。
4. **map 遍历顺序是故意随机的**:Go 每次遍历 map 顺序都不同。本题"组间顺序任意"恰好无感;若题目要求有序输出,必须先收集到切片再显式 sort。

## 拓展:计数法(第二把刀,选学)

除了"排序归一化",还可以**数字母**:`"eat"` → (a:1, e:1, t:1),计数相同即异位词,省掉排序,时间 O(n·k)。

```go
m := make(map[[26]int][]string) // 注意:数组可以做 key,切片不行!
for _, str := range strs {
    var cnt [26]int
    for _, ch := range str {
        cnt[ch-'a']++
    }
    m[cnt] = append(m[cnt], str)
}
```

- **Go 的坑**:切片不能做 map key(不可比较),**数组可以**(可比较)。`map[[26]int][]string` 完全合法 —— "数组 vs 切片"的经典考点。

- k=100 时 log k ≈ 7,两种方法差距很小都能过;面试时说出两种并比较优劣,是加分项。

## 复盘总结

<!-- 每天向老师请教的问题、犯的错误、关键收获,都由老师补充在这里 -->

- **学生自主解出**(2026-09-03):独立设计"排序作 key + map\[string]\[]string 分组",一次通过,仅 1 处 gofmt 空格问题

- **Day 1 → Day 2 的迁移**:核心思想一脉相承 —— "设计特征做 key,让 map 把查找/归类变成 O(1)";升级点:① key 从数字变成"变形后的字符串"(排序归一化),② value 从单个下标升级为 \[]string(一组),③ 学到分组惯用法 `m[key] = append(m[key], str)`

- **新 Go 语法**:string 不可变 → \[]byte 排序 → string 回转;sort.Slice 比较器(首次接触"函数当参数");append(nil, x) 自动分配切片;map 遍历顺序随机

- **复杂度第二课**:总时间 = 循环次数 × 每圈的活儿;排序自带 k log k → 本题 O(n·k log k);暴力两两比较是 O(n²·k) ≈ 10¹⁰,超出"1 秒法则"

- **拓展认知**:计数法 `map[[26]int][]string`;数组可做 key 而切片不行(可比较性考点)

- **自测已布置**:`strs=[""]` 边界走查,待学生推演回答后补充

