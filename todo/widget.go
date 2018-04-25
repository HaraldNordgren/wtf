package todo

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gdamore/tcell"
	"github.com/olebedev/config"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
	"gopkg.in/yaml.v2"
)

// Config is a pointer to the global config object
var Config *config.Config

type Widget struct {
	wtf.TextWidget

	pages    *tview.Pages
	filePath string
	list     *List
}

func NewWidget(pages *tview.Pages) *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget(" 📝 Todo ", "todo"),

		pages:    pages,
		filePath: Config.UString("wtf.mods.todo.filename"),
		list:     &List{selected: -1},
	}

	widget.init()
	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.load()
	widget.display()
	widget.RefreshedAt = time.Now()
}

/* -------------------- Unexported Functions -------------------- */

// edit opens a modal dialog that permits editing the text of the currently-selected item
func (widget *Widget) edit() {
	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				widget.pages.RemovePage("edit")
			}
		})

	widget.pages.AddPage("edit", modal, false, true)
	//widget.app.SetFocus(modal)
}

func (widget *Widget) init() {
	_, err := wtf.CreateFile(widget.filePath)
	if err != nil {
		panic(err)
	}
}

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case " ":
		// Check/uncheck selected item
		widget.list.Toggle()
		widget.persist()
		widget.display()
		return nil
	case "e":
		// Edit selected item
		widget.edit()
		return nil
	case "h":
		// Show help menu
		fmt.Println("HELP!")
		return nil
	case "j":
		widget.list.Next()
		widget.display()
		return nil
	case "k":
		widget.list.Prev()
		widget.display()
		return nil
	case "n":
		// Add a new item
		return nil
	case "o":
		// Open the file
		//widget.openFile()
		wtf.OpenFile(widget.filePath)
		return nil
	}

	switch event.Key() {
	case tcell.KeyCtrlD:
		// Delete the selected item
		widget.list.Delete()
		widget.persist()
		widget.display()
		return nil
	case tcell.KeyCtrlJ:
		// Move selected item down in the list
		widget.list.Demote()
		widget.persist()
		widget.display()
		return nil
	case tcell.KeyCtrlK:
		// Move selected item up in the list
		widget.list.Promote()
		widget.persist()
		widget.display()
		return nil
	case tcell.KeyDown:
		// Select the next item down
		widget.list.Next()
		widget.display()
		return nil
	case tcell.KeyEsc:
		// Unselect the current row
		widget.list.Unselect()
		widget.display()
		return event
	case tcell.KeyUp:
		// Select the next item up
		widget.list.Prev()
		widget.display()
		return nil
	default:
		// Pass it along
		return event
	}
}

// Loads the todo list from Yaml file
func (widget *Widget) load() {
	confDir, _ := wtf.ConfigDir()
	filePath := fmt.Sprintf("%s/%s", confDir, widget.filePath)

	fileData, _ := wtf.ReadFileBytes(filePath)
	yaml.Unmarshal(fileData, &widget.list)
}

// persist writes the todo list to Yaml file
func (widget *Widget) persist() {
	confDir, _ := wtf.ConfigDir()
	filePath := fmt.Sprintf("%s/%s", confDir, widget.filePath)

	fileData, _ := yaml.Marshal(&widget.list)

	err := ioutil.WriteFile(filePath, fileData, 0644)

	if err != nil {
		panic(err)
	}
}