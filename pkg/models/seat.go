package models

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
	"vulcanlabs-assignment/pkg/db"
	"vulcanlabs-assignment/pkg/utils"
	pb "vulcanlabs-assignment/proto"
)

type SeatStatus int

const (
	Available SeatStatus = iota
	Reserved
	Blocked
)

type Seat struct {
	Row        int        `json:"row"`
	Col        int        `json:"col"`
	Blocked    bool       `json:"blocked"`
	Reserved   bool       `json:"reserved"`
	ReservedAt *time.Time `json:"reserved_at"`
}

func (s *Seat) Bytes() ([]byte, error) {
	return json.Marshal(s)
}
func (s *Seat) FromBytes(b []byte) error {
	return json.Unmarshal(b, s)
}

type SeatMap struct {
	Info        []*Seat `json:"info,omitempty"`
	MinDistance int     `json:"min_distance,omitempty"`
	NumRows     int32   `json:"num_rows,omitempty"`
	NumCols     int32   `json:"num_cols,omitempty"`
}

func (m *SeatMap) Bytes() ([]byte, error) {
	return json.Marshal(m)
}
func (m *SeatMap) FromBytes(b []byte) error {
	return json.Unmarshal(b, m)
}
func (m *SeatMap) GetStatusMap() (map[int32]SeatStatus, []int32) {
	statusMap := make(map[int32]SeatStatus)
	availables := make([]int32, 0)
	for i, seat := range m.Info {
		statusMap[int32(i)] = seat.Status()
		if seat.Status() == Available {
			availables = append(availables, int32(i))
		}
	}
	return statusMap, availables
}

func NewSeatMap(r, c int, minDistance *int, reset bool) (*SeatMap, map[int32]SeatStatus, []int32) {
	if !reset {
		s := new(SeatMap)
		err := db.GetKey("seatmap", s)
		if err == nil {
			seatMap, availables := s.GetStatusMap()
			return s, seatMap, availables
		}
		utils.BaseLogger.Info("seatmap recovery encountered failed, creating new one instead")

	}
	m := make([]*Seat, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			s := &Seat{Row: i, Col: j, Reserved: false, ReservedAt: nil}
			m[i*c+j] = s
		}
	}
	distance := 0
	if minDistance != nil {
		distance = *minDistance
	}
	// fill status map
	seatStatusMap := make(map[int32]SeatStatus)
	utils.FillMap(seatStatusMap, int32(r*c), Available)
	availableSeat := utils.SliceWithSize(int32(r * c))
	return &SeatMap{Info: m, MinDistance: distance, NumRows: int32(r), NumCols: int32(c)}, seatStatusMap, availableSeat
}

func Distance(a, b *Seat) int {
	d := math.Abs(float64(a.Row-b.Row)) + math.Abs(float64(a.Col-b.Col))
	//logger.Debugf("|%d-%d| + |%d-%d| = %v \n", a.Row, b.Row, a.Col, b.Col, math.Abs(float64(a.Col-b.Col)))
	return int(d)
}

func (s *Seat) Equal(other *Seat) bool {
	return s.Row == other.Row && s.Col == other.Col
}

func (s *Seat) Proto() *pb.Seat {
	return &pb.Seat{
		Row:      int32(s.Row),
		Col:      int32(s.Col),
		Reserved: s.Reserved,
	}
}
func (s *Seat) Status() SeatStatus {
	if s.Reserved {
		return Reserved
	}
	if s.Blocked {
		return Blocked
	}
	return Available
}

func (m *SeatMap) Proto() *pb.SeatMap {
	result := &pb.SeatMap{
		Seats:       make([]*pb.Seat, 0),
		MinDistance: int32(m.MinDistance),
		NumRow:      int32(m.NumRows),
		NumCol:      int32(m.NumCols),
	}
	for _, c := range m.Info {
		result.Seats = append(result.Seats, &pb.Seat{
			Row:      int32(c.Row),
			Col:      int32(c.Col),
			Reserved: c.Reserved,
			Blocked:  c.Blocked,
		})
	}
	return result
}
func (m *SeatMap) DoReserve(seat *Seat) error {
	position := seat.Row*int(m.NumRows) + seat.Col
	fmt.Printf("DoReserve position:%v,%v,%v\n", position, seat.Row, seat.Col)
	if position > len(m.Info) {
		return fmt.Errorf("position out of range")
	}
	m.Info[position].Reserved = true
	now := time.Now()
	m.Info[position].ReservedAt = &now

	return nil
}
