package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	widgetx "fyne.io/x/fyne/widget"
	"github.com/lostdusty/gobalt/v2"
)

func showFirstRunScreen() fyne.CanvasObject {

	welcomeText := widget.NewRichTextFromMarkdown("# Welcome!\n\nLooks like it's the first time you're launching this app, so we need to do a quick setup before you can download.\n\nCobalt doesn't provide a public api anymore, so you need to host your own instance or use someone's instance, with their permission.")
	welcomeText.Wrapping = fyne.TextWrapWord

	welcomeSep := canvas.NewLine(theme.Color(theme.ColorNamePrimary))
	welcomeSep.StrokeWidth = 1.5

	welcomeTextNext := widget.NewLabel("To proceed, type or select an instance below. Click on the button on the right to refresh the instance list.")

	welcomeInstanceSelector := widgetx.NewCompletionEntry(nil)
	welcomeEntryApiKey := widget.NewPasswordEntry()
	welcomeEntryApiKey.PlaceHolder = "Something like this: 123e4567-e89b-12d3-a456-426655440000"
	welcomeEntryApiKey.ActionItem = &widget.Button{
		Icon: theme.HelpIcon(),
		OnTapped: func() {
			dialog.ShowInformation("Help - Api Keys", "You might need an api key if you wanna use certain instances.\nAsk for the instance owner for one.\nThis field can be blank if a key is not needed.", windowGuibalt)
		},
	}
	revealApiKeyBtn := widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
		welcomeEntryApiKey.Password = !welcomeEntryApiKey.Password
		welcomeEntryApiKey.Refresh()
	})
	formFieldApiKey := container.NewBorder(nil, nil, nil, revealApiKeyBtn, welcomeEntryApiKey)

	finishSetupBtn := widget.NewButtonWithIcon("Proceed", theme.NavigateNextIcon(), func() {
		popupVerifyInstance := dialog.NewCustomWithoutButtons("Testing this instance", loadingContent("Checking instance..."), windowGuibalt)
		popupVerifyInstance.Show()

		if !strings.HasPrefix(welcomeInstanceSelector.Text, "https") {
			log.Println("triggered, new string: " + "https://" + welcomeInstanceSelector.Text)
			welcomeInstanceSelector.Text = "https://" + welcomeInstanceSelector.Text
			welcomeInstanceSelector.Refresh()
		}

		config := gobalt.CreateDefaultSettings()
		config.Url = "https://www.youtube.com/watch?v=aQvGIIdgFDM"
		gobalt.CobaltApi = welcomeInstanceSelector.Text
		gobalt.ApiKey = welcomeEntryApiKey.Text
		//log.Println(welcomeInstanceSelector.Text, welcomeEntryApiKey.Text)
		_, err := gobalt.Run(config)
		if err != nil {
			//log.Println(err)
			popupVerifyInstance.Hide()
			log.Printf("Err: %v | Instance: %v | API-Key: %v | Gobalt config: %v", err, gobalt.CobaltApi, gobalt.ApiKey, config)
			if err.Error() == "error.api.youtube.login" {
				dialog.ShowConfirm("Something is wrong", "This instance is working, but it can't download YouTube videos.\nDo you want to use this instance anyway?", func(b bool) {
					if b {
						appGuibalt.Preferences().SetString("instance", gobalt.CobaltApi)
						appGuibalt.Preferences().SetString("api-key", gobalt.ApiKey)
						appGuibalt.Preferences().SetBool("first-run", true)
						windowGuibalt.SetContent(showMainScreen())
					}
				}, windowGuibalt)
				return
			}

			dialog.ShowError(fmt.Errorf("unable to use selected instance, did you set an api key?\n(Details: %v)", err), windowGuibalt)
			return
		}

		appGuibalt.Preferences().SetString("instance", gobalt.CobaltApi)
		appGuibalt.Preferences().SetString("api-key", gobalt.ApiKey)
		appGuibalt.Preferences().SetBool("first-run", true)

		windowGuibalt.SetContent(showMainScreen())

		popupVerifyInstance.Hide()
	})

	finishSetupBtn.Disable()
	finishSetupBtn.Importance = widget.SuccessImportance
	finishSetupBtn.IconPlacement = widget.ButtonIconTrailingText

	welcomeInstanceSelector.ActionItem = &widget.Button{
		Icon: theme.ViewRefreshIcon(),
		OnTapped: func() {
			log.Println("refreshing instances...")

			popupRefreshing := widget.NewModalPopUp(loadingContent("Refreshing instance list..."), windowGuibalt.Canvas())
			popupRefreshing.Show()

			instances, err := refreshInstances()
			if err != nil {
				popupRefreshing.Hide()
				dialog.ShowError(errors.New("unable to refresh instances.\nverify your internet connection"), windowGuibalt)
			}

			welcomeInstanceSelector.SetOptions(instances)

			popupRefreshing.Hide()

			log.Println("done refreshing, instances:", instances)
		},
	}

	welcomeInstanceSelector.PlaceHolder = "api.cobalt.tools"

	welcomeInstanceSelector.Validator = func(s string) error {
		if s == "" {
			finishSetupBtn.Disable()
			return errors.New("please set one instance to proceed")
		}
		finishSetupBtn.Enable()
		return nil
	}

	welcomeInstanceSelector.OnChanged = func(s string) {
		if len(welcomeInstanceSelector.Options) > 1 {
			welcomeInstanceSelector.ShowCompletion()
		}
	}

	formInstanceInfo := widget.NewForm(&widget.FormItem{
		Text:     "Instance URL",
		Widget:   welcomeInstanceSelector,
		HintText: "Click on the refresh button to get a list of public instances.",
	}, &widget.FormItem{
		Text:     "Instance API Key",
		Widget:   formFieldApiKey,
		HintText: "Click on the \"?\" button for help",
	})

	welcomeLayout := container.NewBorder(nil, finishSetupBtn, nil, nil, container.NewVBox(welcomeText, welcomeSep, welcomeTextNext, formInstanceInfo))

	return welcomeLayout
}

