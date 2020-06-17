package homestay

import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"

// The key of Generators is the prefix of table info url.
// The corresponding value is the Form and Table data.
//
// http://{{config.Domain}}:{{Port}}/{{config.Prefix}}/info/{{key}}
//
// example:
//
// "homestay_channel" => http://localhost:9033/admin/info/homestay_channel
//
// "homestay_room" => http://localhost:9033/admin/info/homestay_room
//
// "homestay_room" => http://localhost:9033/admin/info/homestay_room
//
// "homestay_income_detail" => http://localhost:9033/admin/info/homestay_income_detail
// "homestay_spend_detail" => http://localhost:9033/admin/info/homestay_spend_detail
// "homestay_spend_item" => http://localhost:9033/admin/info/homestay_spend_item
//
// example end
//
var Generators = map[string]table.Generator{
	"homestay_channel": GetHomestayChannelTable,

	"homestay_room": GetHomestayRoomTable,

	"homestay_income_detail": GetHomestayIncomeDetailTable,
	"homestay_spend_detail":  GetHomestaySpendDetailTable,
	"homestay_spend_item":    GetHomestaySpendItemTable,

	// generators end
}
