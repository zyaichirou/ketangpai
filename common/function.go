package common

import (
	"math/rand"
	"time"
)

//RandNumbers
//@title		RandNumbers()
//@description	生成随机的6位数
//@author		zy
//@param
//@return		string
func RandNumbers() string {
	var nums = []byte("0123456789")
	result := make([]byte, 6)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = nums[rand.Intn(len(nums))]
	}

	return string(result)
}


//RandClassId
//@title		RandClassId()
//@description	生成随机的6位数
//@author		zy
//@param
//@return		string
func RandClassId() string{
	var nums = []byte("QWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	result := make([]byte, 6)

	rand.Seed(time.Now().Unix())
	for i := range result{
		result[i] = nums[rand.Intn(len(nums))]
	}

	return string(result)
}
