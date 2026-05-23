package input

import (
	"context"
	"log/slog"
	"slices"

	"github.com/Alia5/VIIPER/apiclient"
)

type ViiperBridge interface {
	AttachViiperDevice(ctx context.Context, d *Device) error
}

type viiperBridge struct {
	client *apiclient.Client
	busID  uint32
}

func NewViiperBridge() ViiperBridge {
	return &viiperBridge{
		client: apiclient.New("localhost:3242"),
	}
}

func (vb *viiperBridge) AttachViiperDevice(ctx context.Context, d *Device) error {
	err := vb.ensureBus(ctx)
	if err != nil {
		return err
	}

	controlStream, apiDevice, err := vb.client.AddDeviceAndConnect(ctx, vb.busID, "xbox360", nil)
	if err != nil {
		return err
	}
	d.viiperDevice = NewViiperDevice(controlStream, apiDevice)
	slog.Info("Created VIIPER device", "viiperDevice", *apiDevice, "busID", vb.busID)

	return nil
}

func (vb *viiperBridge) ensureBus(ctx context.Context) error {
	if vb.busID == 0 {
		slog.Debug("No previews VIIPER bus used, creating new bus")
		bus, err := vb.client.BusCreateCtx(ctx, 0)
		if err != nil {
			return err
		}
		vb.busID = bus.BusID
		slog.Info("Created VIIPER bus", "busID", vb.busID)
		return nil
	}

	busResp, err := vb.client.BusListCtx(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(busResp.Buses, vb.busID) {
		slog.Warn("Created VIIPER bus not found; Recreating")
		bus, err := vb.client.BusCreateCtx(ctx, vb.busID)
		if err != nil {
			return err
		}
		vb.busID = bus.BusID
		slog.Info("Re-Created VIIPER bus", "busID", vb.busID)
	}
	return nil
}
