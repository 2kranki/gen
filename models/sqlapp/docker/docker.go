// Docker supplies functions and objects to deal with Docker
// images and containers dynamically.
//
// This automates the starting of a MySQL docker container.
//
// Found:   https://blog.kowalczyk.info/book/go-cookbook.html
// License: Public Domain


/******
docker image ls -a --format "{{json .}}"

 ******/


package docker

// To run:
// go run main.g

import (
	"../util"
	"strings"
)

const (
	dockerStatusExited  = "exited"
	dockerStatusRunning = "running"
)

var (
	// using https://hub.docker.com/_/mysql/
	// to use the latest mysql, use mysql:8
	DockerImageName = "mysql:5.6"
	// name must be unique across containers runing on this computer
	DockerContainerName = "mysql-db-multi"
	// where mysql stores databases. Must be on local disk so that
	// database outlives the container
	DockerDbDir = "~/data/db-multi"
	// 3306 is standard MySQL port, I use a unique port to be able
	// to run multiple mysql instances for different projects
	DockerDbLocalPort = "7200"
)

//----------------------------------------------------------------------------
//									Docker
//----------------------------------------------------------------------------

// name:tag is the external identifier and id is the internal identifier for
// docker images.
type Docker struct {
}

func (d *Docker) ClientVersion( ) string {
	cmd := util.NewExecCmd("docker", "version", "--format", "{{.Client.Version}}")
	outBytes, err := cmd.RunWithOutput()
	util.PanicIfErr(err, "cmd.RunWithOutput() for '%s' failed with %s", cmd.CommandString(), err)
	s := string(outBytes)
	s = strings.TrimSpace(s)

	return s
}

func NewDocker( ) *Docker {
	d := Docker{}
	return &d
}

