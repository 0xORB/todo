package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Tasks []task

var (
	colTitleIndex       = "ID"
	colTitleComplete    = "Complete"
	colTitleTask        = "Task"
	colTitleCreatedAt   = "Created"
	colTitleCompletedAt = "Finished"
	rowHeader           = table.Row{
		colTitleIndex,
		colTitleComplete,
		colTitleTask,
		colTitleCreatedAt,
		colTitleCompletedAt,
	}
)

type task struct {
	ID          string
	Completed   bool
	Description string
	CreatedAt   time.Time
	CompletedAt time.Time
}

func (t *Tasks) CreateTask(desc string) {
	task := task{
		Completed:   false,
		Description: desc,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*t = append(*t, task)
}

func (t *Tasks) RemoveTask(idx int) error {
	ls := *t
	if err := checkIndex(idx, ls); err != nil {
		return errors.New("Failed to delete task.")
	}
	*t = append(ls[:idx-1], ls[idx:]...)
	return nil
}

func (t *Tasks) MarkTaskComplete(idx int) error {
	ls := *t
	if err := checkIndex(idx, ls); err != nil {
		return err
	}
	ls[idx-1].CompletedAt = time.Now()
	ls[idx-1].Completed = true
	return nil
}

func checkIndex(idx int, t []task) error {
	if idx < 1 || idx > len(t) {
		return errors.New("Invalid Index.")
	}
	return nil
}

func (t *Tasks) LoadTasks(fn string) error {
	f, err := os.ReadFile(fn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.New("Failed to open file.")
	}
	if len(f) == 0 {
		return errors.New("File is empty.")
	}
	err = json.Unmarshal(f, t)
	if err != nil {
		return errors.New("Failed to parse json.")
	}
	return nil
}

func (t *Tasks) StoreTasks(fn string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return errors.New("Failed to marshal task data.")
	}
	os.WriteFile(fn, data, 0644)
	return nil
}

func (t *Tasks) PrintTasks() {
	var (
		mark     string
		compMark interface{}
	)

	tw := table.NewWriter()
	tw.SetTitle("List of Tasks")
	tw.SetIndexColumn(1)
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: colTitleComplete, Align: text.AlignCenter},
	})
	tw.SetStyle(table.StyleBold)
	tw.Style().Options.DrawBorder = true
	tw.Style().Options.SeparateRows = true

	tw.AppendHeader(rowHeader)
	for idx, val := range *t {
		if !val.Completed {
			mark = "_"
			compMark = "-"
		} else {
			mark = "X"
			compMark = val.CompletedAt.Format(time.RFC850)
		}
		tw.AppendRow(
			table.Row{
				idx + 1,
				mark,
				val.Description,
				val.CreatedAt.Format(time.RFC850),
				compMark,
			},
		)
	}
	fmt.Println(tw.Render())
}
