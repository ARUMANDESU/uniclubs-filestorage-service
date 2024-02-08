package image

import (
	"context"
	imageStoragev1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/filestorage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	imageStoragev1.UnimplementedImageStorageServer
	imageService Image
}

func Register(gRPC *grpc.Server, imageService Image) {
	imageStoragev1.RegisterImageStorageServer(gRPC, &serverApi{imageService: imageService})
}

type Image interface {
	Upload(ctx context.Context, image []byte, filename string) (string, error)
	Delete(ctx context.Context, imageUrl string) error
}

func (s serverApi) UploadImage(ctx context.Context, req *imageStoragev1.UploadImageRequest) (*imageStoragev1.UploadImageResponse, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.Image, validation.Required),
		validation.Field(&req.Filename, validation.Required),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	url, err := s.imageService.Upload(ctx, req.GetImage(), req.GetFilename())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &imageStoragev1.UploadImageResponse{ImageUrl: url}, nil
}

func (s serverApi) DeleteImage(ctx context.Context, req *imageStoragev1.DeleteImageRequest) (*empty.Empty, error) {
	err := validation.Validate(&req.ImageUrl, validation.Required)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.imageService.Delete(ctx, req.GetImageUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
