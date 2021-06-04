package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/hack0072008/go-libs/log"
	"github.com/samuel/go-zookeeper/zk"
)

type ZkConn struct {
	Addr string
	conn *zk.Conn
}

var (
	ZkConnPoolMu  sync.Mutex
	ZkConnPool    = make(map[string]*ZkConn, 1)
	ErrInvalidSev = errors.New("invalid conncet addr info")
	ErrorConn     = errors.New("can't connect to remote servers")
)

//deal error
func checkError(err error, errInfo string) bool {
	if err != nil {
		log.Infof("%s : |->{%s}-<|", errInfo, err)
		return false
	}
	return true
}

//Initialize zk link
func NewZkConn(addr string) (*ZkConn, error) {

	//illegal judgement
	if addr == "" {
		err := ErrInvalidSev
		return nil, err
	}
	//eg:127.0.0.0:1234,127.0.0.1:1235
	connAddr := strings.Split(addr, ",")
	rand.Seed(time.Now().UnixNano())
	//add reconnection measure
	retryTick := time.NewTicker(time.Second)
	var (
		ec    <-chan zk.Event
		conn  *zk.Conn
		count int
		err   error
	)
	for _ = range retryTick.C {
		conn, ec, err = zk.Connect(connAddr, 10*time.Second, zk.WithLogger(log.NewUidLog(context.Background())))
		if err != nil {
			if count > 2 {
				log.Fatal("connect to zk failed : ", err)
				return nil, err
			}
			count++
			continue
		}
		retryTick.Stop()
		break
	}
	for {
		select {
		case connEvent, ok := <-ec:
			if ok {
				switch connEvent.State {
				case zk.StateHasSession:
					zkConn := &ZkConn{
						Addr: addr,
						conn: conn,
					}
					return zkConn, nil
				default:
					continue
				}
			} else {
				err = ErrorConn
				return nil, err
			}
		default:
			continue
		}
	}
}

//Get the instance information of the zk cluster connection
func GetZkInstance(zkCluster string) (*ZkConn, error) {
	ZkConnPoolMu.Lock()
	defer ZkConnPoolMu.Unlock()

	zkConn, ok := ZkConnPool[zkCluster]
	if ok {
		//old session exception ， rebuild and add into connection map
		if zkConn.conn.State() != zk.StateHasSession {
			zkConn.conn.Close()
			zkConn, err := NewZkConn(zkCluster)
			if !checkError(err, "Rebuild ZkConn") {
				return nil, err
			}
			ZkConnPool[zkCluster] = zkConn
		}
		return zkConn, nil
	}

	zkConn, err := NewZkConn(zkCluster)
	if !checkError(err, "new zk conn session") {
		return nil, err
	}
	ZkConnPool[zkCluster] = zkConn

	return zkConn, nil
}

//create new link node
func (c *ZkConn) CreateNewNode(path string, data []byte) (string,
	error) {
	if path == "" {
		return "", errors.New("Invalid path")
	}

	flag := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)
	childPath := path

	//create father node
	paths := strings.Split(path, "/")
	var parentPath string
	for _, v := range paths[1 : len(paths)-1] {
		parentPath += "/" + v
		exist, _, err := c.NodeExists(parentPath)
		if !checkError(err, "judgeMent Father Node if exist") {
			return "", err
		}
		if !exist {
			_, err = c.conn.Create(parentPath, nil, 0, acl)
			if !checkError(err, "Create Father Node") {
				return "", err
			}
		}
	}

	//create child node
	exist, _, err := c.NodeExists(childPath)
	if !checkError(err, "JudgeMent Father Node if exist") {
		return "", err
	}
	resPath := ""
	if !exist {
		resPath, err = c.conn.Create(childPath, data, flag, acl)
		if !checkError(err, "Create Child Node") {
			return "", err
		}
	} else {
		err = fmt.Errorf("[%s]  exists", childPath)
	}
	return resPath, nil
}

//set node info
func (c *ZkConn) SetNode(path string, data []byte) error {
	if path == "" {
		return errors.New("Invalid Path")
	}

	//judgement node if exist
	exist, stat, err := c.NodeExists(path)
	if !checkError(err,
		"func SetNode ==> check node "+path+"exist failed") {
		return err
	}
	//not exist
	if !exist {
		return fmt.Errorf("node [%s] dosen't exist,can't be setted", path)
	}

	_, err = c.conn.Set(path, data, stat.Version)
	if !checkError(err, "Set node Info") {
		return err
	}

	return err
}

func (c *ZkConn) GetNode(path string) ([]byte, error) {
	//whether the path is null
	if path == "" {
		return nil, errors.New("Invalid Path")
	}
	log.Infof("================>> %s <<================\n", path)
	data, _, err := c.conn.Get(path)

	return data, err
}

func (c *ZkConn) DeleteNode(path string) error {
	//whether the path is null
	if path == "" {
		return errors.New("Invalid Path")
	}
	//exist or not
	exist, stat, err := c.NodeExists(path)
	if !checkError(err, "Before Delete Node") {
		return err
	}
	if !exist {
		return fmt.Errorf("path [\"%s\"] doesn't exist", path)
	}

	//delete node
	return c.conn.Delete(path, stat.Version)
}

//list child path info
func (c *ZkConn) ListChildren(path string) ([]string, error) {

	children, _, err := c.conn.Children(path)
	return children, err
}

//get the watcher of current node (the full path info)
func (c *ZkConn) GetNodeWatcher(path string) ([]byte, *zk.Stat,
	<-chan zk.Event, error) {
	return c.conn.GetW(path)
}

//get the change watcher of all child nodes of the current node
func (c *ZkConn) GetChildrenWatcher(path string) ([]string, *zk.Stat,
	<-chan zk.Event, error) {
	return c.conn.ChildrenW(path)
}

//if node exist
func (c *ZkConn) NodeExists(path string) (bool, *zk.Stat,
	error) {
	//if node exist
	return c.conn.Exists(path)
}
