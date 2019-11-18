package docker

import "github.com/ory/dockertest/v3"

import "github.com/caspr-io/mu-kit/util"

type Docker struct {
	pool       *dockertest.Pool
	containers []*Container
}

type Container struct {
	*dockertest.Resource
	docker *Docker
}

func StartDocker() (*Docker, error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	return &Docker{pool, []*Container{}}, nil
}

func (d *Docker) RunContainer(image string, version string, env []string) (*Container, error) {
	// pulls the postgres image, creates a container based on it and runs it
	resource, err := d.pool.Run(image, version, env)
	if err != nil {
		return nil, err
	}

	c := &Container{resource, d}
	d.containers = append(d.containers, c)

	return c, nil
}

func (c *Container) WaitForRunning(waitFunc func() error) error {
	return c.docker.pool.Retry(waitFunc)
}

func (d *Docker) Close() error {
	err := new(util.ErrorCollector)

	for _, c := range d.containers {
		e := c.Close()
		if e != nil {
			err.Collect(e)
		}
	}

	return err
}
