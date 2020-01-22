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

    sheet.SetColWidth(2, len(schedule)*3+1, COL_W)
    for i := uint8(0); i < colNum; i++ {
        row[i] = sheet.AddRow()
        row[i].SetHeight(ROW_H)
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
            if trow.Teacher[0] == trow.Teacher[1] {
                cell.Value = trow.Teacher[0]
            } else {
                if trow.Teacher[0] != "" {
                    cell.Value = trow.Teacher[0]
                }
                if trow.Teacher[1] != "" {
                    cell.Value += "\n" + trow.Teacher[1]
                }
            }
            
            cell = row[j+1].AddCell()
            if trow.Subject[0] == trow.Subject[1] {
                cell.Value = trow.Teacher[0]
            } else {
                if trow.Subject[0] != "" {
                    cell.Value = trow.Subject[0] + " (A)"
                }
                if trow.Subject[1] != "" {
                    cell.Value += "\n" + trow.Subject[1] + " (B)"
                }
            }

            sheet.SetColWidth(colWidthForCabinets(i))
            cell = row[j+1].AddCell()
            if trow.Cabinet[0] == trow.Cabinet[1] {
                cell.Value = trow.Cabinet[0]
            } else {
                if trow.Cabinet[0] != "" {
                    cell.Value = trow.Cabinet[0]
                }
                if trow.Cabinet[1] != "" {
                    cell.Value += "\n" + trow.Cabinet[1]
                }
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
    var col = (index+1)*3+1
    return col, col, COL_W_CAB
}
