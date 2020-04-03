// Code generated by svcdec; DO NOT EDIT.

package api

import (
	"context"

	proto "github.com/golang/protobuf/proto"
)

type DecoratedRegistration struct {
	// Service is the service to decorate.
	Service RegistrationServer
	// Prelude is called for each method before forwarding the call to Service.
	// If Prelude returns an error, then the call is skipped and the error is
	// processed via the Postlude (if one is defined), or it is returned directly.
	Prelude func(ctx context.Context, methodName string, req proto.Message) (context.Context, error)
	// Postlude is called for each method after Service has processed the call, or
	// after the Prelude has returned an error. This takes the the Service's
	// response proto (which may be nil) and/or any error. The decorated
	// service will return the response (possibly mutated) and error that Postlude
	// returns.
	Postlude func(ctx context.Context, methodName string, rsp proto.Message, err error) error
}

func (s *DecoratedRegistration) CreateMachines(ctx context.Context, req *MachineList) (rsp *MachineResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "CreateMachines", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.CreateMachines(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "CreateMachines", rsp, err)
	}
	return
}

func (s *DecoratedRegistration) GetMachines(ctx context.Context, req *EntityIDList) (rsp *MachineResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "GetMachines", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.GetMachines(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "GetMachines", rsp, err)
	}
	return
}

func (s *DecoratedRegistration) ListMachines(ctx context.Context, req *ListMachinesRequest) (rsp *MachineResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "ListMachines", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.ListMachines(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "ListMachines", rsp, err)
	}
	return
}

func (s *DecoratedRegistration) UpdateMachines(ctx context.Context, req *MachineList) (rsp *MachineResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "UpdateMachines", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.UpdateMachines(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "UpdateMachines", rsp, err)
	}
	return
}

func (s *DecoratedRegistration) DeleteMachines(ctx context.Context, req *EntityIDList) (rsp *EntityIDResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "DeleteMachines", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.DeleteMachines(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "DeleteMachines", rsp, err)
	}
	return
}
