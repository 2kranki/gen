// Using MySQL in Docker for local testing In Go
//
// This automates the starting of a MySQL docker container.
//
// Found:   https://blog.kowalczyk.info/book/go-cookbook.html
// License: Public Domain

/******
To figure out --format argument, use:
	docker image <command and parameters> --format='{{json .}}'

docker image ls -a --format "{{json .}}"
	generates multiple lines of JSON

 ******/

package docker

import (
	"../util"
	"strings"
)

//----------------------------------------------------------------------------
//								ImageInfo
//----------------------------------------------------------------------------

// name:tag is the external identifier and id is the internal identifier for
// docker images.
type ImageInfo struct {
	id   string
	name string
	tag  string
}

func (i *ImageInfo) Id() string {
	return i.id
}

func (i *ImageInfo) SetId(s string) {
	i.id = s
}

func (i *ImageInfo) Name() string {
	return i.name
}

func (i *ImageInfo) SetName(s string) {
	i.name = s
}

func (i *ImageInfo) Tag() string {
	return i.tag
}

func (i *ImageInfo) SetTag(s string) {
	i.tag = s
}

func (i *ImageInfo) String() string {
	return i.id + "|" + i.name + ":" + i.tag
}

func NewImageInfo() *ImageInfo {
	image := ImageInfo{}
	return &image
}

//----------------------------------------------------------------------------
//								ImageInfos
//----------------------------------------------------------------------------

type ImageInfos struct {
	images []*ImageInfo
}

func (i *ImageInfos) Images() []*ImageInfo {
	return i.images
}

func (i *ImageInfos) FindImage(name, tag string) *ImageInfo {

	for _, img := range i.images {
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

func (i *ImageInfos) PullImage(name, tag string) error {
	var err error
	var ii *ImageInfo

	ii = i.FindImage(name, tag)
	if ii == nil {
		nameTag := name
		if len(tag) > 0 {
			nameTag += ":" + tag
		}
		cmd := util.NewExecCmd("docker", "image", "pull", nameTag)
		err = cmd.Run()
		if err == nil {
			i.Setup()
		}
	}

	return err
}

func (i *ImageInfos) RemoveImage(name, tag string) error {
	var err error
	var ii *ImageInfo

	ii = i.FindImage(name, tag)
	if ii != nil {
		nameTag := name
		if len(tag) > 0 {
			nameTag += ":" + tag
		}
		cmd := util.NewExecCmd("docker", "image", "rm", nameTag)
		err = cmd.Run()
		if err == nil {
			i.Setup()
		}
	}

	return err
}

func (i *ImageInfos) Setup() {

	cmd := util.NewExecCmd("docker", "image", "ls", "-a", "--format", "{{.ID}}|{{.Repository}}|{{.Tag}}")
	s, err := cmd.RunWithOutput()
	util.PanicIfErr(err, "Error - docker image ls -a failed with %s", err)
	if len(s) > 0 {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "|")
			util.PanicIf(len(parts) != 3, "Unexpected output from docker ps:\n%s\n. Expected 3 parts, got %d (%v)\n", line, len(parts), parts)
			ii := NewImageInfo()
			ii.SetId(parts[0])
			ii.SetName(parts[1])
			ii.SetTag(parts[2])
			i.images = append(i.images, ii)
		}
	}
}

func NewImageInfos() *ImageInfos {
	images := &ImageInfos{}
	images.Setup()
	return images
}
