//go:build windows

package php

import (
	"github.com/php-any/origami/data"
	"golang.org/x/sys/windows"
)

func initPlatformConstants(vm data.VM) {
	version := windows.RtlGetVersion()
	for name, value := range map[string]int{
		"PHP_WINDOWS_VERSION_MAJOR":       int(version.MajorVersion),
		"PHP_WINDOWS_VERSION_MINOR":       int(version.MinorVersion),
		"PHP_WINDOWS_VERSION_BUILD":       int(version.BuildNumber),
		"PHP_WINDOWS_VERSION_PLATFORM":    int(version.PlatformId),
		"PHP_WINDOWS_VERSION_SP_MAJOR":    int(version.ServicePackMajor),
		"PHP_WINDOWS_VERSION_SP_MINOR":    int(version.ServicePackMinor),
		"PHP_WINDOWS_VERSION_SUITEMASK":   int(version.SuiteMask),
		"PHP_WINDOWS_VERSION_PRODUCTTYPE": int(version.ProductType),
	} {
		vm.SetConstant(name, data.NewIntValue(value))
	}
}
