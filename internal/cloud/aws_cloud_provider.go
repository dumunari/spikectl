package cloud

type AwsCloudProvider struct {
}

func (p *AwsCloudProvider) CreateResourceGroup() {
	panic("implement me")
}

func (p *AwsCloudProvider) CreateKubernetesCluster() {}

func NewAwsCloudProvider() *AwsCloudProvider {
	return &AwsCloudProvider{}
}
