package service

import (
	"context"
	"time"

	registrypb "github.com/hunkvm/locator/pkg/proto/registry"
	pbtypes "github.com/hunkvm/locator/pkg/proto/types"
	"github.com/hunkvm/locator/pkg/types"

	"google.golang.org/grpc"
)

type RegistryService struct {
	client         registrypb.ServiceRegistryServiceClient
	requestTimeout time.Duration
}

func NewRegistryService(conn *grpc.ClientConn, timeout time.Duration) *RegistryService {
	return &RegistryService{
		client:         registrypb.NewServiceRegistryServiceClient(conn),
		requestTimeout: timeout,
	}
}

func (r *RegistryService) Retrieve(
	newCtx context.Context, selector map[string][]any,
) ([]*types.Service, error) {
	newSelectors := make([]*registrypb.Selector, 0, len(selector))
	for field, values := range selector {
		vals := make([]string, 0, len(values))
		for _, val := range values {
			if v, ok := val.(string); ok {
				vals = append(vals, v)
			}
		}
		newSelectors = append(newSelectors,
			&registrypb.Selector{
				Field:  field,
				Values: vals,
			},
		)
	}
	req := registrypb.RetrieveServiceRequest{Selectors: newSelectors}
	newCtx, cancel := context.WithTimeout(newCtx, r.requestTimeout)
	defer cancel()
	res, err := r.client.Retrieve(newCtx, &req)
	if err != nil {
		return nil, err
	}
	services := make([]*types.Service, 0, len(res.Services))
	for _, srv := range res.Services {
		service := &types.Service{
			ID:       srv.Id,
			Name:     srv.Name,
			Address:  srv.Address,
			Version:  srv.Version,
			Enabled:  srv.Enabled,
			Metadata: srv.Metadata,
		}
		if hc := srv.HealthCheck; hc != nil {
			service.HealthCheck = &types.HealthConfig{
				Service: hc.Service,
				Address: hc.Address,
			}
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *RegistryService) Register(
	ctx context.Context, service types.Service,
) (string, error) {
	request := &registrypb.RegisterServiceRequest{
		Name:     service.Name,
		Address:  service.Address,
		Version:  service.Version,
		Enabled:  service.Enabled,
		Metadata: service.Metadata,
	}
	if hc := service.HealthCheck; hc != nil {
		request.HealthCheck = &pbtypes.HealthConfig{
			Service: hc.Service,
			Address: hc.Address,
		}
	}
	newCtx, cancel := context.WithTimeout(ctx, r.requestTimeout)
	defer cancel()
	response, err := r.client.Register(newCtx, request)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

func (r *RegistryService) Deregister(ctx context.Context, id string) error {
	request := &registrypb.DeregisterServiceRequest{Id: id}
	_, err := r.client.Deregister(ctx, request)
	return err
}
