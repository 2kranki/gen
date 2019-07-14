// Using MySQL in Docker for local testing In Go
//
// This automates the starting of a MySQL docker container.
//
// Found:   https://blog.kowalczyk.info/book/go-cookbook.html
// License: Public Domain


/******
To figure out --format argument, use:
	docker container <command and parameters> --format='{{json .}}'

 ******/


package docker

import (
	"../util"
	"strings"
)

//----------------------------------------------------------------------------
//								ContainerInfo
//----------------------------------------------------------------------------

// name:tag is the external identifier and id is the internal identifier for
// docker images.
type ContainerInfo struct {
	id       	string
	imageName   string
	imageTag	string
	name     	string
	status		string
}

func (c *ContainerInfo) Id( ) string {
	return c.id
}

func (c *ContainerInfo) SetId(s string) {
	c.id = s
}

func (c *ContainerInfo) ImageName( ) string {
	return c.imageName
}

func (c *ContainerInfo) SetImageName(s string) {
	c.imageName = s
}

func (c *ContainerInfo) ImageTag( ) string {
	return c.imageTag
}

func (c *ContainerInfo) SetImageTag(s string) {
	c.imageTag = s
}

func (c *ContainerInfo) Name( ) string {
	return c.name
}

func (c *ContainerInfo) SetName(s string) {
	c.name = s
}

func (c *ContainerInfo) Status( ) string {
	return c.status
}

func (c *ContainerInfo) SetStatus(s string) {
	c.status = s
}

func (i *ContainerInfo) Setup( ) {

	cmd := util.NewExecCmd("docker", "container","inspect", i.id, "--format", "{{json .}}")
	s, err := cmd.RunWithOutput()
	util.PanicIfErr(err, "Error - docker container ps -a failed with %s", err)
	if len(s) > 0 {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "|")
			util.PanicIf(len(parts) != 4, "Unexpected output from docker ps:\n%s\n. Expected 4 parts, got %d (%v)\n", line, len(parts), parts)
			ii := NewContainerInfo()
			ii.SetId(parts[0])
			ii.SetName(parts[1])
			ii.SetTag(parts[2])
			//FIXME: i.containers = append(i.containers, ii)
		}
	}
}

func (c *ContainerInfo) String() string {
	return c.id + "|" + c.name + ":" + c.tag + " " + c.status
}

func NewContainerInfo( ) *ContainerInfo {
	container := &ContainerInfo{}
	return container
}

//----------------------------------------------------------------------------
//							ContainerInfos
//----------------------------------------------------------------------------

type ContainerInfos struct {
	containers  	[]*ContainerInfo
}

func (c *ContainerInfos) Containers( ) []*ContainerInfo {
	return c.containers
}

func (c *ContainerInfos) FindImage(name, tag string) *ContainerInfo {

	for _, img := range c.containers {
		if name == img.Name() {
			if len(tag) > 0 {
				if tag == img.Tag() {
					return img
				}
			} else {
				return img
			}
		}
	}

	return nil
}

func (i *ContainerInfos) PullImage(name, tag string) error {
	var err		error
	var ii		*ContainerInfo

	ii = i.FindImage(name, tag)
	if ii == nil {
		nameTag := name
		if len(tag) > 0 {
			nameTag += ":" + tag
		}
		cmd := util.NewExecCmd("docker", "container", "pull", nameTag)
		err = cmd.Run()
		if err == nil {
			i.Setup()
		}
	}

	return err
}

func (i *ContainerInfos) RemoveImage(name, tag string) error {
	var err		error
	var ii		*ContainerInfo

	ii = i.FindImage(name, tag)
	if ii != nil {
		nameTag := name
		if len(tag) > 0 {
			nameTag += ":" + tag
		}
		cmd := util.NewExecCmd("docker", "container","rm", nameTag)
		err = cmd.Run()
		if err == nil {
			i.Setup()
		}
	}

	return err
}

func (i *ContainerInfos) Setup( ) {

	cmd := util.NewExecCmd("docker", "container","ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Labels}}|{{.Status}}")
	s, err := cmd.RunWithOutput()
	util.PanicIfErr(err, "Error - docker container ps -a failed with %s", err)
	if len(s) > 0 {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "|")
			util.PanicIf(len(parts) != 3, "Unexpected output from docker ps:\n%s\n. Expected 4 parts, got %d (%v)\n", line, len(parts), parts)
			ii := NewContainerInfo()
			ii.SetId(parts[0])
			ii.SetName(parts[1])
			ii.SetTag(parts[2])
			i.containers = append(i.containers, ii)
		}
	}
}

func NewContainerInfos( ) *ContainerInfos {
	images := &ContainerInfos{}
	images.Setup()
	return images
}






