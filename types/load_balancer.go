package types

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNoOwner = errors.New("load balancer does not have an owner")
)

/* LB Apps Table represents DB relationship of LBs and apps */
// do not change the tags, they're snake_case on purpose
type LbApp struct {
	LbID  string `json:"lb_id"`
	AppID string `json:"app_id"`
}

/* Load Balancers Table */
type (
	LoadBalancer struct {
		ID                string         `json:"id"`
		Name              string         `json:"name"`
		UserID            string         `json:"userID"`
		ApplicationIDs    []string       `json:"applicationIDs,omitempty"`
		RequestTimeout    int            `json:"requestTimeout"`
		Gigastake         bool           `json:"gigastake"`
		GigastakeRedirect bool           `json:"gigastakeRedirect"`
		StickyOptions     StickyOptions  `json:"stickinessOptions"`
		Applications      []*Application `json:"applications"`
		Users             []UserAccess   `json:"users"`
		CreatedAt         time.Time      `json:"createdAt"`
		UpdatedAt         time.Time      `json:"updatedAt"`
	}
	StickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    bool     `json:"stickiness"`
	}
	UserAccess struct {
		ID       string   `json:"id,omitempty"`
		UserID   string   `json:"userID"`
		RoleName RoleName `json:"roleName"`
		Email    string   `json:"email"`
		Accepted bool     `json:"accepted"`
	}
	/* Update structs */
	UpdateLoadBalancer struct {
		Name          string               `json:"name,omitempty"`
		StickyOptions *UpdateStickyOptions `json:"stickinessOptions,omitempty"`
		Remove        bool                 `json:"remove,omitempty"`
	}
	UpdateStickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    *bool    `json:"stickiness"`
	}
	UpdateUserAccess struct {
		ID       string   `json:"id,omitempty"`
		UserID   string   `json:"userID"`
		RoleName RoleName `json:"roleName"`
	}

	RoleName        string
	PermissionsEnum string
)

const (
	RoleOwner  RoleName = "OWNER"
	RoleAdmin  RoleName = "ADMIN"
	RoleMember RoleName = "MEMBER"

	ReadEndpoint  PermissionsEnum = "read:endpoint"
	WriteEndpoint PermissionsEnum = "write:endpoint"
)

var (
	ValidRoleNames = map[RoleName]bool{
		RoleOwner:  true,
		RoleAdmin:  true,
		RoleMember: true,
	}

	ValidPermissions = map[PermissionsEnum]bool{
		ReadEndpoint:  true,
		WriteEndpoint: true,
	}
)

func (lb *LoadBalancer) GetOwnerEmail() (string, error) {
	for _, user := range lb.Users {
		if user.RoleName == RoleOwner {
			return user.Email, nil
		}
	}

	return "", ErrNoOwner
}

func (e *PermissionsEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PermissionsEnum(s)
	case string:
		*e = PermissionsEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for PermissionsEnum: %T", src)
	}
	return nil
}

func (s *StickyOptions) IsEmpty() bool {
	if !s.Stickiness {
		return true
	}
	return len(s.StickyOrigins) == 0
}
