package engine

import (
	"fmt"

	triplesec "github.com/keybase/go-triplesec"
	"github.com/keybase/go/libkb"
)

/*
type SignupUI interface {
	libkb.LogUI
	GPGUI
}
*/

type SignupEngine struct {
	pwsalt     []byte
	tspkey     libkb.TSPassKey
	uid        libkb.UID
	me         *libkb.User
	signingKey libkb.GenericKey
	logui      libkb.LogUI
	gpgui      GPGUI
}

func NewSignupEngine(logui libkb.LogUI, gpgui GPGUI) *SignupEngine {
	return &SignupEngine{logui: logui, gpgui: gpgui}
}

func (s *SignupEngine) Init() error {
	return nil
}

func (s *SignupEngine) CheckRegistered() (err error) {
	G.Log.Debug("+ libkb.SignupJoinEngine::CheckRegistered")
	if cr := G.Env.GetConfig(); cr == nil {
		err = fmt.Errorf("No configuration file available")
	} else if u := cr.GetUid(); u != nil {
		err = libkb.AlreadyRegisteredError{Uid: *u}
	}
	G.Log.Debug("- libkb.SignupJoinEngine::CheckRegistered -> %s", libkb.ErrToOk(err))
	return err
}

func (s *SignupEngine) PostInviteRequest(arg libkb.InviteRequestArg) error {
	return libkb.PostInviteRequest(arg)
}

type SignupEngineRunArg struct {
	Username   string
	Email      string
	InviteCode string
	Passphrase string
	DeviceName string
	SkipGPG    bool
}

// XXX this might need to return more than error...passphraseok, postok, writeok.
func (s *SignupEngine) Run(arg SignupEngineRunArg) error {
	if err := s.genTSPassKey(arg.Passphrase); err != nil {
		return err
	}

	if err := s.join(arg.Username, arg.Email, arg.InviteCode); err != nil {
		return err
	}

	if err := s.registerDevice(arg.DeviceName); err != nil {
		return err
	}

	if err := s.genDetKeys(); err != nil {
		return err
	}

	if !arg.SkipGPG {
		if err := s.checkGPG(); err != nil {
			return err
		}
	}

	return nil
}

func (s *SignupEngine) genTSPassKey(passphrase string) error {
	salt, err := libkb.RandBytes(triplesec.SaltLen)
	if err != nil {
		return err
	}
	s.pwsalt = salt

	s.tspkey, err = libkb.NewTSPassKey(passphrase, salt)
	return err
}

func (s *SignupEngine) join(username, email, inviteCode string) error {
	joinEngine := NewSignupJoinEngine()

	arg := SignupJoinEngineRunArg{
		Username:   username,
		Email:      email,
		InviteCode: inviteCode,
		PWHash:     s.tspkey.PWHash(),
		PWSalt:     s.pwsalt,
	}
	res := joinEngine.Run(arg)
	if res.Err != nil {
		return res
	}

	s.uid = *res.Uid
	user, err := libkb.LoadUser(libkb.LoadUserArg{Uid: res.Uid, PublicKeyOptional: true})
	if err != nil {
		return err
	}
	s.me = user
	return nil
}

func (s *SignupEngine) registerDevice(deviceName string) error {
	eng := NewDeviceEngine(s.me, s.logui)
	err := eng.Run(deviceName)
	if err != nil {
		return err
	}
	s.signingKey = eng.RootSigningKey()
	return nil
}

func (s *SignupEngine) genDetKeys() error {
	eng := NewDetKeyEngine(s.me, s.signingKey, s.logui)
	return eng.Run(s.tspkey.EdDSASeed(), s.tspkey.DHSeed())
}

func (s *SignupEngine) checkGPG() error {
	eng := NewGPG(s.gpgui)
	return eng.Run()
}
