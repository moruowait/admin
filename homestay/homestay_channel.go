package homestay

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetHomestayChannelTable(ctx *context.Context) table.Table {

	homestayChannelTable := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := homestayChannelTable.GetInfo()

	info.AddField("序号", "id", db.Int)
	info.AddField("名称", "name", db.Varchar)
	info.AddField("描述", "desc", db.Varchar)

	info.SetTable("homestay_channel").SetTitle("民宿渠道").SetDescription("民宿渠道")

	formList := homestayChannelTable.GetForm()

	formList.AddField("序号", "id", db.Int, form.Default).FieldNotAllowAdd()
	formList.AddField("名称", "name", db.Varchar, form.Text)
	formList.AddField("描述", "desc", db.Varchar, form.TextArea)

	formList.SetTable("homestay_channel").SetTitle("民宿渠道").SetDescription("民宿渠道")

	return homestayChannelTable
}
