package util



import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

/*
UCloudSign 根据一定规则对API请求体进行签名
将结构体的字段，按照 key 升序排序，组成字符串，拼接 privateKey 之后,做 SHA1 求摘要，作为签名
 */
func UCloudSign(privateKey string, params map[string]interface{}) string {
	var builder strings.Builder
	keys := make(sort.StringSlice, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(fmt.Sprintf("%v", params[key]))
	}
	builder.WriteString(privateKey)
	r := sha1.Sum([]byte(builder.String()))
	return hex.EncodeToString(r[:])
}
