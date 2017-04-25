package storage

import (
	"fmt"
	"os"
	"encoding/csv"
	"github.com/aziule/tasks/task"
	"strconv"
	"io"
	"io/ioutil"
)

const FILE_NAME = "tasks.csv"

func init() {
	if _, err := os.Stat(FILE_NAME); err != nil {
		file, err := os.Create(FILE_NAME);

		if err != nil {
			fmt.Println("Impossible to create the file", FILE_NAME)
			os.Exit(1)
		}

		defer file.Close()
	}
}

func Update(t *task.Task) error {
	file, err := os.OpenFile(FILE_NAME, os.O_RDWR, 0660)

	defer file.Close()

	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	tasks := []task.Task{}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		currentTask := csvToTask(record)

		if currentTask.Id == t.Id {
			currentTask.Text = t.Text
		}

		tasks = append(tasks, *currentTask)
	}

	if err := ioutil.WriteFile(FILE_NAME, []byte{}, 0664); err != nil {
		return err
	}

	return addMultiple(tasks)
}

func Add(t *task.Task) error {
	return addMultiple([]task.Task{*t})
}

func addMultiple(tasks []task.Task) error {
	file, err := os.OpenFile(FILE_NAME, os.O_RDWR | os.O_APPEND, 0660)

	defer file.Close()

	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range tasks {
		if t.Id == 0 {
			taskId, err := nextId()

			if err != nil {
				return err
			}

			t.Id = taskId
		}

		if err := writer.Write(taskToCsv(&t)); err != nil {
			return err
		}
	}

	return nil
}

func nextId() (int, error) {
	file, err := os.Open(FILE_NAME)

	defer file.Close()

	if err != nil {
		return 0, err
	}

	reader := csv.NewReader(file)
	lastRecord := []string{}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		lastRecord = record
	}

	if len(lastRecord) == 0 {
		return 1, nil
	}

	lastId, _ := strconv.Atoi(lastRecord[0])

	return lastId + 1, nil
}

func taskToCsv(t *task.Task) []string {
	return []string{
		strconv.Itoa(t.Id),
		t.Text,
	}
}

func csvToTask(record []string) *task.Task {
	taskId, _ := strconv.Atoi(record[0])

	return &task.Task{
		taskId,
		record[1],
	}
}
