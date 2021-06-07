package util

import (
	"math/rand"
	"net"
	"strings"
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


