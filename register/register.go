package register

import(
	"strings"
	"time"
	"log"
	"encoding/json"

	"../common"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/deckarep/golang-set"
	"github.com/fsouza/go-dockerclient"

)

type Worker struct {
	ip			string
	servSet		mapset.Set
	keysAPI		client.KeysAPI
	dockerClt	*docker.Client
}

func NewWorker(IP string, servTypes []string, etcdEpts []string, dockerEpt string) *Worker {
	etcdCfg := client.Config {
		Endpoints:					etcdEpts,
		Transport:					client.DefaultTransport,
		HeaderTimeoutPerRequest:	time.Second,
	}

	etcdClient, err := client.New(etcdCfg)
	if err != nil {
		log.Println("Failed to connect to etcd: ", err)
		return nil
	}

	dockerClt, err := docker.NewClient(dockerEpt)
	if err != nil {
		log.Println("Failed to connect to docker daemon: ", err)
		return nil
	}

	w := &Worker {
		ip:			IP,
		servSet:	mapset.NewThreadUnsafeSet(),
		keysAPI:	client.NewKeysAPI(etcdClient),
		dockerClt:	dockerClt,
	}
	for _, s := range servTypes {
		w.servSet.Add(s)
	}

	return w
}

func (w *Worker) HeartBeat() {
	conts, err := w.dockerClt.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Print("Error in getting containers: ", err)
		return
	}

	for _, c := range conts {
		if len(c.Names) == 0 {
			continue
		}
		servName := c.Names[0][1:]
		if !w.servSet.Contains(servName) {
			continue
		}
		info := common.NodeInfo {
			 IP:		w.ip,
			 Service:	servName,
		}
		if strings.Contains(c.Status, "Exit") {
			info.Status = "DOWN"
		}else {
			info.Status = "UP"
		}

		value, _ := json.Marshal(info)
		_, err := w.keysAPI.Set(context.Background(),
				"/service/"+servName+w.ip,
				string(value),
				&client.SetOptions{TTL: 2*time.Second},
		)
		if err != nil {
			log.Print("Failed to report node info: ", err)
			break
		}
		if strings.Contains(c.Status, "Exit") {
			w.dockerClt.RestartContainer(c.ID, 1)
		}
	}
}