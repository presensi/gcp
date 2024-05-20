package helper

import (
	"strings"

	"github.com/gocroot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WebHook(WAKeyword, WAPhoneNumber, WAAPIQRLogin, WAAPIMessage string, msg model.IteungMessage, db *mongo.Database) (resp model.Response, err error) {
	if IsLoginRequest(msg, WAKeyword) { //untuk whatsauth request login
		resp, err = HandlerQRLogin(msg, WAKeyword, WAPhoneNumber, db, WAAPIQRLogin)
	} else { //untuk membalas pesan masuk
		resp, err = HandlerIncomingMessage(msg, WAPhoneNumber, db, WAAPIMessage)
	}
	return
}

func RefreshToken(dt *model.WebHook, WAPhoneNumber, WAAPIGetToken string, db *mongo.Database) (res *mongo.UpdateResult, err error) {
	profile, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	var resp model.User
	if profile.Token != "" {
		resp, err = PostStructWithToken[model.User]("Token", profile.Token, dt, WAAPIGetToken)
		if err != nil {
			return
		}
		profile.Phonenumber = resp.PhoneNumber
		profile.Token = resp.Token
		res, err = ReplaceOneDoc(db, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
		if err != nil {
			return
		}
	}
	return
}

func IsLoginRequest(msg model.IteungMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) // && msg.From_link
}

func GetUUID(msg model.IteungMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg model.IteungMessage, WAKeyword string, WAPhoneNumber string, db *mongo.Database, WAAPIQRLogin string) (resp model.Response, err error) {
	dt := &model.WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Delay:       msg.From_link_delay,
	}
	structtoken, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	resp, err = PostStructWithToken[model.Response]("Token", structtoken.Token, dt, WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg model.IteungMessage, WAPhoneNumber string, db *mongo.Database, WAAPIMessage string) (resp model.Response, err error) {
	dt := &model.TextMessage{
		To:       msg.Chat_number,
		IsGroup:  false,
		Messages: GetRandomReplyFromMongo(msg, db),
	}
	if msg.Chat_server == "g.us" { //jika pesan datang dari group maka balas ke group
		dt.IsGroup = true
	}
	botnumber, err := IsBotNumber(msg.Phone_number, db)
	if err != nil {
		return
	}
	if botnumber {
		dt.Messages = dt.Messages + " terdeteksi nomor bot"
	} else {
		dt.Messages = dt.Messages + "bukan nomor bot"
	}
	if !botnumber { //ignore pesan datang dari sesama bot di profile
		var profile model.Profile
		profile, err = GetAppProfile(WAPhoneNumber, db)
		if err != nil {
			return
		}
		resp, err = PostStructWithToken[model.Response]("Token", profile.Token, dt, WAAPIMessage)
		if err != nil {
			return
		}
	}
	return
}

func IsBotNumber(phonenumber string, db *mongo.Database) (status bool, err error) {
	profilebot, err := GetAppProfile(phonenumber, db)
	if err != nil {
		return
	}
	if (profilebot != model.Profile{}) {
		return true, nil
	}
	return
}

func GetRandomReplyFromMongo(msg model.IteungMessage, db *mongo.Database) string {
	rply, err := GetRandomDoc[model.Reply](db, "reply", 1)
	if err != nil {
		return "Koneksi Database Gagal: " + err.Error()
	}
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", msg.Alias_name)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func GetAppProfile(phonenumber string, db *mongo.Database) (apitoken model.Profile, err error) {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, err = GetOneDoc[model.Profile](db, "profile", filter)

	return
}
