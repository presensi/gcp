package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/idname"
	"github.com/gocroot/mod/lldikti"
	"github.com/gocroot/mod/panduan"

	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Modulename string, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "idname-masuk":
		reply = idname.IDNameMasuk(Pesan, db)
	case "idname-pulang":
		reply = idname.IDNamePulang(Pesan, db)
	case "lldikti":
		reply = lldikti.Lldikti(Pesan)
	case "panduan":
		reply = panduan.Panduan(Pesan)
	}
	return
}
