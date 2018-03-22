package models

//TimelineUser .
type TimelineUser struct {
	UserProfile
	CreatedAt int64  `json:"create_at" gorm:"column:create_at"`
	Duration  uint64 `json:"recent_duration" gorm:"column:recent_duration"`
}

//TableName .
func (TimelineUser) TableName() string {
	return "time_line"
}

//QueryAll .
func (t *TimelineUser) QueryAll(gender, skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser
	return tl, db.Offset(skip).Limit(limit).Where("gender = ?", gender).Find(&tl).Error
}

//QueryRecent .
func (t *TimelineUser) QueryRecent(timestamp int64, gender, skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser

	if gender == 1 {
		return tl, db.Offset(skip).Limit(limit).Where("create_at > ? and gender = ?", timestamp, gender).Find(&tl).Error
	}
	return tl, db.Offset(skip).Limit(limit).Where("create_at > ? and gender = ? and user_type = ?", timestamp, gender, UserTypeAnchor).Find(&tl).Error
}

//QueryHot .
func (t *TimelineUser) QueryHot(skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser
	return tl, db.Offset(skip).Limit(limit).Where("user_status = ? and gender = ?", UserStatusHot, GenderWoman).Find(&tl).Error
}
