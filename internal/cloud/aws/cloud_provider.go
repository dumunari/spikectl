package aws

type AwsCloudProvider struct {
}

func (a AwsCloudProvider) InstantiateKubernetesCluster() error {
	//TODO implement me
	panic("implement me")
	return nil
}

func NewAwsCloudProvider() *AwsCloudProvider {
	return &AwsCloudProvider{}
}