func showMainScreen() fyne.CanvasObject {
	//downloadOptions := gobalt.CreateDefaultSettings()
	btnConfig := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		windowGuibalt.SetContent(showConfigScreen())
	})

	btnDownloadQueue := widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		//queue screen
	})
	headerTitleApp := widget.NewLabelWithStyle("Guibalt", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	headerContainer := container.NewBorder(nil, nil, btnConfig, btnDownloadQueue, headerTitleApp)
	container.NewThemeOverride(headerContainer, themeNoBg{})

	cardProgress := widget.NewProgressBar()
	cardTitle := widget.NewRichTextFromMarkdown("## download.title")
	cardSubtitle := widget.NewLabel("download.description")
	cardImg := &canvas.Image{Resource: theme.FyneLogo(), FillMode: canvas.ImageFillContain}
	layoutCard := container.NewGridWithColumns(2, cardImg, container.NewVBox(cardTitle, cardSubtitle, cardProgress))

	downloadEntry := widget.NewEntry()
	downloadEntry.PlaceHolder = "https://www.youtube.com/watch?v=Q7M5v8UAZAI"
	btnPasteFromClip := widget.NewButtonWithIcon("", theme.ContentPasteIcon(), nil)
	//btnPasteFromClip.Importance = widget.HighImportance
	downloadEntry.ActionItem = btnPasteFromClip

	btnDownload := widget.NewButtonWithIcon("", theme.DownloadIcon(), nil)
	btnDownload.Importance = widget.HighImportance

	downloadInputLayout := container.NewBorder(nil, nil, nil, btnDownload, downloadEntry)

	return container.NewVBox(headerContainer, layout.NewSpacer(), layoutCard, downloadInputLayout, layout.NewSpacer())
}

func showConfigScreen() fyne.CanvasObject {
	return widgetLoading()
}

func refreshInstances() ([]string, error) {
	getInstances, err := gobalt.GetCobaltInstances()
	if err != nil {
		return nil, err
	}
	var instancesList []string
	for _, v := range getInstances {
		instancesList = append(instancesList, "https://"+v.API)
	}
	return instancesList, nil
}

func widgetLoading() fyne.CanvasObject {
	load, _ := widgetx.NewAnimatedGifFromResource(resourceLoadingGif)
	load.SetMinSize(fyne.NewSquareSize(16))
	load.Start()
	return load
}

func loadingContent(text string) fyne.CanvasObject {
	return container.NewHBox(widgetLoading(), widget.NewLabel(text))
}
