package browsertest

import "strings"

func IsNetworkChangedError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NETWORK_CHANGED")
}
