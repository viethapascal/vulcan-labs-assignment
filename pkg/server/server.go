package server

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"vulcanlabs-assignment/pkg/utils"
	pb "vulcanlabs-assignment/proto"
)

func StartServer(port, rows, cols, distance int, reset bool) {
	logger := utils.Logger("server-app")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	var opt []grpc.ServerOption
	grpcServer := grpc.NewServer(opt...)
	//seatMap := models.NewSeatMap(rows, cols, &distance, reset)
	pb.RegisterSeatReservationServiceServer(grpcServer, NewReservationServer(rows, cols, distance, reset))
	logger.Infof("Starting server at :%v. Config: row = %d, col = %d, min_distance = %d\n", lis.Addr(), rows, cols, distance)

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error(err)
		return
	}
}
