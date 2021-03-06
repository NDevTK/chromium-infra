// Code generated by svcdec; DO NOT EDIT.

package dashboard

import (
	"context"

	proto "github.com/golang/protobuf/proto"
)

type DecoratedChopsAnnouncements struct {
	// Service is the service to decorate.
	Service ChopsAnnouncementsServer
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

func (s *DecoratedChopsAnnouncements) CreateLiveAnnouncement(ctx context.Context, req *CreateLiveAnnouncementRequest) (rsp *CreateLiveAnnouncementResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "CreateLiveAnnouncement", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.CreateLiveAnnouncement(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "CreateLiveAnnouncement", rsp, err)
	}
	return
}

func (s *DecoratedChopsAnnouncements) RetireAnnouncement(ctx context.Context, req *RetireAnnouncementRequest) (rsp *Announcement, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "RetireAnnouncement", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.RetireAnnouncement(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "RetireAnnouncement", rsp, err)
	}
	return
}

func (s *DecoratedChopsAnnouncements) SearchAnnouncements(ctx context.Context, req *SearchAnnouncementsRequest) (rsp *SearchAnnouncementsResponse, err error) {
	if s.Prelude != nil {
		var newCtx context.Context
		newCtx, err = s.Prelude(ctx, "SearchAnnouncements", req)
		if err == nil {
			ctx = newCtx
		}
	}
	if err == nil {
		rsp, err = s.Service.SearchAnnouncements(ctx, req)
	}
	if s.Postlude != nil {
		err = s.Postlude(ctx, "SearchAnnouncements", rsp, err)
	}
	return
}
