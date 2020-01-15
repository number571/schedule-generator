package schedule

func (gen *Generator) teacherToReserved(teacher string) {
	if _, ok := gen.Reserved.Teachers[teacher]; ok {
		return
	}
	gen.Reserved.Teachers[teacher] = make([]bool, gen.NumTables)
}

func (gen *Generator) cabinetToReserved(cabnum string) {
	if _, ok := gen.Reserved.Cabinets[cabnum]; ok {
		return
	}
	gen.Reserved.Cabinets[cabnum] = make([]bool, gen.NumTables)
}
