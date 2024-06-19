/*
Copyright © 2024 Bridge Digital
*/
package predefined

import (
	"github.com/mgutz/ansi"
)

func BuildAnsw(q string, a string) string {
	aFormat := ansi.ColorFunc("cyan")
	aResult := aFormat(a)

	qFormat := ansi.ColorFunc("default+hb")
	qResult := qFormat(q)

	mFormat := ansi.ColorFunc("green+hb")
	mResult := mFormat("? ")

	result := mResult + qResult + aResult

	return result
}

func BuildOk(m string) string {
	formattedMsg := ansi.ColorFunc("green+hb")
	mResult := formattedMsg(m)

	return mResult
}

func BuildError(m string) string {
	formattedMsg := ansi.ColorFunc("red+hb")
	mResult := formattedMsg(m)

	return mResult
}

func BuildWarning(m string) string {
	formattedMsg := ansi.ColorFunc("yellow+hb")
	mResult := formattedMsg(m)

	return mResult
}
