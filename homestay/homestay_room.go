package homestay

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetHomestayRoomTable(ctx *context.Context) table.Table {

	homestayRoomTable := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := homestayRoomTable.GetInfo()

	info.AddField("Id", "id", db.Int)
	info.AddField("Number", "number", db.Int)
	info.AddField("Address", "address", db.Varchar)
	info.AddField("Monthly_rent", "monthly_rent", db.Int)
	info.AddField("Monthly_manage_fee", "monthly_manage_fee", db.Int)
	info.AddField("Start_time", "start_time", db.Timestamp)
	info.AddField("End_time", "end_time", db.Timestamp)

	info.SetTable("homestay_room").SetTitle("民宿房间").SetDescription("民宿房间")

	formList := homestayRoomTable.GetForm()

	formList.AddField("Id", "id", db.Int, form.Default).FieldNotAllowAdd()
	formList.AddField("Number", "number", db.Int, form.Number)
	formList.AddField("Address", "address", db.Varchar, form.Text)
	formList.AddField("Monthly_rent", "monthly_rent", db.Int, form.Number)
	formList.AddField("Monthly_manage_fee", "monthly_manage_fee", db.Int, form.Number)
	formList.AddField("Start_time", "start_time", db.Timestamp, form.Datetime)
	formList.AddField("End_time", "end_time", db.Timestamp, form.Datetime)

	formList.SetTable("homestay_room").SetTitle("民宿房间").SetDescription("民宿房间")

	return homestayRoomTable
}
