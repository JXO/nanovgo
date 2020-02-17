package demo

import (
	"fmt"
	"github.com/jxo/davinci/ui"
	"github.com/jxo/davinci/vg"
	"math"
	"strconv"
)

func ButtonDemo(screen *ui.Screen) {
	window := ui.NewWindow(screen, "Button demo")
	window.SetPosition(15, 15)
	window.SetLayout(ui.NewGroupLayout())

	ui.NewLabel(window, "Push buttons").SetFont("sans-bold")

	b1 := ui.NewButton(window, "Plain button")
	b1.SetCallback(func() {
		fmt.Println("pushed!")
	})

	b2 := ui.NewButton(window, "Styled")
	b2.SetBackgroundColor(vg.RGBA(0, 0, 255, 25))
	b2.SetIcon(ui.IconRocket)
	b2.SetCallback(func() {
		fmt.Println("pushed!")
	})

	ui.NewLabel(window, "Toggle button").SetFont("sans-bold")
	b3 := ui.NewButton(window, "Toggle me")
	b3.SetFlags(ui.ToggleButtonType)
	b3.SetChangeCallback(func(state bool) {
		fmt.Println("Toggle button state:", state)
	})

	ui.NewLabel(window, "Radio buttons").SetFont("sans-bold")
	b4 := ui.NewButton(window, "Radio button 1")
	b4.SetFlags(ui.RadioButtonType)
	b5 := ui.NewButton(window, "Radio button 2")
	b5.SetFlags(ui.RadioButtonType)

	ui.NewLabel(window, "A tool palette").SetFont("sans-bold")
	tools := ui.NewWidget(window)
	tools.SetLayout(ui.NewBoxLayout(ui.Horizontal, ui.Middle, 0, 6))

	ui.NewToolButton(tools, ui.IconCloud)
	ui.NewToolButton(tools, ui.IconFastForward)
	ui.NewToolButton(tools, ui.IconCompass)
	ui.NewToolButton(tools, ui.IconInstall)

	ui.NewLabel(window, "Popup buttons").SetFont("sans-bold")
	b6 := ui.NewPopupButton(window, "Popup")
	b6.SetIcon(ui.IconExport)
	popup := b6.Popup()
	popup.SetLayout(ui.NewGroupLayout())

	ui.NewLabel(popup, "Arbitrary widgets can be placed here").SetFont("sans-bold")
	ui.NewCheckBox(popup, "A check box")
	b7 := ui.NewPopupButton(popup, "Recursive popup")
	b7.SetIcon(ui.IconFlash)
	popup2 := b7.Popup()

	popup2.SetLayout(ui.NewGroupLayout())
	ui.NewCheckBox(popup2, "Another check box")
}

func BasicWidgetsDemo(screen *ui.Screen, images []ui.Image) (*ui.PopupButton, *ui.ImagePanel, *ui.ProgressBar) {
	window := ui.NewWindow(screen, "Basic widgets")
	window.SetPosition(230, 15)
	window.SetLayout(ui.NewGroupLayout())

	ui.NewLabel(window, "Message dialog").SetFont("sans-bold")

	tools := ui.NewWidget(window)
	tools.SetLayout(ui.NewBoxLayout(ui.Horizontal, ui.Middle, 0, 6))

	b1 := ui.NewButton(tools, "Info")
	b1.SetCallback(func() {

	})
	b2 := ui.NewButton(tools, "Warn")
	b2.SetCallback(func() {

	})
	b3 := ui.NewButton(tools, "Ask")
	b3.SetCallback(func() {

	})

	ui.NewLabel(window, "Image panel & scroll panel").SetFont("sans-bold")
	imagePanelButton := ui.NewPopupButton(window, "Image Panel")
	imagePanelButton.SetIcon(ui.IconFolder)
	popup := imagePanelButton.Popup()
	imgPanel := ui.NewImagePanel(popup)
	imgPanel.SetImages(images)
	popup.SetFixedSize(245, 150)

	ui.NewLabel(window, "File dialog").SetFont("sans-bold")

	tools2 := ui.NewWidget(window)
	tools2.SetLayout(ui.NewBoxLayout(ui.Horizontal, ui.Middle, 0, 6))

	b4 := ui.NewButton(tools2, "Open")
	b4.SetCallback(func() {

	})
	b5 := ui.NewButton(tools2, "Save")
	b5.SetCallback(func() {

	})

	ui.NewLabel(window, "Combo box").SetFont("sans-bold")
	ui.NewComboBox(window, []string{"Combo box item 1", "Combo box item 2", "Combo box item 3"})

	ui.NewLabel(window, "Check box").SetFont("sans-bold")
	cb1 := ui.NewCheckBox(window, "Flag 1")
	cb1.SetCallback(func(checked bool) {
		fmt.Println("Check box 1 state:", checked)
	})
	cb1.SetChecked(true)

	cb2 := ui.NewCheckBox(window, "Flag 2")
	cb2.SetCallback(func(checked bool) {
		fmt.Println("Check box 2 state:", checked)
	})
	ui.NewLabel(window, "Progress bar").SetFont("sans-bold")
	progress := ui.NewProgressBar(window)

	ui.NewLabel(window, "Slider and text box").SetFont("sans-bold")
	panel := ui.NewWidget(window)
	panel.SetLayout(ui.NewBoxLayout(ui.Horizontal, ui.Middle, 0, 20))
	slider := ui.NewSlider(panel)
	slider.SetValue(0.5)
	slider.SetFixedWidth(80)

	textBox := ui.NewTextBox(panel)
	textBox.SetFixedSize(60, 25)
	textBox.SetFontSize(20)
	textBox.SetAlignment(ui.TextRight)
	textBox.SetValue("50")
	textBox.SetUnits("%")

	slider.SetCallback(func(value float32) {
		textBox.SetValue(strconv.FormatInt(int64(value*100), 10))
	})
	slider.SetFinalCallback(func(value float32) {
		fmt.Printf("Final slider value: %d\n", int(value*100))
	})

	return imagePanelButton, imgPanel, progress
}

