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

//QueryAll 活跃专区
func (t *TimelineUser) QueryAll(faker bool, gender, skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser

	q := db
	if faker {
		q = q.Where("user_type = ?", UserTypeFaker)
	} else {
		q = q.Where("gender = ?", gender)

		if gender == GenderWoman {
			q = q.Where("user_type = ?", UserTypeAnchor)
		} else {
			q = q.Where("user_type <> ?", UserTypeFaker)
		}
	}
	return tl, q.Order("online_status = 1 or online_status = 2 desc, recent_duration desc").Offset(skip).Limit(limit).Find(&tl).Error
}

//QueryRecent 新人专区
func (t *TimelineUser) QueryRecent(faker bool, timestamp int64, gender, skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser

	q := db

	if faker {
		q = q.Where("user_type = ?", UserTypeFaker)
	} else {
		q = q.Where("gender = ?", gender).Where("create_at > ?", timestamp)

		if gender == GenderWoman {
			q = q.Where("user_type = ?", UserTypeAnchor)
		} else {
			q = q.Where("user_type <> ?", UserTypeFaker)
		}
	}
	return tl, q.Order("online_status = 1 or online_status = 2 desc, recent_duration desc").Offset(skip).Limit(limit).Find(&tl).Error
}

//QuerySuggestion 推荐专区
func (t *TimelineUser) QuerySuggestion(faker bool, gender, skip, limit int) ([]TimelineUser, error) {
	var tl []TimelineUser

	q := db.Where("user_status = ?", UserStatusHot)

	if faker {
		q = q.Where("user_type = ?", UserTypeFaker)
	} else {
		q = q.Where("gender = ?", gender)

		if gender == GenderWoman {
			q = q.Where("user_type = ?", UserTypeAnchor)
		} else {
			q = q.Where("user_type <> ?", UserTypeFaker)
		}
	}
	return tl, q.Order("online_status = 1 or online_status = 2 desc, recent_duration desc").Offset(skip).Limit(limit).Find(&tl).Error
}
