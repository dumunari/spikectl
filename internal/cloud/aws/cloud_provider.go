package aws

type CloudProvider struct {
}

func (a CloudProvider) InstantiateKubernetesCluster() error {
	//TODO implement me
	panic("implement me")
	return nil
}

func NewAwsCloudProvider() *CloudProvider {
	return &CloudProvider{}
}
