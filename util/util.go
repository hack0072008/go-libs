package util

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/cihub/seelog"
)


/*
生成 mac 地址 改成大写形式
 */
func GenMacAddrs(prefix []byte, count int, exclude []string) []string {
	macs := make(map[string]bool, count)
	excludeMap := make(map[string]bool, len(exclude))
	for _, mac := range exclude {
		excludeMap[mac] = true
	}

	fields := prefix[:]
	prefixLen := len(fields)

	for len(macs) < count {
		fields = fields[:prefixLen]
		for j, r := 0, rand.Uint64(); len(fields) < 6; j++ {
			fields = append(fields, byte(r>>uint64(j*8)))
		}
		mac := net.HardwareAddr(fields).String()
		if _, exist := excludeMap[mac]; exist {
			continue
		}
		macs[mac] = true
	}

	result := make([]string, 0, len(macs))
	for mac := range macs {
		result = append(result, strings.ToUpper(mac))
	}

	return result
}


/**
*生成随机字符
**/
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, fmt.Sprintf("%d",rand.Intn(26)+65))
		} else {
			rs = append(rs, fmt.Sprintf("%d",rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}


// 参数非法校验
func ParamIllegalCheck(param ...interface{}) error {
	// 参数为空校验
	for _, v := range param {
		switch v.(type) {
		case string:
			if v == "" {
				return errors.New("Param is Empty or Miss")
			}
		case int:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		case int16:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		case int32:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		case int64:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		case float64:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		case float32:
			if v == 0 {
				return errors.New("Param is Empty or Miss")
			}
		default:
			if v == nil {
				return errors.New("Param is Empty or Miss")
			}
			if p, ok := v.([]interface{}); ok {
				seelog.Debug("[]interface{}")
				if len(p) == 0 {
					return errors.New("Param is Empty or Miss")
				}
			}
			if p, ok := v.(map[string]interface{}); ok {
				if len(p) == 0 {
					return errors.New("Param is Empty or Miss")
				}
			}
		}
	}
	return nil
}

// to int
func ParamIntChange(param interface{}) (int, error) {
	switch param.(type) {
	case string:
		result, err := strconv.Atoi(param.(string))
		if err != nil {
			return 0, fmt.Errorf("Transfer string to int  Failed: %s", err)
		}
		return result, nil
	case float64:
		result := int(param.(float64))
		return result, nil
	case float32:
		result := int(param.(float32))
		return result, nil
	case int:
		return param.(int), nil
	default:
		return 0, errors.New("Int Not Fount Param Type. ")
	}
	return 0, errors.New("PASS.")
}

// 参数是否为 string 类型
func ParamIllegalStringCheck(param ...interface{}) error {
	for _, v := range param {
		if _, ok := v.(string); !ok {
			return errors.New("Param is not Type of string")
		}
	}
	return nil
}

// 参数是否为 float64 类型
func ParamIllegalFloat64Check(param ...interface{}) error {
	for _, v := range param {
		if _, ok := v.(float64); !ok {
			return errors.New("Param is not Type of float64")
		}
	}
	return nil
}

// to string
func ParamStringChange(param interface{}) (string, error) {
	switch param.(type) {
	case int:
		result := strconv.Itoa(param.(int))
		return result, nil
	case float64:
		result := strconv.FormatFloat(param.(float64), 'f', -1, 64)
		return result, nil
	case uint64:
		result := strconv.FormatUint(param.(uint64), 10)
		return result, nil
	// case float32:
	// 	return "", errors.New("Param Type float32")
	case string:
		return param.(string), nil
	default:
		if p, ok := param.([]interface{}); ok {
			var result string
			for _, value := range p {
				if _, ok := value.(string); ok {
					result = result + value.(string) + ","
				}
			}
			return trimSuffix(result, ","), nil
		}
		if p, ok := param.([]string); ok {
			var result string
			for _, value := range p {
				result = result + value + ","
			}
			return trimSuffix(result, ","), nil
		}
		return "", errors.New("String Not Fount Param Type. ")
	}
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// 随机三位数
func RandNum3() string {
	rand.Seed(time.Now().UnixNano())
	n := fmt.Sprintf("%v%v%v", strconv.Itoa(rand.Intn(9)), strconv.Itoa(rand.Intn(9)), strconv.Itoa(rand.Intn(9)))
	return n
}

// 随机四位数
func RandNum4() string {
	rand.Seed(time.Now().UnixNano())
	n := fmt.Sprintf("%v%v%v%v", strconv.Itoa(rand.Intn(9)), strconv.Itoa(rand.Intn(9)), strconv.Itoa(rand.Intn(9)), strconv.Itoa(rand.Intn(9)))
	return n
}

// 参数是否为 []interface{} 类型
func ParamIllegalListInterfaceCheck(param ...interface{}) error {
	for _, v := range param {
		if _, ok := v.([]interface{}); !ok {
			return errors.New("Param is not Type of list interface")
		}
	}
	return nil
}

// 参数是否为 []string{} 类型
func ParamIllegalListStringCheck(param ...interface{}) error {
	for _, v := range param {
		if _, ok := v.([]string); !ok {
			return errors.New("Param is not Type of list string ")
		}
	}
	return nil
}

// 参数是否为 map[string]interface{} 类型
func ParamIllegalMapCheck(param ...interface{}) error {
	for _, v := range param {
		if _, ok := v.(map[string]interface{}); !ok {
			return errors.New("Param is not Type of map")
		}
	}
	return nil
}


// id转化为ip, 例如:10064131001 -> 10.64.131.1
func IdToIp(id uint64) string {
	idS := fmt.Sprintf("%d", id)
	if len(idS) == 11 {
		idS = fmt.Sprintf("0%s", idS)
	} else if len(idS) == 10 {
		idS = fmt.Sprintf("00%s", idS)
	}
	idStr := []byte(idS)

	var arr []int
	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(string(idStr[i*3 : (i+1)*3]))
		if err != nil {
			panic(fmt.Errorf("%v id", err))
		}
		arr = append(arr, n)
	}

	var str = ""
	for i := 0; i < 4; i++ {
		str += fmt.Sprintf("%d", arr[i])
		if i != 3 {
			str += "."
		}
	}
	return str
}

