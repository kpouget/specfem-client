package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"os"
	"time"

	specfemv1 "github.com/openshift-psap/specfem-client/apis/specfem/v1alpha1"
	errs "github.com/pkg/errors"

)

var ArtifactDir string

func SaveArtifact(app *specfemv1.SpecfemApp, dirname, filename string, content []byte) error {
	dir, err := getArtifactDir(app, dirname)
	if err != nil {
		return err
	}

	dest := dir+"/"+filename
	log.Printf("artifacts: Saving %s ...\n", dest)
	err = os.WriteFile(dest, content, 0644)
	if err != nil {
        return errs.Wrap(err, fmt.Sprintf("Could not write the artifact file '%s'", dest))
	}

	return nil
}

func getArtifactDir(app *specfemv1.SpecfemApp, dirname string) (string, error) {
	if ArtifactDir != "" {
		return ArtifactDir, nil
	}
	currentTime := time.Now()

	parentDir := fmt.Sprintf("/tmp/specfem-client_%s", currentTime.Format("20060602"))
	if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return "", errs.Wrap(err, fmt.Sprintf("Could not create the artifact directory '%s'", parentDir))
	}

    files, err := ioutil.ReadDir(parentDir)
    if err != nil {
        return "", errs.Wrap(err, fmt.Sprintf("Could not read the content of the artifact directory '%s'", parentDir))
    }

	cnt := 0
    for _, file := range files {
		if ! file.IsDir() {
			continue
		}
		cnt += 1
    }
	ArtifactDir = fmt.Sprintf("%s/%03d_specfem/%s", parentDir, cnt, dirname)
	if err := os.MkdirAll(ArtifactDir, os.ModePerm); err != nil {
		return "", errs.Wrap(err, fmt.Sprintf("Could not create the artifact directory '%s'", ArtifactDir))
	}

	return ArtifactDir, nil
}
