package services

import (
	//"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	//"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	//"github.com/revel/revel"
	//uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type RoleService struct {
}

func (srv *RoleService) Getlargestrole(role []models.Role) int {
	if len(role) < 1 {
		return 0
	}
	max := role[0].Weight // assume first value is the smallest
	for _, value := range role {
		if value.Weight > max {
			max = value.Weight // found another smaller value, replace previous value in max
		}
	}

	largestWeight:=int(max)

	return largestWeight
}
