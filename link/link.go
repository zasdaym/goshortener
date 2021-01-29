package link

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type link struct {
	ID        primitive.ObjectID
	Slug      string
	Target    string
	CreatedAt time.Time
}

type service interface {
	getLinks(ctx context.Context, req getLinksRequest) (links []*link, err error)
	createLink(ctx context.Context, req createLinkRequest) (err error)
}

type getLinksRequest struct {
	Limit int64 `query:"limit"`
	Page  int64 `query:"page"`
}

type createLinkRequest struct {
	Slug   string `json:"slug"`
	Target string `json:"target"`
}

func (req *createLinkRequest) validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.Slug, validation.Required),
		validation.Field(&req.Target, validation.Required),
	)
}
