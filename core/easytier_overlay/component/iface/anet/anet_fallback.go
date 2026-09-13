package anet

import "net"

func firstNonEmptyIfaces(primary []net.Interface, primaryErr error, fallback []net.Interface, fallbackErr error) ([]net.Interface, error) {
	if primaryErr == nil && len(primary) > 0 {
		return primary, nil
	}
	if fallbackErr == nil && len(fallback) > 0 {
		return fallback, nil
	}
	if primaryErr != nil {
		if fallbackErr != nil {
			return nil, primaryErr
		}
		return fallback, nil
	}
	return primary, nil
}

func firstNonEmptyAddrs(primary []net.Addr, primaryErr error, fallback []net.Addr, fallbackErr error) ([]net.Addr, error) {
	if primaryErr == nil && len(primary) > 0 {
		return primary, nil
	}
	if fallbackErr == nil && len(fallback) > 0 {
		return fallback, nil
	}
	if primaryErr != nil {
		if fallbackErr != nil {
			return nil, primaryErr
		}
		return fallback, nil
	}
	return primary, nil
}

func firstValidIface(primary *net.Interface, primaryErr error, fallback *net.Interface, fallbackErr error) (*net.Interface, error) {
	if primaryErr == nil && primary != nil {
		return primary, nil
	}
	if fallbackErr == nil && fallback != nil {
		return fallback, nil
	}
	if primaryErr != nil {
		if fallbackErr != nil {
			return nil, primaryErr
		}
		return fallback, fallbackErr
	}
	return primary, primaryErr
}

func chooseIfaces(preferNetlink bool, netlink []net.Interface, netlinkErr error, std []net.Interface, stdErr error) ([]net.Interface, error) {
	if preferNetlink {
		return firstNonEmptyIfaces(netlink, netlinkErr, std, stdErr)
	}
	return firstNonEmptyIfaces(std, stdErr, netlink, netlinkErr)
}

func chooseAddrs(preferNetlink bool, netlink []net.Addr, netlinkErr error, std []net.Addr, stdErr error) ([]net.Addr, error) {
	if preferNetlink {
		return firstNonEmptyAddrs(netlink, netlinkErr, std, stdErr)
	}
	return firstNonEmptyAddrs(std, stdErr, netlink, netlinkErr)
}

func chooseIface(preferNetlink bool, netlink *net.Interface, netlinkErr error, std *net.Interface, stdErr error) (*net.Interface, error) {
	if preferNetlink {
		return firstValidIface(netlink, netlinkErr, std, stdErr)
	}
	return firstValidIface(std, stdErr, netlink, netlinkErr)
}
