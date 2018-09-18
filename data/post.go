package data

import (
	"git.cory.red/DankBotList/data"
	"github.com/jinzhu/gorm"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Post struct {
	gorm.Model
	PostURLID  uint64 `gorm:"not null;unique_index"`
	Title      string `gorm:"not null;unique_index" sql:"index"`
	Author     User   `gorm:"-"`
	URL        string `gorm:"not null;unique_index"`
	Summary    string
	FullTextID uint `gorm:"not null;unique_index"`
	Published  bool
}

type PostText struct {
	gorm.Model
	OldVersion uint   `gorm:"not null;unique_index"`
	FullText   string `gorm:"type:LONGTEXT"`
}

func (p *Post) GetFullText(h *data.Handler) (string, uint, error) {
	var postText PostText
	err := h.Engine.First(&postText, "id = ?", p.FullTextID).Error
	if err != nil {
		return "", 0, err
	}

	if postText.OldVersion != 0 {
		var patches [][]diffmatchpatch.Patch
		dmp := diffmatchpatch.New()

		for postText.OldVersion != 0 {
			err := h.Engine.First(&postText, "id = ?", postText.OldVersion).Error
			if err != nil {
				return "", 0, err
			}

			if postText.OldVersion != 0 {
				patch, err := dmp.PatchFromText(postText.FullText)
				if err != nil {
					return "", 0, err
				}
				patches = append(patches, patch)
			}
		}

		for _, v := range patches {
			postText.FullText, _ = dmp.PatchApply(v, postText.FullText)
		}
	}

	return postText.FullText, postText.OldVersion, nil
}

func (p *Post) SetFullText(h *data.Handler, newText string) error {
	oldText, oldId, err := p.GetFullText(h)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	patchText := dmp.PatchToText(dmp.PatchMake(oldText, dmp.DiffMain(oldText, newText, true)))

	postText := PostText{
		OldVersion: oldId,
		FullText:   patchText,
	}

	if err := h.Engine.Create(&postText).Error; err != nil {
		return err
	}

	return h.Engine.Model(p).Updates(map[string]interface{}{
		"full_text_id": postText.ID,
	}).Error

}
