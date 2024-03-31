package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type SubjectService struct {
}

var SubjectServiceApp = new(SubjectService)

func (SubjectService) FindAndUseSubjectByTopic(topic string) string {
	subject := &ai.Subject{}

	err := global.GVA_DB.Where("topic=?", topic).Where("use_times=0").Last(&subject).Error
	if err != nil {
		subject.UseTimes = subject.UseTimes + 1
		global.GVA_DB.Save(subject)
	}
	return subject.Subject
}
