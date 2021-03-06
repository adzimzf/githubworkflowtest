package config

import (
	"fmt"
	"net/url"
)

type ProxySetterStation interface {
	Execute(p *Proxy) error
}

type ProxySetter struct{}

func NewProxySetterStations() *ProxySetter {
	return &ProxySetter{}
}

func (ps *ProxySetter) Execute(p *Proxy) error {
	twoFAStt := NewSetTwoFAStation(nil)
	userNameStt := NewSetUserNameStation(twoFAStt)
	authConnectorStt := NewSetAuthConnectorStation(userNameStt)
	addrStt := NewSetAddressStation(authConnectorStt)
	envStt := NewSetEnvStation(addrStt)
	return envStt.Execute(p)
}

type SetEnvStation struct {
	next ProxySetterStation
}

func NewSetEnvStation(next ProxySetterStation) *SetEnvStation {
	return &SetEnvStation{next}
}

func (s *SetEnvStation) Execute(p *Proxy) error {
	var err error
	p.Env, err = prompt("Environment", func(env string) error {
		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetAddressStation struct {
	next ProxySetterStation
}

func NewSetAddressStation(next ProxySetterStation) *SetAddressStation {
	return &SetAddressStation{next}
}

func (s *SetAddressStation) Execute(p *Proxy) error {
	var err error
	p.Address, err = prompt("Proxy Address (with http protocol)", func(address string) error {
		_, err := url.ParseRequestURI(address)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetUserNameStation struct {
	next ProxySetterStation
}

func NewSetUserNameStation(next ProxySetterStation) *SetUserNameStation {
	return &SetUserNameStation{next}
}

func (s *SetUserNameStation) Execute(p *Proxy) error {
	var err error
	p.UserName, err = prompt("Username (teleport username)", func(userName string) error {
		if p.AuthConnector == "" {
			return fmt.Errorf("Username OR Auth Connector is required")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetTwoFAStation struct {
	next ProxySetterStation
}

func NewSetTwoFAStation(next ProxySetterStation) *SetTwoFAStation {
	return &SetTwoFAStation{next}
}

func (s *SetTwoFAStation) Execute(p *Proxy) error {
	isTwoFA, err := prompt("Is Need 2FA (Y/y/N/n)", func(towFA string) error {
		if towFA == "Y" || towFA == "y" || towFA == "N" || towFA == "n" {
			return nil
		}
		return fmt.Errorf("invalid formatting")
	})

	if err != nil {
		return err
	}

	if isTwoFA == "Y" || isTwoFA == "y" {
		p.TwoFA = true
	}

	return determineNext(s.next, p)
}

type SetAuthModeStation struct {
	next ProxySetterStation
}

func NewSetAuthConnectorStation(next ProxySetterStation) *SetAuthModeStation {
	return &SetAuthModeStation{next}
}

func (s *SetAuthModeStation) Execute(p *Proxy) error {
	var err error
	p.AuthConnector, err = prompt("Auth Connector", func(towFA string) error {
		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

func determineNext(next ProxySetterStation, p *Proxy) error {
	if next != nil {
		return next.Execute(p)
	}

	return nil
}
