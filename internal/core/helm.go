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
	coreConfig   []config.CoreComponent
	settings     *cli.EnvSettings
	actionConfig *action.Configuration
}

func InstallCoreComponents(config *config.Spike) {
	fmt.Println("[🐶] Checking Core Components...")

	helm := Helm{}
	helm.coreConfig = config.Spike.CoreConfig

	helm.initConfiguration("")

	installedReleases := helm.checkInstalledComponents()

	for _, coreComponent := range helm.coreConfig {
		if !slices.Contains(installedReleases, coreComponent.ReleaseName) {
			helm.installCoreComponents(coreComponent)
		}
	}
}

func (h *Helm) initConfiguration(namespace string) {
	h.settings = cli.New()
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

func (h *Helm) installCoreComponents(coreComponent config.CoreComponent) {
	h.initConfiguration(coreComponent.Namespace)

	client := action.NewInstall(h.actionConfig)
	client.Namespace = coreComponent.Namespace
	client.ReleaseName = coreComponent.ReleaseName
	client.RepoURL = coreComponent.Repository
	client.Version = coreComponent.Version
	client.CreateNamespace = true
	client.IsUpgrade = true

	chrt_path, err := client.LocateChart(coreComponent.Chart, h.settings)
	if err != nil {
		panic(err)
	}

	myChart, err := loader.Load(chrt_path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[🐶] %s is being installed...\n", coreComponent.ReleaseName)
	_, err = client.Run(myChart, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[🐶] %s(using chart version %s) successfully installed on namespace %s\n", coreComponent.ReleaseName, coreComponent.Version, coreComponent.Namespace)
}
