package entity

import "time"

type Song struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;column:id"`
	ArtistID    uint      `gorm:"index;not null;column:artist_id"`
	Name        string    `gorm:"size:200;not null;column:name"`
	Album       string    `gorm:"size:200;column:album"`
	Lyric       string    `gorm:"type:text;column:lyric"`
	Duration    string    `gorm:"size:10;column:duration"` // mm:ss 或 ss
	Style       string    `gorm:"size:100;column:style"`
	CoverURL    string    `gorm:"size:500;column:cover_url"`
	AudioURL    string    `gorm:"size:500;column:audio_url"`
	ReleaseTime time.Time `gorm:"type:date;column:release_time"`
}

func (Song) TableName() string { return "tb_song" }
