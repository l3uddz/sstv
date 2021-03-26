package main

import (
	"context"
	"fmt"
	"github.com/koron/go-ssdp"
	"github.com/l3uddz/sstv"
	"github.com/rs/zerolog/log"
	"time"
)

func startSSDP(deviceID string, publicURL string, ctx context.Context) error {
	// credits: https://github.com/tellytv/telly/blob/d507e7ec3f81bc683f904eae80ab34d3142de91b/routes.go#L136
	ad, err := ssdp.Advertise(
		"upnp:rootdevice",
		fmt.Sprintf("uuid:%s::upnp:rootdevice", deviceID),
		sstv.JoinURL(publicURL, "device.xml"),
		"sstv",
		1800)
	if err != nil {
		return fmt.Errorf("ssdp create: %w", err)
	}

	log.Info().Msg("Advertising presence via ssdp")

	// advertiser
	go func(a *ssdp.Advertiser) {
		// graceful ssdp shutdown
		defer func() {
			if err := a.Bye(); err != nil {
				log.Error().
					Err(err).
					Msg("Failed advertising ssdp shutdown")
			}
			if err := a.Close(); err != nil {
				log.Error().
					Err(err).
					Msg("Failed graceful ssdp shutdown")
			}
		}()

		heartbeat := time.Tick(15 * time.Second)
		// advertise presence
		for {
			select {
			case <-heartbeat:
				if err := a.Alive(); err != nil {
					log.Error().
						Err(err).
						Msg("Failed advertising presence via ssdp")
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(ad)
	return nil
}
