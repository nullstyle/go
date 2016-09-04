package sci

import "log"

func (p *ValueParser) log(args ...interface{}) {
	log.Print(args...)
}

func (p *ValueParser) logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
