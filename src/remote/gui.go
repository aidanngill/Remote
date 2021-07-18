package main

import (
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"github.com/dustin/go-humanize"
)

var fileListBoxModel *ui.TableModel
var consoleBoxAddress *ui.Entry
var serverBoxDevices *ui.Combobox

type ModelHandler struct {
	ID string
}

func newModelHandler() *ModelHandler {
	return new(ModelHandler)
}

func (mh *ModelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // Database Entry ID
		ui.TableString(""), // File Base Name
		ui.TableString(""), // File Size
		ui.TableString(""), // Task Status
		ui.TableInt(0),     // Download Progress
	}
}

func (mh *ModelHandler) NumRows(m *ui.TableModel) int {
	// Upon calling RowInserted(), another "ghost" row is inserted,
	// and so we have to remove one from the length for it not to
	// crash. Unsure as to why this happens.
	length := len(packages)

	if length > 0 {
		length--
	}

	return length
}

func (mh *ModelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	pkg := packages[row]

	switch column {
	case 0:
		return ui.TableString(fmt.Sprint(pkg.Task.ID))
	case 1:
		return ui.TableString(path.Base(pkg.Path))
	case 2:
		return ui.TableString(humanize.Bytes(pkg.Size))
	case 3:
		return ui.TableString(strings.Title(pkg.Task.Status))
	case 4:
		return ui.TableInt(pkg.Progress)
	}

	panic("unreachable")
}

func (mh *ModelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {}

func setupUI() {
	err := ui.Main(func() {
		window := ui.NewWindow("Remote", 870, 260, false)
		window.SetMargined(true)

		boxFull := ui.NewHorizontalBox()
		boxFull.SetPadded(true)

		// Defined earlier as it must be disabled by earlier elements.
		fileListButton := ui.NewButton("Add file")
		fileListButton.Disable()

		// <-- Server Group
		serverGroup := ui.NewGroup("Server")

		serverBox := ui.NewVerticalBox()
		serverBoxDevices = ui.NewCombobox()

		serverBox.Append(ui.NewLabel("Address"), false)
		serverBox.Append(serverBoxDevices, false)

		devices := findNetworkDevices()

		for _, device := range devices {
			serverBoxDevices.Append(device.String())
		}

		serverGroup.SetChild(serverBox)
		serverGroup.SetMargined(true)
		// Server Group -->

		// <-- Console Group
		consoleGroup := ui.NewGroup("Console")

		consoleBox := ui.NewVerticalBox()
		consoleBoxAddress = ui.NewEntry()

		consoleBoxAddress.OnChanged(func(e *ui.Entry) {
			addr := net.ParseIP(consoleBoxAddress.Text())

			if addr != nil {
				fileListButton.Enable()
			} else {
				fileListButton.Disable()
			}
		})

		consoleBox.Append(ui.NewLabel("Address"), false)
		consoleBox.Append(consoleBoxAddress, false)

		consoleGroup.SetChild(consoleBox)
		consoleGroup.SetMargined(true)
		// Console Group -->

		// <-- File List Group
		fileListGroup := ui.NewGroup("Files")

		fileListBox := ui.NewVerticalBox()

		fileListBoxModelHandler := newModelHandler()
		fileListBoxModel = ui.NewTableModel(fileListBoxModelHandler)
		fileListBoxTable := ui.NewTable(&ui.TableParams{
			Model:                         fileListBoxModel,
			RowBackgroundColorModelColumn: 1,
		})

		fileListButton.OnClicked(func(*ui.Button) {
			filePath := ui.OpenFile(window)

			if filePath == "" {
				return
			}

			pkg, err := createPackage(filePath)

			if err != nil {
				ui.MsgBoxError(window, "User Error", fmt.Sprintf("Error encountered, %s.", err.Error()))
				return
			}

			packages = append(packages, pkg)
			fileListBoxModel.RowInserted(int(pkg.Row))

			go func() {
				task, err := createTask(pkg)

				packages[pkg.Row].Task = task
				fileListBoxModel.RowChanged(int(pkg.Row))

				if err != nil {
					ui.MsgBoxError(window, "Console Error", fmt.Sprintf("Error encountered, %s.", err.Error()))
					return
				}

				for {
					progress, err := getTaskProgress(task)

					if err != nil {
						continue
					}

					packages[pkg.Row].Progress = uint64(progress)
					fileListBoxModel.RowChanged(int(pkg.Row))

					if progress >= 100 {
						break
					}

					time.Sleep(1 * time.Second)
				}
			}()
		})

		fileListBoxTable.AppendTextColumn("ID", 0, ui.TableModelColumnNeverEditable, nil)
		fileListBoxTable.AppendTextColumn("Name", 1, ui.TableModelColumnNeverEditable, nil)
		fileListBoxTable.AppendTextColumn("Size", 2, ui.TableModelColumnNeverEditable, nil)
		fileListBoxTable.AppendTextColumn("Status", 3, ui.TableModelColumnNeverEditable, nil)
		fileListBoxTable.AppendProgressBarColumn("Progress", 4)

		fileListBox.Append(fileListBoxTable, true)

		fileListGroup.SetChild(fileListBox)
		fileListGroup.SetMargined(true)
		// File List Group -->

		boxLeft := ui.NewVerticalBox()
		boxLeft.SetPadded(true)

		boxLeft.Append(serverGroup, false)
		boxLeft.Append(consoleGroup, false)
		boxLeft.Append(fileListButton, false)

		boxFull.Append(boxLeft, false)
		boxFull.Append(fileListGroup, true)

		window.SetChild(boxFull)

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()
	})

	if err != nil {
		panic(err)
	}
}
