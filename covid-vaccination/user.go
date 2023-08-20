package covidvaccination

import "fmt"



type UserService struct {
	UserData         []*User
	UserCount        int
	UserIdToIndexMap map[string]int
}

func GetUserService( userSize int ) *UserService{
	return &UserService{
		UserData: make([]*User, userSize),
		UserCount: 0,
		UserIdToIndexMap: make(map[string]int),
	}
}

type User struct {
	Id         string
	Name       string
	Age        int
	DistrictId string
	HasBooking bool
	BookingId  string
}

func (s *UserService) RegisterUser(user *User) error {
	s.UserData[s.UserCount] = user
	s.UserIdToIndexMap[user.Id] = s.UserCount
	s.UserCount++
	return nil
}

func (s *UserService) GetUser(userId string) (*User, error) {
	if _, ok := s.UserIdToIndexMap[userId]; !ok {
		return nil, fmt.Errorf("user not found")
	}
	user := s.UserData[s.UserIdToIndexMap[userId]]
	return user, nil
}
