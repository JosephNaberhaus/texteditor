package main

import (
	"fmt"
	"github.com/JosephNaberhaus/go-text-editor"
	"github.com/ahmetb/go-cursor"
	"github.com/eiannone/keyboard"
	"os"
)

func main() {
	editor := go_text_editor.NewEditor()
	editor.SetWidth(20)

	err := keyboard.Open()
	if err != nil {
		fmt.Printf("Error starting keyboard listener")
		os.Exit(2)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			fmt.Printf("Error closing keyboard listener")
		}
	}()

	for {
		fmt.Printf(cursor.ClearEntireScreen())
		fmt.Printf(cursor.MoveTo(0, 0))
		fmt.Printf("Press ESC to exit...\n")
		fmt.Printf(editor.String())
		fmt.Printf(cursor.MoveTo(editor.CursorRow()+2, editor.CursorColumn()+1))

		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Printf("Error getting key")
			break
		}

		if char != 0 {
			editor.Write(string(char))
		} else {
			switch key {
			case keyboard.KeyEsc:
				return
			case keyboard.KeyArrowLeft:
				editor.Left()
			case keyboard.KeyArrowRight:
				editor.Right()
			case keyboard.KeyArrowUp:
				editor.Up()
			case keyboard.KeyArrowDown:
				editor.Down()
			case keyboard.KeyBackspace:
				fallthrough
			case keyboard.KeyBackspace2:
				editor.Backspace()
			case keyboard.KeyEnter:
				editor.Newline()
			case keyboard.KeySpace:
				editor.Write(" ")
			case keyboard.KeyHome:
				editor.Home()
			case keyboard.KeyEnd:
				editor.End()
			}
		}
	}
}
