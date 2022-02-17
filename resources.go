package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	specfemv1 "github.com/openshift-psap/specfem-client/apis/specfem/v1alpha1"
)

var podResource         = schema.GroupVersionResource{Version: "v1", Resource: "pods"}
var pvcResource         = schema.GroupVersionResource{Version: "v1", Resource: "persistentvolumeclaims"}
var cmResource          = schema.GroupVersionResource{Version: "v1", Resource: "configmaps"}
var svcResource         = schema.GroupVersionResource{Version: "v1", Resource: "services"}
var jobResource         = schema.GroupVersionResource{Version: "v1", Resource: "jobs",         Group: "batch"}
var buildconfigResource = schema.GroupVersionResource{Version: "v1", Resource: "buildconfigs", Group: "build.openshift.io"}
var imagestreamResource    = schema.GroupVersionResource{Version: "v1", Resource: "imagestreams",    Group: "image.openshift.io"}
var imagestreamtagResource = schema.GroupVersionResource{Version: "v1", Resource: "imagestreamtags", Group: "image.openshift.io"}

var routeResource       = schema.GroupVersionResource{Version: "v1", Resource: "routes",       Group: "route.openshift.io"}

var secretResource       = schema.GroupVersionResource{Version: "v1", Resource: "secrets"}

var NAMESPACE = ""

var USE_UBI_BASE_IMAGE = true

type TemplateCfg struct {
	ConfigMaps struct {
		HelperFile struct {
			ConfigMapName string
			ManifestName string
		}
	}
	SecretNames struct {
		DockerCfgPush string
	}
	MesherSolver struct {
		Stage string
		Image string
		Nreplicas int
	}
	BaseImage struct {
		Image string
	}
}

func yamlNamespace() (string, YamlResourceTmpl) {
	return "000_namespace.yaml", NoTemplateCfg
}

func yamlImageStream() (string, YamlResourceTmpl) {
	return "001_imagestream.yaml", NoTemplateCfg
}

func yamlPVC() (string, YamlResourceTmpl) {
	return "002_pvc.yaml", NoTemplateCfg
}

func yamlBaseImageBuildConfig() (string, YamlResourceTmpl) {
	return "010_buildconfig_base.yaml", func(app *specfemv1.SpecfemApp) *TemplateCfg {
		cfg := &TemplateCfg{}

		var img string
		if app.Spec.Resources.UseUbiImage {
			if app.Spec.Specfem.GpuPlatform != "" {
				img = "docker.io/nvidia/cuda:11.1-devel-ubi8"
			} else {
				img = "registry.access.redhat.com/ubi8/ubi"
			}
		} else {
			img = "docker.io/ubuntu:eoan"
		}

		cfg.BaseImage.Image = img

		return cfg
	}
}

func yamlSpecfemImageBuildConfig() (string, YamlResourceTmpl) {
	return "011_buildconfig_specfem.yaml", NoTemplateCfg
}

func yamlProjectImageBuildConfig() (string, YamlResourceTmpl) {
	return "012_buildconfig_project.yaml", NoTemplateCfg
}

func yamlRunDecomposeMesherJob() (string, YamlResourceTmpl) {
	return "020_job_decompose_mesh.yaml", NoTemplateCfg
}

func yamlRunGenerateDbMpiJob() (string, YamlResourceTmpl) {
	return "021_mpijob_generate_db.yaml", NoTemplateCfg
}

func yamlRunSetupSymlinksJob() (string, YamlResourceTmpl) {
	return "022_job_setup_symlinks.yaml", NoTemplateCfg
}

func yamlRunSolverMpiJob() (string, YamlResourceTmpl) {
	return "023_mpijob_solver.yaml", NoTemplateCfg
}

func yamlSaveSolverOutputJob() (string, YamlResourceTmpl) {
	return "030_job_save-solver-output.yaml", NoTemplateCfg
}

func yamlDecomposeMeshScript() (string, YamlResourceTmpl) {
	return yamlFileConfigMap("run_decompose_mesh.sh")
}

func yamlGenerateDbScript() (string, YamlResourceTmpl) {
	return yamlFileConfigMap("run_generate_db.sh")
}

func yamlSetupSymlinksScript() (string, YamlResourceTmpl) {
	return yamlFileConfigMap("run_setup_symlinks.sh")
}

func yamlSolverScript() (string, YamlResourceTmpl) {
	return yamlFileConfigMap("run_solver.sh")
}

func yamlFileConfigMap(filename string) (string, YamlResourceTmpl) {
	return "999_configmap_file.yaml", func(app *specfemv1.SpecfemApp) *TemplateCfg {
		cfg := &TemplateCfg{}
		cfg.ConfigMaps.HelperFile.ManifestName = filename
		return cfg
	}
}
