package main

import (
	"strings"
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
	ConfigMap struct {
		Name string
		Filename string
	}
	Job struct {
		Name string
		ConfigMap string
		Entrypoint string
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

type ImageResource struct {
	Image string
	YamlSpec YamlResourceSpec
}

var ImageResources = []ImageResource{
	ImageResource{"specfem:base", yamlBaseImageBuildConfig},
	ImageResource{"specfem:specfem", yamlSpecfemImageBuildConfig},
	ImageResource{"specfem:project", yamlProjectImageBuildConfig},
}

type StageResource struct {
	Stage string
	RunFunction func(app *specfemv1.SpecfemApp, filename, stage string) error
	Script string
}

var StageResources = []StageResource{
	StageResource{"decompose", RunScriptJob, "run_decompose_mesher.sh"},
	StageResource{"generate-db", RunScriptMpiJob, "run_mpi_generate_db.sh"},
	StageResource{"setup-symlinks", RunScriptJob, "run_setup_symlinks.sh"},
	StageResource{"solver", RunScriptMpiJob, "run_mpi_solver.sh"},
}

// --- //

func yamlSingleFileConfigMap(filename string) (string, YamlResourceTmpl) {
	return "999_configmap_file.yaml", func(app *specfemv1.SpecfemApp) *TemplateCfg {
		cfg := &TemplateCfg{}

		cfg.ConfigMap.Name = stringToName(filename)
		cfg.ConfigMap.Filename = filename
		return cfg
	}
}

func yamlScriptJob(filename string) (string, YamlResourceTmpl) {
	return "999_job_template.yaml", func(app *specfemv1.SpecfemApp) *TemplateCfg {
		cfg := &TemplateCfg{}

		cfg.Job.Entrypoint = filename
		cfg.Job.ConfigMap = stringToName(filename)
		cfg.Job.Name = strings.TrimSuffix(cfg.Job.ConfigMap, "-sh")

		return cfg
	}
}

func yamlMpiScriptJob(filename string) (string, YamlResourceTmpl) {
	return "999_mpijob_template.yaml", func(app *specfemv1.SpecfemApp) *TemplateCfg {
		cfg := &TemplateCfg{}

		cfg.Job.Entrypoint = filename
		cfg.Job.ConfigMap = stringToName(filename)
		cfg.Job.Name = strings.TrimSuffix(cfg.Job.ConfigMap, "-sh")

		return cfg
	}
}

func stringToName(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "_", "-"), ".", "-")
}