func MiscWidgetsDemo(screen *ui.Screen) {
	window := ui.NewWindow(screen, "Misc. widgets")
	window.SetPosition(445, 15)
	window.SetLayout(ui.NewGroupLayout())

	ui.NewLabel(window, "Color wheel").SetFont("sans-bold")
	ui.NewColorWheel(window)

	ui.NewLabel(window, "Color picker").SetFont("sans-bold")
	ui.NewColorPicker(window)

	ui.NewLabel(window, "Function graph").SetFont("sans-bold")
	graph := ui.NewGraph(window, "Some function")
	graph.SetHeader("E = 2.35e-3")
	graph.SetFooter("Iteration 89")
	fValues := make([]float32, 100)
	for i := 0; i < 100; i++ {
		x := float64(i)
		fValues[i] = 0.5 * float32(0.5*math.Sin(x/10.0)+0.5*math.Cos(x/23.0)+1.0)
	}
	graph.SetValues(fValues)
}

func GridDemo(screen *ui.Screen) {
	window := ui.NewWindow(screen, "Grid of small widgets")
	window.SetPosition(445, 358)
	layout := ui.NewGridLayout(ui.Horizontal, 2, ui.Middle, 15, 5)
	layout.SetColAlignment(ui.Maximum, ui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	{
		ui.NewLabel(window, "Regular text :").SetFont("sans-bold")
		textBox := ui.NewTextBox(window, "日本語")
		textBox.SetFont("japanese")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
	}
	{
		ui.NewLabel(window, "Floating point :").SetFont("sans-bold")
		textBox := ui.NewTextBox(window, "50.0")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("GiB")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	}
	{
		ui.NewLabel(window, "Positive integer :").SetFont("sans-bold")
		textBox := ui.NewTextBox(window, "50")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	{
		ui.NewLabel(window, "Float box :").SetFont("sans-bold")
		floatBox := ui.NewFloatBox(window, 10.0)
		floatBox.SetEditable(true)
		floatBox.SetFixedSize(100, 20)
		floatBox.SetUnits("GiB")
		floatBox.SetDefaultValue(0.0)
		floatBox.SetFontSize(16)
	}
	{
		ui.NewLabel(window, "Int box :").SetFont("sans-bold")
		intBox := ui.NewIntBox(window, true, 50)
		intBox.SetEditable(true)
		intBox.SetFixedSize(100, 20)
		intBox.SetUnits("MHz")
		intBox.SetDefaultValue(0)
		intBox.SetFontSize(16)
	}
	{
		ui.NewLabel(window, "Checkbox :").SetFont("sans-bold")
		checkbox := ui.NewCheckBox(window, "Check me")
		checkbox.SetFontSize(16)
		checkbox.SetChecked(true)
	}
	{
		ui.NewLabel(window, "Combobox :").SetFont("sans-bold")
		combobox := ui.NewComboBox(window, []string{"Item 1", "Item 2", "Item 3"})
		combobox.SetFontSize(16)
		combobox.SetFixedSize(100, 20)
	}
	{
		ui.NewLabel(window, "Color button :").SetFont("sans-bold")

		popupButton := ui.NewPopupButton(window, "")
		popupButton.SetBackgroundColor(vg.RGBA(255, 120, 0, 255))
		popupButton.SetFontSize(16)
		popupButton.SetFixedSize(100, 20)
		popup := popupButton.Popup()
		popup.SetLayout(ui.NewGroupLayout())

		colorWheel := ui.NewColorWheel(popup)
		colorWheel.SetColor(popupButton.BackgroundColor())

		colorButton := ui.NewButton(popup, "Pick")
		colorButton.SetFixedSize(100, 25)
		colorButton.SetBackgroundColor(colorWheel.Color())

		colorWheel.SetCallback(func(color vg.Color) {
			colorButton.SetBackgroundColor(color)
		})

		colorButton.SetChangeCallback(func(pushed bool) {
			if pushed {
				popupButton.SetBackgroundColor(colorButton.BackgroundColor())
				popupButton.SetPushed(false)
			}
		})
	}
}

func SelectedImageDemo(screen *ui.Screen, imageButton *ui.PopupButton, imagePanel *ui.ImagePanel) {
	window := ui.NewWindow(screen, "Selected image")
	window.SetPosition(685, 15)
	window.SetLayout(ui.NewGroupLayout())

	img := ui.NewImageView(window)
	img.SetPolicy(ui.ImageSizePolicyExpand)
	img.SetFixedSize(300, 300)
	img.SetImage(imagePanel.Images()[0].ImageID)

	imagePanel.SetCallback(func(index int) {
		img.SetImage(imagePanel.Images()[index].ImageID)
	})

	cb := ui.NewCheckBox(window, "Expand")
	cb.SetCallback(func(checked bool) {
		if checked {
			img.SetPolicy(ui.ImageSizePolicyExpand)
		} else {
			img.SetPolicy(ui.ImageSizePolicyFixed)
		}
	})
	cb.SetChecked(true)
}
