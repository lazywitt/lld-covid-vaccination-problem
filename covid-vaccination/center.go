package covidvaccination

import "fmt"

type CenterService struct {
	CenterData           []*Center
	CenterCount          int
	CenterIdToIndexMap   map[string]int
	DistrictIdToIndexMap map[string][]int
}

func GetCenterService(bookingSize int) *CenterService {
	return &CenterService{
		CenterData:           make([]*Center, bookingSize),
		CenterIdToIndexMap:   make(map[string]int),
		DistrictIdToIndexMap: make(map[string][]int),
	}
}

type Center struct {
	Id                    string
	DefaultCapacity       int
	DistrictId            string
	DistrictName          string
	StateName             string
	BookingCapacityPerDay map[int]int
	BookingCountPerDay    map[int]int
}

func (s *CenterService) RegisterCenter(center *Center) error {
	s.CenterData[s.CenterCount] = center
	s.CenterIdToIndexMap[center.Id] = s.CenterCount
	s.DistrictIdToIndexMap[center.DistrictId] = append(s.DistrictIdToIndexMap[center.DistrictId], s.CenterCount)
	s.CenterCount++
	return nil
}

func (s *CenterService) GetCentersFromDistrict(districtId string) ([]Center, error) {
	centersId := s.DistrictIdToIndexMap[districtId]
	centers := []Center{}
	for i := range centersId {
		centers = append(centers, *s.CenterData[i])
	}
	return centers, nil
}

func (s *CenterService) GetCenter(centerId string) (*Center, error) {
	if index, ok := s.CenterIdToIndexMap[centerId]; ok {
		return s.CenterData[index], nil
	}
	return nil, fmt.Errorf("center not found")
}
