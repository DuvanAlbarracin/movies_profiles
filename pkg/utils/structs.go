package utils

import (
	"reflect"

	"github.com/DuvanAlbarracin/movies_profiles/pkg/models"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/proto"
)

func SetModifyMap(s any) (map[string]any, []any) {
	resMap := make(map[string]any)
	var args []any
	ref := reflect.ValueOf(s)

	for i := 0; i < ref.NumField(); i++ {
		f := ref.Field(i)
		if f.IsZero() {
			continue
		}

		if f.Kind() == reflect.Pointer {
			f = f.Elem()
		}

		v := ref.Type().Field(i)
		fValue := f.Interface()
		resMap[v.Name] = fValue
		args = append(args, fValue)
	}

	return resMap, args
}

func SetProtoProfile(profile models.Profile) *proto.Profile {
	var g proto.Gender
	r := proto.Role(*profile.Role)

	pp := proto.Profile{
		Id:        profile.Id,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       profile.Age,
		City:      profile.City,
		Role:      &r,
	}

	if profile.Gender != nil {
		g = proto.Gender(*profile.Gender)
		pp.Gender = &g
	}

	return &pp
}
