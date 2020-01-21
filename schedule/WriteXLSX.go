package schedule

import (
    // "fmt"
    "strconv"
    "github.com/tealeg/xlsx"
)

func (gen *Generator) WriteXLSX(filename string, schedule []*Schedule) error {
    const (
        colWidth = 30
        rowHeight = 30
    )

    var (
        colNum = gen.NumTables + 1
        row = make([]*xlsx.Row, colNum)
        cell *xlsx.Cell
    )

    file := xlsx.NewFile()
    sheet, err := file.AddSheet("Schedule")
    if err != nil {
        return err
    }

    sheet.SetColWidth(2, len(schedule)*3+1, colWidth)
    for i := uint8(0); i < colNum; i++ {
        row[i] = sheet.AddRow()
        row[i].SetHeight(rowHeight)
        cell = row[i].AddCell()
        if i == 0 {
            cell.Value = "Пара"
            continue
        }
        cell.Value = strconv.Itoa(int(i))
    }

    for i, sch := range schedule {
        savedCell := row[0].AddCell()
        savedCell.Value = "Группа " + sch.Group

        cell = row[0].AddCell()
        cell = row[0].AddCell()

        savedCell.Merge(2, 0)

        for j, trow := range sch.Table {

            cell = row[j+1].AddCell()
            switch {
            case trow.Teacher[0] == trow.Teacher[1]: cell.Value = trow.Teacher[0]
            case trow.Teacher[0] != "": cell.Value = trow.Teacher[0]
            case trow.Teacher[1] != "": cell.Value += "\n" + trow.Teacher[1]
            }

            cell = row[j+1].AddCell()
            switch {
            case trow.Subject[0] == trow.Subject[1]: cell.Value = trow.Subject[0]
            case trow.Subject[0] != "": cell.Value = trow.Subject[0] + " (A)"
            case trow.Subject[1] != "": cell.Value += "\n" + trow.Subject[1] + " (B)"
            }

            sheet.SetColWidth(colWidthForCabinets(i))
            cell = row[j+1].AddCell()
            switch {
            case trow.Cabinet[0] == trow.Cabinet[1]: cell.Value = trow.Cabinet[0]
            case trow.Cabinet[0] != "": cell.Value = trow.Cabinet[0]
            case trow.Cabinet[1] != "": cell.Value += "\n" + trow.Cabinet[1]
            }

        }
    }

    err = file.Save(filename)
    if err != nil {
        return err
    }

    return nil
}

func colWidthForCabinets(index int) (int, int, float64) {
    const colWidth = 10
    var col = (index+1)*3+1
    return col, col, colWidth
}
