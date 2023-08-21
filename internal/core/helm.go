package core

import (
	"fmt"
	"log"
	"os"

	"github.com/dumunari/spikectl/internal/config"
	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

type Helm struct {
	coreComponents []config.CoreComponent
	settings       *cli.EnvSettings
	actionConfig   *action.Configuration
	kubeConfig     config.KubeConfig
}

func InstallCoreComponents(spikeConfig *config.Spike, kubeConfig config.KubeConfig) {
	fmt.Println("[🐶] Checking Core Components...")

	helm := Helm{}
	helm.coreComponents = spikeConfig.Spike.CoreComponents
	helm.kubeConfig = kubeConfig

	helm.initConfiguration("", helm.kubeConfig)

	installedReleases := helm.checkInstalledComponents()

	for _, coreComponent := range helm.coreComponents {
		if !slices.Contains(installedReleases, coreComponent.ReleaseName) {
			helm.installCoreComponent(coreComponent)
		}
	}

	fmt.Println("[🐶] All Core Components successfully installed.")
}

func (h *Helm) initConfiguration(namespace string, kubeConfig config.KubeConfig) {
	h.settings = cli.New()
	h.settings.KubeAPIServer = kubeConfig.EndPoint
	h.settings.KubeCaFile = kubeConfig.CaFile
	h.settings.KubeToken = kubeConfig.Token
	h.actionConfig = new(action.Configuration)
	if err := h.actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}

func (h *Helm) checkInstalledComponents() []string {
	var installedReleases []string

	client := action.NewList(h.actionConfig)
	client.Short = true
	client.AllNamespaces = true
	client.Deployed = true
	client.Filter = "spike-"

	releases, _ := client.Run()

	for _, release := range releases {
		fmt.Printf("[🐶] Found release %s at namespace %s\n", release.Name, release.Namespace)
		installedReleases = append(installedReleases, release.Name)
	}

	return installedReleases
}

func (h *Helm) installCoreComponent(coreComponent config.CoreComponent) {
	if err := h.actionConfig.Init(h.settings.RESTClientGetter(), coreComponent.Namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewInstall(h.actionConfig)
	client.Namespace = coreComponent.Namespace
	client.ReleaseName = coreComponent.ReleaseName
	client.RepoURL = coreComponent.Repository
	client.Version = coreComponent.ChartVersion
	client.CreateNamespace = true
	client.IsUpgrade = true

	fetch_chart, err := client.LocateChart(coreComponent.Chart, h.settings)
	if err != nil {
		panic(err)
	}

	chart, err := loader.Load(fetch_chart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[🐶] %s is being installed...\n", coreComponent.ReleaseName)
	_, err = client.Run(chart, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[🐶] %s(using chart version %s) successfully installed on namespace %s\n", coreComponent.ReleaseName, coreComponent.ChartVersion, coreComponent.Namespace)
}
