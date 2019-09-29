// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test files package

package docker

import (
	"testing"
)

func TestNewContainerInfos(t *testing.T) {
	//var err			error
	//var expected	string
	var containers		*ContainerInfos

	t.Log("TestNewContainerInfos()")

	containers = NewContainerInfos()
	if containers == nil {
		t.Errorf("Error: Unable to allocate ContainerInfos!\n")
	}
	//t.Logf("\tcontainers: %+v\n", containers.Containers())

	t.Log("\tend: TestNewContainerInfos")
}

func TestNewImageInfos(t *testing.T) {
	//var err			error
	//var expected	string
	var imgs		*ImageInfos

	t.Log("TestNewImageInfos()")

	imgs = NewImageInfos()
	if imgs == nil {
		t.Errorf("Error: Unable to allocate ImageInfos!\n")
	}
	t.Logf("\timages: %+v\n", imgs.Images())

	t.Log("\tend: TestNewImageInfos")
}

func TestPullImage(t *testing.T) {
	var err			error
	//var expected	string
	var imgs		*ImageInfos

	t.Log("TestNewImageInfos()")
	t.Log("\tWarning: This may take a little bit of time!")

	imgs = NewImageInfos()
	if imgs == nil {
		t.Errorf("Error: Unable to allocate ImageInfos!\n")
	}

	err = imgs.PullImage("docker", "")
	if err != nil {
		t.Errorf("Error: Unable to pull docker: %s\n", err.Error())
	}

	err = imgs.RemoveImage("docker", "")
	if err != nil {
		t.Errorf("Error: Unable to remove docker: %s\n", err.Error())
	}

	t.Log("\tend: TestNewImageInfos")
}

