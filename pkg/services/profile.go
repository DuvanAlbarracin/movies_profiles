package services

import (
	"context"
	"net/http"

	"github.com/DuvanAlbarracin/movies_profiles/pkg/db"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/models"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/proto"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	H db.Handler
	proto.UnimplementedProfilesServiceServer
}

func (s *Server) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	var g int
	var lNSanitized string
	var cSanitized string

	r := int(req.Role)
	fNSanitized := utils.SanitizeString(req.FirstName)

	profile := models.Profile{
		FirstName: &fNSanitized,
		Role:      &r,
	}

	if req.LastName != nil {
		lNSanitized = utils.SanitizeString(*req.LastName)
		profile.LastName = &lNSanitized
	}

	if req.Age != nil {
		profile.Age = req.Age
	}

	if req.Gender != nil {
		g = int(*req.Gender)
		profile.Gender = &g
	}

	if req.City != nil {
		cSanitized = utils.SanitizeString(*req.City)
		profile.City = &cSanitized
	}

	if err := db.CreateProfile(s.H.Conn, &profile); err != nil {
		st := status.New(codes.Internal, "Cannot create the user")
		return nil, st.Err()
	}

	return &proto.CreateResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	_, err := db.FindProfileById(s.H.Conn, req.Id)
	if err != nil {
		var code codes.Code
		var message string
		if err.Error() == "no rows in result set" {
			code = codes.NotFound
			message = "Profile not found"
		} else {
			code = codes.Internal
			message = "There was an error finding the profile"
		}

		return nil, status.New(code, message).Err()
	}

	if err = db.DeleteProfile(s.H.Conn, req.Id); err != nil {
		st := status.New(codes.Internal, "There was an error deleting the profile")
		return nil, st.Err()
	}

	return &proto.DeleteResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) Modify(ctx context.Context, req *proto.ModifyRequest) (*proto.ModifyResponse, error) {
	var g int
	var r int

	reqProfile := models.Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		City:      req.City,
	}
	if req.Gender != nil {
		g = int(*req.Gender)
		reqProfile.Gender = &g
	} else {
		reqProfile.Gender = nil
	}

	if req.Role != nil {
		r = int(*req.Role)
		reqProfile.Role = &r
	} else {
		reqProfile.Role = nil
	}

	profileMap, args := utils.SetModifyMap(reqProfile)
	if err := db.ModityProfile(s.H.Conn, req.Id, profileMap, args); err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &proto.ModifyResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetById(ctx context.Context, req *proto.GetByIdRequest) (*proto.GetByIdResponse, error) {
	profile, err := db.FindProfileById(s.H.Conn, req.Id)
	if err != nil {
		var code codes.Code
		var message string
		if err.Error() == "no rows in result set" {
			code = codes.NotFound
			message = "Profile not found"
		} else {
			code = codes.Internal
			message = "There was an error finding the profile"
		}

		return nil, status.New(code, message).Err()
	}

	pp := utils.SetProtoProfile(profile)
	var resp proto.GetByIdResponse
	resp.Profile = pp

	return &resp, nil
}

func (s *Server) GetAll(ctx context.Context, req *proto.GetAllRequest) (*proto.GetAllResponse, error) {
	profiles, err := db.GetAllProfiles(s.H.Conn)
	if err != nil {
		st := status.New(codes.Internal, "There was an error getting all profiles: "+err.Error())
		return nil, st.Err()
	}

	if len(profiles) == 0 {
		st := status.New(codes.NotFound, "There is no profiles")
		return nil, st.Err()
	}

	var protoProfiles []*proto.Profile
	for _, profile := range profiles {
		protoProfile := utils.SetProtoProfile(profile)
		protoProfiles = append(protoProfiles, protoProfile)
	}

	return &proto.GetAllResponse{
		Profiles: protoProfiles,
	}, nil
}
