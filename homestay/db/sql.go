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

// 民宿房间信息表
type HomestayRoom struct {
	ID               int       `json:"id" gorm:"id"`                                 // 序号
	Number           int       `json:"number" gorm:"number"`                         // 房号
	Address          string    `json:"address" gorm:"address"`                       // 地址
	MonthlyRent      float64   `json:"monthly_rent" gorm:"monthly_rent"`             // 月租
	MonthlyManageFee float64   `json:"monthly_manage_fee" gorm:"monthly_manage_fee"` // 物业费
	StartTime        time.Time `json:"start_time" gorm:"start_time"`                 // 承租时间
	EndTime          time.Time `json:"end_time" gorm:"end_time"`                     // 到期时间
}

// 获取民宿房间信息表
func (db *HomeStayDB) GetHomestayRoom() ([]*HomestayRoom, error) {
	var rooms []*HomestayRoom
	sql := `select * from homestay_room`
	if err := db.DB.Raw(sql).Scan(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// 获取民宿房间信息表
func (db *HomeStayDB) GetHomestayRoomMap() (map[int]*HomestayRoom, error) {
	var r []*HomestayRoom
	sql := `select * from homestay_room`
	if err := db.DB.Raw(sql).Scan(&r).Error; err != nil {
		return nil, err
	}
	var rooms = make(map[int]*HomestayRoom, len(r))
	for _, v := range r {
		rooms[v.Number] = v
	}
	return rooms, nil
}
