package log

import (
	"fmt"
	"github.com/rossmcdonald/telegram_hook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

// Create a new instance of the logger. You can have any number of instances.
var Log = logrus.New()

func Logging() {

	// If you wish to add the calling method as a field, instruct the logger via:
	// Note that this does add measurable overhead
	//Log.SetReportCaller(true)

	fmt.Println("-f")

	hook, err := telegram_hook.NewTelegramHook(
		"checkerr",
		"1084260162:AAGyxgi6R_kcnx-TA7caQoTZvrPP2P-FN5c",
		"-1001270332944",
		telegram_hook.WithAsync(true),
		telegram_hook.WithTimeout(10*time.Second),
	)
	if err != nil {
		fmt.Println("× Не получилось создать telegram hook:", err)
	} else {
		fmt.Println("✔ Telegram hook успешно добавлен")
		Log.AddHook(hook)
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	Log.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	//file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err == nil {
	//	fmt.Println("✔ Логирование установлено в файл")
	//	Log.Out = file
	//} else {
	//	fmt.Println("× Failed to log to file, using default stderr")
	//	Log.Info("Failed to log to file, using default stderr")
	//}

	// Use logrus for standard log output
	// Note that `log` here references stdlib's log
	// Not logrus imported under the name `log`.
	log.SetOutput(Log.Writer())
}
