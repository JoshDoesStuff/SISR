package steaminputbindings

import (
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/Alia5/SISR/helper"
)

type Enforcer interface {
	GetForcedInputAppID() uint32
	ForceInputAppID(appID uint32) error
	ForceOwnAppID() error
}

type enforcer struct {
	forcedAppID uint32

	ownAppID uint32

	mtx sync.Mutex
}

func NewEnforcer() Enforcer {
	var ownAppID uint32
	ownAppIDStr := os.Getenv("SteamAppId")
	if ownAppIDStr != "" {
		appID64, err := strconv.ParseUint(ownAppIDStr, 10, 32)
		if err != nil {
			slog.Error("Failed to parse steamAppIDEnv", "SteamAppID", ownAppIDStr, "error", err)
		} else {
			ownAppID = uint32(appID64)
		}
	}
	return &enforcer{
		ownAppID: ownAppID,
	}
}

func (e *enforcer) GetForcedInputAppID() uint32 {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	return e.forcedAppID
}

func (e *enforcer) ForceInputAppID(appID uint32) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	err := helper.OpenURL("steam://forceinputappid/" + strconv.FormatUint(uint64(appID), 10))
	if err != nil {
		return err
	}

	e.forcedAppID = appID
	return nil
}

func (e *enforcer) ForceOwnAppID() error {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	if e.ownAppID == 0 {
		return ErrNoSteamAppID
	}
	err := helper.OpenURL("steam://forceinputappid/" + strconv.FormatUint(uint64(e.ownAppID), 10))
	if err != nil {
		return err
	}

	e.forcedAppID = e.ownAppID
	return nil
}
