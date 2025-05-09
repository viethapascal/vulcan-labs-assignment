package server

import (
	context "context"
	"errors"
	"go.uber.org/zap"
	"slices"
	"sync"
	"vulcanlabs-assignment/pkg/db"
	"vulcanlabs-assignment/pkg/models"
	"vulcanlabs-assignment/pkg/utils"
	pb "vulcanlabs-assignment/proto"
)

/*
A[M*n]
A[i,j] => position = i * M + j
*/

type ReservationServer struct {
	pb.UnimplementedSeatReservationServiceServer
	seatMap       *models.SeatMap
	seatStatusMap map[int32]models.SeatStatus
	mtx           sync.Mutex
	logger        *zap.SugaredLogger
	availableSeat []int32
}

func NewReservationServer(rows, cols, distance int, reset bool) *ReservationServer {
	logger := utils.Logger("grpc-server")
	logger.Info("Starting new ReservationServer")
	seatMap, seatStatusMap, avail := models.NewSeatMap(rows, cols, &distance, reset)
	// save seatmap
	db.SetKey("seatmap", seatMap)
	return &ReservationServer{
		seatMap:       seatMap,
		mtx:           sync.Mutex{},
		logger:        logger,
		seatStatusMap: seatStatusMap,
		availableSeat: avail,
	}
}

func (r *ReservationServer) GetAvailableSeats(ctx context.Context, empty *pb.Empty) (*pb.ReserveResponse, error) {
	data := make([]*pb.Seat, 0)
	for p, status := range r.seatStatusMap {
		if status == models.Available {
			data = append(data, &pb.Seat{
				Row:      p / r.seatMap.NumCols,
				Col:      p % r.seatMap.NumCols,
				Reserved: false,
			})
		}
	}
	return &pb.ReserveResponse{
		Success: true,
		Message: "GetAvailableSeats",
		Data:    data,
	}, nil
}

func (r *ReservationServer) SaveState() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	err := db.SetKey("seatmap", r.seatMap)
	if err != nil {
		r.logger.Warnw("Error saving seatmap", "err", err)
	}
	r.logger.Info("Saved seatmap")
}

func (r *ReservationServer) Reserve(ctx context.Context, req *pb.ReserveRequest) (*pb.ReserveResponse, error) {
	r.logger.Infof("Received Reserve. %v", req.Seat)
	defer r.logger.Info("Finished Reserve")
	defer r.SaveState()
	/* Reservation Logic
	- Check if Available Seat remain >= len(req.Seat)
	- Update reserved seats to Reserved status
	- Update block status
	*/
	if len(req.Seat) == 0 {
		return nil, errors.New("seat must not be empty")
	}
	if len(r.availableSeat) < len(req.Seat) {
		return nil, errors.New("no available seat")
	}
	toBeBlocked := make([]int32, 0)
	for _, seat := range req.Seat {
		if (seat.Row < 0 || seat.Row >= r.seatMap.NumCols) || (seat.Col < 0 || seat.Col >= r.seatMap.NumCols) {
			return nil, errors.New("invalid seat")
		}
		p := seat.Row*r.seatMap.NumCols + seat.Col
		if r.seatMap.Info[p].Blocked || r.seatMap.Info[p].Reserved {
			return nil, errors.New("seat already reserved or blocked")
		}

		// Update Seat status
		if err := r.seatMap.DoReserve(&models.Seat{Col: int(seat.Col), Row: int(seat.Row)}); err != nil {
			return nil, err
		}
		// Update Seatmap
		r.seatStatusMap[p] = models.Reserved
		// Block all point within Manhattan radius
		blockedList := utils.PointsWithinManhattan(int(seat.Row), int(seat.Col), r.seatMap.MinDistance, r.seatMap.NumCols)
		r.logger.Info("Blocked seat ", blockedList)
		toBeBlocked = append(toBeBlocked, blockedList...)
		// Ignore reserved seat
		slices.DeleteFunc(toBeBlocked, func(i int32) bool {
			return i == p
		})
	}

	for _, b := range toBeBlocked {
		slices.DeleteFunc(r.availableSeat, func(v int32) bool {
			return v == b
		})
		if b >= 0 && b < int32(r.seatMap.NumRows*r.seatMap.NumCols) {
			r.seatMap.Info[b].Blocked = true
		}
	}
	r.logger.Info("New avail seat ", r.availableSeat)

	//
	//r.logger.Infof("Reserving: %d - %d. Available: %d\n", seat.Row, seat.Col, len(r.availableSeat))
	//input := &models.Seat{Row: int(seat.Row), Col: int(seat.Col)}
	//// d = x1-x2 + y1-y2
	// (p/N) - x2 + P%N -y2 <= d
	// A-x + B-y <= d
	//tobeRemove := make([]int, 0)
	//for i, s := range r.availableSeat {
	//	//fmt.Println("Seat:", s.Row, s.Col, models.Distance(s, input))
	//	if r.seatMap.MinDistance == 0 {
	//		if s.Equal(input) {
	//			tobeRemove = append(tobeRemove, i)
	//			break
	//		}
	//	} else {
	//		if s.Equal(input) {
	//			tobeRemove = append(tobeRemove, i)
	//		}
	//		if !s.Equal(input) && models.Distance(s, input) >= r.seatMap.MinDistance {
	//			tobeRemove = append(tobeRemove, i)
	//		}
	//	}
	//}
	//// Update available seat
	//newList := make([]*models.Seat, 0)
	//for i, t := range tobeRemove {
	//
	//	if i == t {
	//		newList = append(newList, r.availableSeat[i])
	//	}
	//}
	//r.availableSeat = newList
	return &pb.ReserveResponse{
		Success: true,
		Message: "Reserved",
		Data:    nil,
	}, nil
}

func (r *ReservationServer) GetSeatMap(ctx context.Context, empty *pb.Empty) (*pb.SeatMap, error) {
	r.logger.Info("GetSeatMap")
	//convert to pb.SeatMap
	return r.seatMap.Proto(), nil
}
