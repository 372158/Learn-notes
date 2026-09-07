# 1. 两数之和 (Two Sum)

> 难度：简单 | 标签：数组、哈希表
> 来源：[力扣 HOT 100 - Day 1](https://leetcode.cn/problems/two-sum/description/?envType=study-plan-v2\&envId=top-100-liked)
> 状态：✅ 已完成(学生独立解出,2026-09-01)

## 题目描述

给定一个整数数组 `nums` 和一个整数目标值 `target`,请你在该数组中找出 **和为目标值** `target` 的那 **两个** 整数,并返回它们的数组下标。

你可以假设每种输入只会对应一个答案,并且你**不能使用两次相同的元素**。

你可以按任意顺序返回答案。

## 示例

**示例 1:**

```text
输入:nums = [2,7,11,15], target = 9
输出:[0,1]
解释:因为 nums[0] + nums[1] == 9 ,返回 [0, 1] 。
```

**示例 2:**

```text
输入:nums = [3,2,4], target = 6
输出:[1,2]
```

**示例 3:**

```text
输入:nums = [3,3], target = 6
输出:[0,1]
```

## 提示

- `2 <= nums.length <= 10^4`

- `-10^9 <= nums[i] <= 10^9`

- `-10^9 <= target <= 10^9`

- 只会存在一个有效答案

**进阶:** 你可以想出一个时间复杂度小于 `O(n^2)` 的算法吗?

***

## 第一课:复杂度入门(2026-09-01)

复杂度只回答一个问题:**数据变大时,步数涨得多快**。三级台阶:

| 记号    | 生活画面          | 程序长相      | 数据×100 后  |
| ----- | ------------- | --------- | --------- |
| O(1)  | 翻字典到指定页码,永远一下 | 没有循环,固定几步 | 步数不变      |
| O(n)  | 数全场人数,一个一个数   | 一层 for 循环 | 慢 100 倍   |
| O(n²) | 全场每两人握一次手     | 两层嵌套 for  | 慢 10000 倍 |

- 数步数:暴力两层循环是 n+(n-1)+…+1 ≈ n²/2 步;丢掉常数和低阶项、只留最高阶 → 记 O(n²)。复杂度看的是"涨多快",不是具体数字。

- **1 秒法则**:计算机一秒约做 10⁸ 次简单操作。n=10⁴ 时:O(n)≈10⁴ 步(瞬间完成),O(n²)≈10⁸ 步(贴极限)。题面"进阶:能否 < O(n²)"翻译过来就是——想办法砍掉一层循环。

## 思路:哈希表

**什么时候想到哈希法?** 当需要查询一个元素是否出现过(是否在集合里)时,第一时间想到哈希法。本题需要一个集合来存放遍历过的元素,遍历时去询问集合"某个元素是否出现过"。

**为什么用 map 而不是数组或 set?**

- 本题不仅要知道元素有没有遍历过,还要知道元素对应的**下标**,需要 key-value 结构:key 存元素,value 存下标。

- 数组做哈希:大小受限,元素少而哈希值大会浪费内存。

- set 做哈希:只能存一个 key,无法记录下标。

以 C++ 为例,三种 map 的对比:

| 映射                   | 底层实现 | 是否有序   | 数值是否可以重复 | 查询效率     | 增删效率     |
| -------------------- | ---- | ------ | -------- | -------- | -------- |
| `std::map`           | 红黑树  | key 有序 | key 不可重复 | O(log n) | O(log n) |
| `std::multimap`      | 红黑树  | key 有序 | key 可重复  | O(log n) | O(log n) |
| `std::unordered_map` | 哈希表  | key 无序 | key 不可重复 | O(1)     | O(1)     |

本题不需要 key 有序,选 `unordered_map`(哈希表)效率最高。

**遍历过程:**

1. 遍历数组,对当前元素 `nums[i]`,去 map 中查询是否有匹配的 `target - nums[i]`;
2. 有 → 找到答案,返回两个下标;
3. 没有 → 把当前元素和下标存入 map,继续遍历。

## 代码实现

**C++:**

```cpp
class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
        std::unordered_map <int,int> map;
        for(int i = 0; i < nums.size(); i++) {
            // 遍历当前元素,并在map中寻找是否有匹配的key
            auto iter = map.find(target - nums[i]);
            if(iter != map.end()) {
                return {iter->second, i};
            }
            // 如果没找到匹配对,就把访问过的元素和下标加入到map中
            map.insert(pair<int, int>(nums[i], i));
        }
        return {};
    }
};
```

**Java:**

```java
public int[] twoSum(int[] nums, int target) {
    int[] res = new int[2];
    if(nums == null || nums.length == 0){
        return res;
    }
    Map<Integer, Integer> map = new HashMap<>();
    for(int i = 0; i < nums.length; i++){
        int temp = target - nums[i];
        if(map.containsKey(temp)){
            res[1] = i;
            res[0] = map.get(temp);
        }
        map.put(nums[i], i);
    }
    return res;
}
```

**Go:**

```go
func twoSum(nums []int, target int) []int {
    hashTable := map[int]int{}
    for i, x := range nums {
        if p, ok := hashTable[target-x]; ok {
            return []int{p, i}
        }
        hashTable[x] = i
    }
    return nil
}
```

**Python:**

```python
class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        hashtable = {}
        for i, num in enumerate(nums):
            if target - num in hashtable:
                return [hashtable[target - num], i]
            hashtable[num] = i
        return []
```

## 复杂度分析

- **时间复杂度:** O(n),其中 n 是数组长度,只需遍历一次数组,哈希表查询为 O(1)。

- **空间复杂度:** O(n),哈希表最多存储 n 个元素。

## 总结

本题的四个重点:

1. 为什么会想到用哈希表 —— 需要快速判断"某元素是否出现过";
2. 哈希表为什么用 map —— 既要判断存在性,又要记录下标;
3. 本题 map 是用来存什么的 —— 存放已经访问过的元素;
4. map 的 key 和 value 分别存什么 —— key 存数组元素,value 存元素下标。

把这四点想清楚了,本题才算理解透彻。

## 我的代码(第一版,Go)

```go
func twoSum(nums []int, target int) []int {
    m := map[int]int{}              // 字典:key=数组元素, value=下标
    for i, num := range nums {
        need := target - num        // 我缺的另一半
        if j, ok := m[need]; ok {   // 先查:另一半之前出现过吗?
            return []int{j, i}
        }
        m[num] = i                  // 后存:把自己记进字典
    }
    return nil
}
```

**Go ↔ C++ 对照**(同一逻辑的两种说法):

| Go                   | C++                                     | 含义                  |
| -------------------- | --------------------------------------- | ------------------- |
| `m := map[int]int{}` | `unordered_map<int,int> m;`             | 字典:key=元素, value=下标 |
| `j, ok := m[need]`   | `auto it = m.find(need); it != m.end()` | 查:key 存在吗           |
| `j`                  | `it->second`                            | 查到的下标               |
| `m[num] = i`         | `m[nums[i]] = i;`                       | 存:记下自己和位置           |
| `return nil`         | `return {};`                            | 没找到时的兜底             |

> 纯 C 没有内置哈希表,得手写哈希函数、开桶数组、挂链表——所以刷题一般用带标准库的语言,把精力留给算法本身。

## 复盘总结

<!-- 每天向老师请教的问题、犯的错误、关键收获,都由老师补充在这里 -->

- **学生基础**(2026-09-01):数据结构零基础,第一次接触复杂度分析 → 第一课先补 Big-O 直觉

- **复杂度速记**:O(1) 固定几步;O(n) 一层循环扫一遍;O(n²) 两层循环、每对元素配一次。只看最高阶、忽略常数:数据大 100 倍,O(n) 慢 100 倍,O(n²) 慢 10000 倍

- **本题暴力法**:外层 i × 内层 j,约 n²/2 次比较。n=10⁴ 时约 5×10⁷ 次,贴着"1 秒 ≈ 10⁸ 次操作"的经验极限 → 题目"进阶"就是在提示找 O(n²) 以下的算法

- **学生自主推出哈希表**(2026-09-01):准确说出"查 = 知道位置直接定位;找 = 不清楚在哪,只能一个个看",并用字典类比提出"数 → 下标"的映射 —— 这就是哈希表:key 存数组元素,value 存下标

- **学生验证**(2026-09-01):手动模拟 `nums=[3,2,4]、target=6`,先查后存 → 返回 `[1,2]` ✓;并独立发现"先存后查会让 `6-3=3` 查到自己"的原因 ✓

- **第一版代码批改**(Go):逻辑全对。修正 2 处编译问题:① 循环前需声明字典 `m := map[int]int{}`;② Go 要求函数有兜底 return,末尾补 `return nil`。风格 2 处:gofmt 会把 `[]int{j,i}` 写作 `[]int{j, i}`;变量名 `com` 改成 `need`(我需要的数)更达意

- **Go 惯用法积累**:`for i, num := range nums` 同时给下标和值,天生适配本题;`v, ok := m[k]`(comma-ok)是 Go 判断"key 是否存在"的标准写法

- **C++ 对照**:`unordered_map` + `find()/end()` 迭代器,`it->second` 对应 Go 的 `j`。纯 C 没有内置哈希表、需要手写,所以刷题一般用带标准库的语言

- **巩固自测**:`nums=[3,3]、target=6` 用"先查后存"走一遍 → i=0 时查不到、存 3→0;i=1 时 need=3 查到 j=0,返回 \[0,1] ✓ —— 先查后存天然兼容重复元素(只要答案合法)

