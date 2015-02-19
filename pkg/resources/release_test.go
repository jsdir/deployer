package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChangedServices(t *testing.T) {
	release := Release{
		Services: map[string]string{
			"service1": "1",
			"service2": "2",
			"service3": "foo",
			"service4": "4",
		},
	}

	changedServices := release.getChangedServices(&Release{
		Services: map[string]string{
			"service1": "1",
			"service2": "2",
			"service3": "3",
		},
	})

	assert.Len(t, changedServices, 2)
	assert.Contains(t, changedServices, "service3")
	assert.Contains(t, changedServices, "service4")
}
