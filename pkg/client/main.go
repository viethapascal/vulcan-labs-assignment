package client

import (
	context "context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"slices"
	"time"
	"vulcanlabs-assignment/pkg/utils"
	pb "vulcanlabs-assignment/proto"
)

type CineClient struct {
	logger *zap.SugaredLogger
	c      pb.SeatReservationServiceClient
	ctx    context.Context
	conn   *grpc.ClientConn
}

func NewCineClient(port int) (*CineClient, error) {
	logger := utils.Logger("cine-client")
	logger.Infof("Connecting to Cine at port %v", port)
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Errorf("Failed to create client: %v", err)
		return nil, err
	}
	seatClient := pb.NewSeatReservationServiceClient(conn)
	return &CineClient{
		ctx:    context.Background(),
		logger: logger,
		c:      seatClient,
		conn:   conn,
	}, nil

}

func (c *CineClient) GetSeatMap() (*pb.SeatMap, error) {
	ctx := context.Background()
	//defer cancel()
	resp, err := c.c.GetSeatMap(ctx, &pb.Empty{})
	if err != nil {
		c.logger.Errorw("Failed to get seat map", "error", err)
		return nil, err
	}
	data := make([][]string, resp.NumRow)
	// headers
	headers := make([]string, resp.NumCol+1)
	for i := range resp.NumCol {
		headers[i+1] = fmt.Sprintf("Column %d", i)
	}
	//data = append(data, headers)
	chunked := slices.Chunk(resp.Seats, int(resp.NumCol))
	noNum := 0
	for r := range chunked {
		row := make([]string, resp.NumCol+1)
		row[0] = fmt.Sprintf("ROW %d", noNum)
		for idx, j := range r {
			if j.Reserved {
				row[idx+1] = "R"
			} else if j.Blocked {
				row[idx+1] = "B"
			} else {
				row[idx+1] = "-"
			}
		}
		data[noNum] = row

		noNum++
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetRowLine(true) // Enable row line
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetRowSeparator("-")
	table.AppendBulk(data)
	table.Render()
	fmt.Println("R: Reserved; B: Blocked; -: Available")
	return resp, err
}

func (c *CineClient) Reserve(input [][2]int32) (*pb.SeatMap, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()
	payload := make([]*pb.Seat, 0)
	for _, d := range input {
		payload = append(payload, &pb.Seat{Row: d[0], Col: d[1]})
	}
	req := &pb.ReserveRequest{
		Seat: payload,
	}
	_, err := c.c.Reserve(ctx, req)
	if err != nil {
		c.logger.Infof("Failed to reserve:. %v. Error: %v", input, err.Error())
		return nil, err
	}
	c.logger.Infow("Successfully reserved", "seats", input)
	return nil, nil
}

func (c *CineClient) Close() {
	c.conn.Close()
}
