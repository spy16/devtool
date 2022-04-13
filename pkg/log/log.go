package log

import "log"

func Fatalf(format string, args ...any) {
	log.Fatalf("💣  "+format, args...)
}

func Infof(format string, args ...any) {
	log.Printf("ℹ️  "+format, args...)
}

func Warnf(format string, args ...any) {
	log.Printf("⚠️️  "+format, args...)
}
