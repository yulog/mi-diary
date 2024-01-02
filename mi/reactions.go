// https://mholt.github.io/json-to-go/ で生成後、修正
package mi

import "time"

type Reactions []struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
	Type      string    `json:"type"`
	Note      Note      `json:"note,omitempty"`
}
type AvatarDecorations struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
type Instance struct {
	Name            string `json:"name"`
	SoftwareName    string `json:"softwareName"`
	SoftwareVersion string `json:"softwareVersion"`
	IconURL         string `json:"iconUrl"`
	FaviconURL      string `json:"faviconUrl"`
	ThemeColor      string `json:"themeColor"`
}

// type Emojis struct {
// }
type User struct {
	ID                string              `json:"id"`
	Name              any                 `json:"name"`
	Username          string              `json:"username"`
	Host              any                 `json:"host"`
	AvatarURL         string              `json:"avatarUrl"`
	AvatarBlurhash    any                 `json:"avatarBlurhash"`
	AvatarDecorations []AvatarDecorations `json:"avatarDecorations"`
	IsBot             bool                `json:"isBot"`
	IsCat             bool                `json:"isCat"`
	Instance          Instance            `json:"instance,omitempty"`
	Emojis            any                 `json:"emojis"`
	OnlineStatus      string              `json:"onlineStatus"`
	BadgeRoles        []any               `json:"badgeRoles"`
}

// type User struct {
// 	AvatarDecorations []any    `json:"avatarDecorations"`
// }
// type User struct {
// 	Name              string              `json:"name"`
// 	AvatarBlurhash    string              `json:"avatarBlurhash"`
// 	BadgeRoles        []BadgeRoles        `json:"badgeRoles"`
// }
// type User struct {
// 	Name              string   `json:"name"`
// 	Host              string   `json:"host"`
// }
// type BadgeRoles struct {
// 	Name         string `json:"name"`
// 	IconURL      string `json:"iconUrl"`
// 	DisplayOrder int    `json:"displayOrder"`
// }

// type Reactions struct {
// 	HayaHantei int `json:":haya_hantei@.:"`
// }
// type ReactionEmojis struct {
// }
type Note struct {
	ID                 string    `json:"id"`
	CreatedAt          time.Time `json:"createdAt"`
	UserID             string    `json:"userId"`
	User               User      `json:"user"`
	Text               string    `json:"text"`
	Cw                 any       `json:"cw"`
	Visibility         string    `json:"visibility"`
	LocalOnly          bool      `json:"localOnly"`
	ReactionAcceptance string    `json:"reactionAcceptance"`
	RenoteCount        int       `json:"renoteCount"`
	RepliesCount       int       `json:"repliesCount"`
	Reactions          any       `json:"reactions"`
	ReactionEmojis     any       `json:"reactionEmojis"`
	Emojis             any       `json:"emojis"`
	FileIds            []any     `json:"fileIds"`
	// Files              []any     `json:"files"`
	Files        []File   `json:"files"`
	ReplyID      any      `json:"replyId"`
	Mentions     []string `json:"mentions,omitempty"`
	URI          string   `json:"uri"`
	URL          string   `json:"url"`
	RenoteID     any      `json:"renoteId"`
	ClippedCount int      `json:"clippedCount"`
	Reply        Reply    `json:"reply,omitempty"`
	MyReaction   string   `json:"myReaction"`
}

type Properties struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
type File struct {
	ID           string     `json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Md5          string     `json:"md5"`
	Size         int        `json:"size"`
	IsSensitive  bool       `json:"isSensitive"`
	Blurhash     string     `json:"blurhash"`
	Properties   Properties `json:"properties"`
	URL          string     `json:"url"`
	ThumbnailURL string     `json:"thumbnailUrl"`
	Comment      any        `json:"comment"`
	FolderID     string     `json:"folderId"`
	Folder       any        `json:"folder"`
	UserID       any        `json:"userId"`
	User         any        `json:"user"`
}

// type Note0 struct {
// 	ReactionAcceptance any       `json:"reactionAcceptance"`

// 	ReplyID            string    `json:"replyId"`
// }

// type Emojis struct {
// 	Ta     string `json:"_ta"`
// 	I      string `json:"_i"`
// 	Ni     string `json:"_ni"`
// 	Xya    string `json:"_xya"`
// 	N      string `json:"_n"`
// 	Pu     string `json:"_pu"`
// 	Github string `json:"github"`
// }

// type Reactions struct {
// 	NAMING_FAILED                     int `json:"🤯"`
// 	Sugoihanashi                      int `json:":sugoihanashi@.:"`
// 	SurprisedAi                       int `json:":surprised_ai@.:"`
// 	YosanoParty                       int `json:":yosano_party@.:"`
// 	IgyoNoKesin                       int `json:":igyo_no_kesin@.:"`
// 	Bikkuri9IneverseCom               int `json:":bikkuri@9ineverse.com:"`
// 	KokyakuGaHontouNiHituyoudattaMono int `json:":kokyaku_ga_hontou_ni_hituyoudatta_mono@.:"`
// }
// type ReactionEmojis struct {
// 	Bikkuri9IneverseCom string `json:"bikkuri@9ineverse.com"`
// }
// type Emojis struct {
// }

// type Reactions struct {
// 	SetsufuroCheering int `json:":setsufuro_cheering@.:"`
// }
// type ReactionEmojis struct {
// }
type Reply struct {
	ID                 string    `json:"id"`
	CreatedAt          time.Time `json:"createdAt"`
	UserID             string    `json:"userId"`
	User               User      `json:"user"`
	Text               string    `json:"text"`
	Cw                 any       `json:"cw"`
	Visibility         string    `json:"visibility"`
	LocalOnly          bool      `json:"localOnly"`
	ReactionAcceptance any       `json:"reactionAcceptance"`
	RenoteCount        int       `json:"renoteCount"`
	RepliesCount       int       `json:"repliesCount"`
	Reactions          any       `json:"reactions"`
	ReactionEmojis     any       `json:"reactionEmojis"`
	FileIds            []any     `json:"fileIds"`
	Files              []any     `json:"files"`
	ReplyID            string    `json:"replyId"`
	RenoteID           any       `json:"renoteId"`
	Mentions           []string  `json:"mentions"`
}

// type Reactions struct {
// 	NAMING_FAILED      int `json:"👍"`
// 	Polarbear          int `json:":polarbear@.:"`
// 	WakaruMkYopoWork   int `json:":wakaru@mk.yopo.work:"`
// 	OtokuSocialSda1Net int `json:":otoku@social.sda1.net:"`
// }
// type ReactionEmojis struct {
// 	WakaruMkYopoWork   string `json:"wakaru@mk.yopo.work"`
// 	OtokuSocialSda1Net string `json:"otoku@social.sda1.net"`
// }
