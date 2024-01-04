// /users/reactions
// https://mholt.github.io/json-to-go/ で生成後、修正
package mi

import "time"

type Reactions []Reaction

type Reaction struct {
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
	Files              []File    `json:"files"`
	Tags               []string  `json:"tags"`
	ReplyID            any       `json:"replyId"`
	Mentions           []string  `json:"mentions,omitempty"`
	URI                string    `json:"uri"`
	URL                string    `json:"url"`
	RenoteID           any       `json:"renoteId"`
	ClippedCount       int       `json:"clippedCount"`
	Reply              Reply     `json:"reply,omitempty"`
	MyReaction         string    `json:"myReaction"`
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
