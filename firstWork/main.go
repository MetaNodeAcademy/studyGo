package main

import "fmt"

//题目1
//给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。
//找出那个只出现了一次的元素。可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
//例如通过 map 记录每个元素出现的次数，
//然后再遍历 map 找到出现次数为1的元素。

// 题目2 是否是回文数
// 数字判断
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	tmp, renum := x, 0
	//将数字倒转
	for tmp > 0 {
		renum = renum*10 + tmp%10
		tmp /= 10
	}
	return x == renum
}

// 字符串判断
func isPalindromeString(x string) bool {
	if x == "" {
		return false
	}
	tempStr := []rune(x)
	for i := 0; i < len(tempStr)/2; i++ {
		if tempStr[i] != tempStr[len(tempStr)-i-1] {
			return false
		}
	}
	return true
}

//题目3
//判断是否包含有效括号

func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	stack := []rune{}
	for _, s1 := range s {
		if s1 == '(' || s1 == '[' || s1 == '{' {
			stack = append(stack, s1)
		} else {
			if len(stack) == 0 {
				return false
			}

			if s1 == ')' && stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
			} else if s1 == ']' && stack[len(stack)-1] == '[' {
				stack = stack[:len(stack)-1]
			} else if s1 == '}' && stack[len(stack)-1] == '{' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		}

	}
	//核心思路，使用堆栈模式，遇到相反即出栈，栈内为空即为全部出栈，则为正确
	return len(stack) == 0
}

// 题目4
// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	answer := ""
	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {
			if i == len(strs[j]) || strs[j][i] != strs[0][i] {
				return answer
			}
		}
		answer += string(strs[0][i])
	}
	return answer
	//嵌套循环比较，不一致时退出
}

// 题目5
func plusOne(digits []int) []int {
	//先把数组转成成数字再+1
	num := 0
	newDigits := []int{}
	for i := 0; i < len(digits); i++ {
		num = num*10 + digits[i]
	}
	num++
	for num > 0 {
		newDigits = append([]int{num % 10}, newDigits...)
		num /= 10
	}
	return newDigits
}

// 题目6
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] == nums[i+1] {
			nums = append(nums[:i], nums[i+1:]...)
			i--
		}
	}
	return len(nums)
}

// 题目7
func mergeInterval(nums [][]int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}
	for i := 0; i < len(nums)-1; i++ {
		if nums[i][1] > nums[i+1][0] {
			nums[i][1] = nums[i+1][1]
			nums = append(nums[:i+1], nums[i+2:]...)
			i--
		}
	}
	return nums
}

// 题目8
func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}
func main() {
	//题目1解法
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 5}
	fmt.Println(nums)
	numsMap := make(map[int]int)
	for i, num := range nums {
		if numsMap[num] == 1 {
			fmt.Println("找到了重复数值为:", num)
		} else {
			fmt.Println("暂时未找到重复数")
			numsMap[num] = 1
		}
		fmt.Println(i, num)
	}

	//题目2，判断一个整数是否是回文数
	flag := isPalindrome(12321)
	fmt.Println(flag)
	fmt.Println(isPalindromeString("12321"))

	// 题目3
	fmt.Println(isValid("{{{[][]([])}}}"))

	//题目4
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))

	//题目5
	fmt.Println(plusOne([]int{1, 2, 9}))
	//题目6
	fmt.Println(removeDuplicates([]int{1, 1, 2, 3, 4, 5, 5, 6, 7, 8, 9}))
	//题目7
	fmt.Println(mergeInterval([][]int{{1, 3}, {2, 6}, {8, 10}, {9, 18}}))
	//题目8
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
}
