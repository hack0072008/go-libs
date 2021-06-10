package util

import (
	"github.com/hack0072008/go-libs/log"
	"net"
	"os"
	"strconv"
	"time"
)

/*
 心跳发送
*/
func HeartBeatSender(conn *net.TCPConn) {
	for i := 0; i < 10; i++ {
		words := strconv.Itoa(i) + " Hello I'm MyHeartbeat Client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			log.Errorf("connecct to %s %s %s", conn.RemoteAddr().String(), "Fatal error: ", err.Error())
			os.Exit(1)
		}
		log.Infof("服务端接收了:%d", msg)
		time.Sleep(2 * time.Second)
	}
	for i := 0; i < 2; i++ {
		time.Sleep(12 * time.Second)
	}
	for i := 0; i < 10; i++ {
		words := strconv.Itoa(i) + " Hi I'm MyHeartbeat Client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			log.Errorf("connect to %s %s %s", conn.RemoteAddr().String(), "Fatal error: ", err.Error())
			os.Exit(1)
		}
		log.Infof("服务端接收了:%d", msg)
		time.Sleep(2 * time.Second)
	}

}

/*
 心跳保活：每次接收到心跳数据就 SetDeadline 延长一个时间段 timeout。如果没有接到心跳数据，5秒后连接关闭。
*/
func HeartBeatLifeAdd(conn net.Conn, bytes chan byte, timeout int) {
	select {
	case fk := <-bytes:
		log.Infof("client %s %s %s %s", conn.RemoteAddr().String(), "心跳:第", string(fk), "times")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break

	case <-time.After(5 * time.Second):
		log.Warn("conn dead now")
		conn.Close()
	}
}
