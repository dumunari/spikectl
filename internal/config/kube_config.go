package config

// TODO: not sure if this should liv here
const KUBE_CA_PATH string = "/tmp/spikectl/"
const KUBE_CA_FILE_PATH string = KUBE_CA_PATH + "kubecafile"

type KubeConfig struct {
	EndPoint string `json:"endpoint"`
	CaFile   string `json:"ca_file"`
	Token    string `json:"token"`
}
