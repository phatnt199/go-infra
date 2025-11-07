package mappings

import (
	datamodel "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	dtoV1 "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/mapper"
)

func ConfigureProductsMappings() error {
	err := mapper.CreateMap[*models.User, *dtoV1.UserDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dtoV1.UserDto, *models.User]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*datamodel.UserDataModel, *models.User]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.User, *datamodel.UserDataModel]()
	if err != nil {
		return err
	}

	return nil
}
