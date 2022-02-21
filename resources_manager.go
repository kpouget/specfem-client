package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	errs "github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	specfemv1 "github.com/openshift-psap/specfem-client/apis/specfem/v1alpha1"
)

func PrepareNamespace(app *specfemv1.SpecfemApp) error {
	_, err := CreateYamlResource(app, yamlNamespace, "---")
	if err != nil {
		return errs.Wrap(err, "Cannot create resource for yamlNamespace")
	}

	_, err = CreateYamlResource(app, yamlImageStream, "namespace")
	if err != nil {
		return errs.Wrap(err, "Cannot create resource for yamlImageStream")
	}

	_, err = CreateYamlResource(app, yamlPVC, "namespace")
	if err != nil {
		return errs.Wrap(err, "Cannot create resource for yamlPVC")
	}

	return nil
}


func RunSaveSolverOutput(app *specfemv1.SpecfemApp) error {
	stage := "solver"
	script_name := "run_save_solver_output.sh"

	yamlConfigMap := func() (string, YamlResourceTmpl) {
		return yamlSingleFileConfigMap(script_name)
	}

	yamlJob := func() (string, YamlResourceTmpl) {
		return yamlScriptJob(script_name)
	}

	if _, err := CreateYamlResource(app, yamlConfigMap, stage); err != nil {
		return err
	}

	jobName, err := CreateYamlResource(app, yamlJob, stage)
	if err != nil {
		return err
	}

	if delete_mode {
		CleanupJobPods(app, yamlJob)
	}

	if jobName == "" {
		return nil
	}

	var logs *string = nil
	err = WaitWithJobLogs(jobName, "", &logs)
	if err != nil {
		return err
	}
	if logs == nil {
		return fmt.Errorf("Failed to get logs for job/%s", jobName)
	}

	date_uid := time.Now().Format("20060102_150405")

    SAVELOG_FILENAME := fmt.Sprintf("/tmp/specfem_%s.log", date_uid)

	output_f, err := os.Create(SAVELOG_FILENAME)

	if err != nil {
		return err
	}

	defer output_f.Close()

	output_f.WriteString(*logs)

	log.Printf("Saved solver logs into '%s'", SAVELOG_FILENAME)

	return nil
}

/* --- */

func CreateImage(app *specfemv1.SpecfemApp, yamlSpecFct YamlResourceSpec, imageName, stage string) error {

	if ! delete_mode {
		if err := CheckImageTag(imageName, stage); err == nil {
			log.Printf("Found " + imageName + " image, don't build it.")
			return nil
		}
	}

	if err := CreateAndWaitYamlBuildConfig(app, yamlSpecFct, stage); err != nil {
		return err
	}

	if err := CheckImageTag(imageName, stage); err != nil {
		return err
	}

	return nil

}

func HasMpiWorkerPods(app *specfemv1.SpecfemApp, stage string) (int, error) {
	pods, err := client.ClientSet.CoreV1().Pods(app.ObjectMeta.Namespace).List(context.TODO(),
		metav1.ListOptions{LabelSelector: "mpi_role_type=worker,mpi_job_name=mpi-"+stage})

	if err != nil {
		return 0, err
	}

	return len(pods.Items), nil
}

func RunScriptJob(app *specfemv1.SpecfemApp, script_name, stage string) error {
	yamlConfigMap := func() (string, YamlResourceTmpl) {
		return yamlSingleFileConfigMap(script_name)
	}

	yamlJob := func() (string, YamlResourceTmpl) {
		return yamlScriptJob(script_name)
	}

	_, err := CreateYamlResource(app, yamlConfigMap, stage)
	if err != nil {
		return err
	}

	jobName, err := CreateYamlResource(app, yamlJob, stage)
	if err != nil || jobName == "" {
		return err
	}

	if err = WaitWithJobLogs(jobName, "", nil); err != nil {
		return err
	}

	return nil
}

func RunScriptMpiJob(app *specfemv1.SpecfemApp, script_name, stage string) error {
	yamlConfigMap := func() (string, YamlResourceTmpl) {
		return yamlSingleFileConfigMap(script_name)
	}

	yamlMpiJob := func() (string, YamlResourceTmpl) {
		return yamlMpiScriptJob(script_name)
	}

	_, err := CreateYamlResource(app, yamlConfigMap, stage)
	if err != nil {
		return err
	}

	if delete_mode {
		goto skip_pod_check
	}
	for {
		pod_cnt, err := HasMpiWorkerPods(app, stage)
		if err != nil {
			return err
		}

		fmt.Printf("found %d worker pods from previous mpijob/mpi-%s ...\n", pod_cnt, stage)
		if pod_cnt == 0 {
			break
		}

		time.Sleep(2 * time.Second)
		// loop
	}

skip_pod_check:
	mpijobName, err := CreateYamlResource(app, yamlMpiJob, stage)
	if err != nil || mpijobName == "" {
		return err
	}

	if err = WaitMpiJob(mpijobName); err != nil {
		return err
	}

	log.Printf("MPI %s done!", stage)

	return nil
}
