package models

import "gorm.io/gorm"

type UserSegment struct {
	gorm.Model
	SegmentID int
	Segment   Segment
	UserID    int
	User      User
}
