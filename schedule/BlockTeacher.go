package schedule

import (
	"errors"
)

func (gen *Generator) BlockTeacher(teacher string) error {
	if !gen.InTeachers(teacher) {
		return errors.New("teacher undefined")
	}
	gen.Blocked[teacher] = true
	return nil
}
