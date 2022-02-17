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

	_, err = CreateYamlResource(app, yamlImageStream, "all")
	if err != nil {
		return errs.Wrap(err, "Cannot create resource for yamlImageStream")
	}

	_, err = CreateYamlResource(app, yamlPVC, "all")
	if err != nil {
		return errs.Wrap(err, "Cannot create resource for yamlPVC")
	}

	return nil
}

func CreateBaseImage(app *specfemv1.SpecfemApp) error {
	return CreateImage(app, yamlBaseImageBuildConfig, "specfem:base", "all")
}

func CreateSpecfemImage(app *specfemv1.SpecfemApp) error {
	return CreateImage(app, yamlSpecfemImageBuildConfig, "specfem:specfem", "all")
}

func CreateProjectImage(app *specfemv1.SpecfemApp) error {
	return CreateImage(app, yamlProjectImageBuildConfig, "specfem:project", "all")
}

func RunDecomposeMesherJob(app *specfemv1.SpecfemApp) error {
	if _, err := CreateYamlResource(app, yamlDecomposeMeshScript, "decompose"); err != nil {
		return err
	}

	return RunSeqJob(app, yamlRunDecomposeMesherJob, "decompose")
}

func RunGenerateDbMpiJob(app *specfemv1.SpecfemApp) error {
	if _, err := CreateYamlResource(app, yamlGenerateDbScript, "generate-db"); err != nil {
		return err
	}

	return RunMpiJob(app, yamlRunGenerateDbMpiJob, "generate-db")
}

func RunSetupSymlinksJob(app *specfemv1.SpecfemApp) error {
	if _, err := CreateYamlResource(app, yamlSetupSymlinksScript, "symlinks"); err != nil {
		return err
	}

	return RunSeqJob(app, yamlRunSetupSymlinksJob, "symlinks")
}

func RunSolverMpiJob(app *specfemv1.SpecfemApp) error {
	if _, err := CreateYamlResource(app, yamlSolverScript, "solver"); err != nil {
		return err
	}
	return RunMpiJob(app, yamlRunSolverMpiJob, "solver")
}

func RunSaveSolverOutput(app *specfemv1.SpecfemApp) error {
	jobName, err := CreateYamlResource(app, yamlSaveSolverOutputJob, "solver")
	if err != nil {
		return err
	}

	if delete_mode {
		CleanupJobPods(app, yamlSaveSolverOutputJob)
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

    SAVELOG_FILENAME := fmt.Sprintf("/tmp/specfem.solver-%dproc-%dcores-%dnex_%s.log",
		app.Spec.Exec.Nproc, app.Spec.Exec.Ncore, app.Spec.Specfem.Nex, date_uid)

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

func RunSeqJob(app *specfemv1.SpecfemApp, yamlSpecFct YamlResourceSpec, stage string) error {
	jobName, err := CreateYamlResource(app, yamlSpecFct, stage)
	if err != nil || jobName == "" {
		return err
	}

	if err = WaitWithJobLogs(jobName, "", nil); err != nil {
		return err
	}

	return nil
}

func RunMpiJob(app *specfemv1.SpecfemApp, yamlSpecFct YamlResourceSpec, stage string) error {
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
	mpijobName, err := CreateYamlResource(app, yamlSpecFct, stage)
	if err != nil || mpijobName == "" {
		return err
	}

	if err = WaitMpiJob(mpijobName); err != nil {
		return err
	}

	log.Printf("MPI %s done!", stage)

	return nil
}
