package image

import (
	"context"
	imageStoragev1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/filestorage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type serverApi struct {
	imageStoragev1.UnimplementedImageStorageServer
}

func Register(gRPC *grpc.Server) {
	imageStoragev1.RegisterImageStorageServer(gRPC, &serverApi{})
}

func (s serverApi) UploadImage(ctx context.Context, request *imageStoragev1.UploadImageRequest) (*imageStoragev1.UploadImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) DeleteImage(ctx context.Context, request *imageStoragev1.DeleteImageRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}
