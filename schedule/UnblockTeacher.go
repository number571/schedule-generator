package schedule

import (
	"errors"
)

func (gen *Generator) UnblockTeacher(teacher string) error {
	if !gen.InBlocked(teacher) {
		return errors.New("teacher undefined")
	}
	delete(gen.Blocked, teacher)
	return nil
}
