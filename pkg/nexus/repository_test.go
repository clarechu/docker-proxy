package nexus

import (
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_List(t *testing.T) {
	app := models.NexusApp{}
	rep := NewRepository(app)
	ports, err := rep.GetPortByDocker()
	assert.Equal(t, nil, err)
	assert.Equal(t, []int{9001}, ports)
}
