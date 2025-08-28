package exql

import (
	"fmt"
	"time"

	"github.com/vedadiyan/exql/lang"
)

func DateTime(args []lang.Value) lang.Value {
	format := time.RFC1123
	switch len := len(args); {
	case len == 0:
		{
			break
		}
	case len == 1:
		{
			val, ok := args[0].(lang.StringValue)
			if !ok {
				return fmt.Errorf("expected StringValue but got %T", args[0])
			}
			format = string(val)
		}
	default:
		{
			return fmt.Errorf("DateTime does not take %d arguments", len)
		}
	}

	return lang.StringValue(time.Now().Format(format))
}
