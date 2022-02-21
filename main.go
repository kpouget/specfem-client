package main

import (
	"flag"
	"log"

	specfemv1 "github.com/openshift-psap/specfem-client/apis/specfem/v1alpha1"
)

var DELETE_KEYS = []string{
	"---",
	"all",
	"decompose",
	"generate-db",
	"setup-symlinks",
	"solver",
}
var delete_mode = false

var to_delete = map[string]bool{}

func initDelete() {
	for _, key := range DELETE_KEYS {
		delete_it := (key == *flag_delete) || delete_mode
		to_delete[key] = delete_it
		if delete_it {
			if !delete_mode {
				log.Println("Stages to delete:")
			}
			delete_mode = true
			log.Println("- ", key)
		}
	}

	if *flag_delete == "" {
		return
	}

	if !delete_mode {
		log.Fatalf("FATAL: wrong delete flag option: %v\n", *flag_delete)
	}
}

var flag_delete = flag.String("delete", "", "solver,mesher,config,all|none")
var flag_cfg = flag.String("config", "", "name of the config file (in config/[name].yaml)")

func main() {
	var err error
	flag.Parse()

	initDelete()

	if err = InitClient(); err != nil {
		log.Fatalf("FATAL: %+v\n", err)
	}

	if err = FetchManifests(); err != nil {
		log.Fatalf("FATAL: %+v\n", err)
	}

	var configName string
	if *flag_cfg == "" {
		configName = "specfem-sample"
	} else {
		configName = *flag_cfg
	}

	var app *specfemv1.SpecfemApp

	if app, err = getSpecfemConfig(configName); err != nil {
		log.Fatalf("FATAL: failed to get the application configuration: %+v\n", err)
	}

	if err = checkSpecfemConfig(app); err != nil {
		log.Fatalf("FATAL: config error: %+v\n", err)
	}

	NAMESPACE = app.ObjectMeta.Namespace

	if err = RunSpecfem(app); err != nil {
		log.Fatalf("FATAL: %v\n", err)
	}

	log.Println("Done :)")
}


func RunSpecfem(app *specfemv1.SpecfemApp) error {

	if err := PrepareNamespace(app); err != nil {
		return err
	}

	for _, resource := range ImageResources {
		if err := CreateImage(app, resource.YamlSpec, resource.Image, "all"); err != nil {
			return err
		}
	}

	for _, resource := range StageResources {
		if err := resource.RunFunction(app, resource.Script, resource.Stage); err != nil {
			return err
		}
	}

	if err := RunSaveSolverOutput(app); err != nil {
		return err
	}

	log.Println("All done!")

	return nil
}
