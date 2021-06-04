package zookeeper

import (
	"encoding/json"
	"fmt"
)

type ZkConfig struct {
	Region     interface{}       `json:"region"`
	Names      map[string]string `json:"names"`
	Connstring string            `json:"connstring"`
}

type MongoConfig struct {
	Url    string `json:"uri"`
	RsName string `json:"rs_name"`
}

type ConfigMessage struct {
	RegionId   int        `json:"region_id"`
	ZkNodePath string     `json:"zk_node_path"`
	ZkConfig   []ZkConfig `json:"zk_config"`
}

type NodeConfig struct {
	MongoDB       map[string]string `json:"mongo_db"`
	NetListenAddr string            `json:"net_listen_addr"`
	NetListenPort uint32            `json:"net_listen_port"`
	RegionId      uint32            `json:"region_id"`

	ZkConfig   []ZkConfig `json:"zk_config"`
	ZkNodePath string     `json:"zk_node_path"`
}

//zk config info :   Configuration information that should be obtained from zk configuration at startup
var cfg = NodeConfig{
	MongoDB: map[string]string{
		"rs_name": "myrs",
		"uri":     "mongodb://10.0.0.1:27017/mydb,mongodb://10.0.0.2:27017/mydb,mongodb://10.0.0.3:27017/mydb",
	},
	NetListenAddr: "@eth0",
	NetListenPort: 8080,
	RegionId:      10001,
	ZkConfig: []ZkConfig{
		{
			/*Connstring: "10.0.0.1:2181,10.0.0.2:2181,10.0.0.3:2181",*/
			Connstring: "127.0.0.1:2181",
			Names: map[string]string{
				"access":  "/NS/region10009/access",
				"manager": "/NS/region10009/manager",
			},
			Region: 10001,
		},
	},
	ZkNodePath: "/NS/region10009/access",
}

//
var config = ConfigMessage{
	RegionId:   10009,
	ZkNodePath: "/cm/_A_uimage3-access/_S_set99/access",
	ZkConfig: []ZkConfig{
		{
			Region: 10009,
			Names: map[string]string{
				"access":  "/NS/region10009/access",
				"manager": "/NS/region10009/manager",
			},
			Connstring: "127.0.0.1:2181",
		},
	},
}

func main() {

	//============ test for zk instance

	//get zk conn instance
	conn, err := GetZkInstance(config.ZkConfig[0].Connstring)
	if err != nil {
		fmt.Println("Get zk instance failed ....")
		return
	}

	//create a new node to store the configuration information needed at startup
	resCfg, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("marshall zk config info faild ...")
		return
	}
	resPath, err := conn.CreateNewNode(config.ZkNodePath, resCfg)
	if err != nil {
		fmt.Println("Create New Node failed ....")
		return
	}
	fmt.Printf("Create New node success : |->{%s}<-| \n", resPath)

	data, err := conn.GetNode(config.ZkNodePath)
	if err != nil {
		fmt.Println("failed to get node info : ", err)
		return
	}
	fmt.Println("zk node info : ", string(data))

	//polling for zk configuration information
}
