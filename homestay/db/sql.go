package db

import "time"

// 民宿收入明细表
type HomestayIncomeDetail struct {
	ID         int       `json:"id" gorm:"id"`                   // 序号
	Channel    string    `json:"channel" gorm:"channel"`         // 渠道
	Room       int       `json:"room" gorm:"room"`               // 房间
	Money      float64   `json:"money" gorm:"money"`             // 收入
	IncomeTime time.Time `json:"income_time" gorm:"income_time"` // 收入时间
}

// 获取民宿收入详情数据
func (db *HomeStayDB) GetHomestayIncomeDetailData(startTime, endTime string) ([]*HomestayIncomeDetail, error) {
	var details []*HomestayIncomeDetail
	sql := `select 
				a.id,b.name as channel,c.number as room,a.money,a.income_time
			from 
				homestay_income_detail a 
			left join 
 				homestay_channel b on a.channel_id = b.id 
			left join 
				homestay_room c on a.room_id = c.id where a.income_time between ? and ?
			order by a.channel_id,a.room_id`
	if err := db.DB.Raw(sql, startTime, endTime).Scan(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}
