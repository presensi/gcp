package lldikti

import "github.com/whatsauth/itmodel"

func Lldikti(Pesan itmodel.IteungMessage) (reply string) {

	return "Hai.. hai.. ini dia nama alias kaka:\n" + Pesan.Alias_name + "\nsimpan dan catat baik-baik ya.. \nNomor telepon kaka: " + Pesan.Phone_number + "\nmakasih"
}
