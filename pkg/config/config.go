package config

type Config struct {
	KubeConfig   string
	Create       bool
	Delete       bool
	Count        int
	StorageClass string
	PVCSize      string
	Namespace    string
	ImageName    string
	ImageTag     string
}
