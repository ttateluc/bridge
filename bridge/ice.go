package bridge

import (
	"fmt"

	"github.com/duality-solutions/web-bridge/configs/settings"
	webrtc "github.com/pion/webrtc/v2"
)

// newIceSetting create a new WebRTC ICE Server setting
func newIceSetting(config *settings.Configuration) (*webrtc.ICEServer, error) {
	iceServer := *config.IceServers()
	if (len(*config.IceServers())) == 0 {
		return nil, fmt.Errorf("no ICE service URL found")
	}
	urls := []string{iceServer[0].URL}
	iceSettings := webrtc.ICEServer{
		URLs:           urls,
		Username:       iceServer[0].UserName,
		Credential:     iceServer[0].Credential,
		CredentialType: webrtc.ICECredentialTypePassword,
	}
	return &iceSettings, nil
}

func connectToIceServicesOption(config *settings.Configuration, detached bool) (*webrtc.PeerConnection, error) {
	iceServer, err := newIceSetting(config)
	if err != nil {
		return nil, fmt.Errorf("NewIceSetting %v", err)
	}
	configICE := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{*iceServer},
	}

	s := webrtc.SettingEngine{}
	if detached {
		s.DetachDataChannels()
	}

	// Create an API object with the engine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))
	peerConnection, err := api.NewPeerConnection(configICE)
	if err != nil {
		return nil, fmt.Errorf("NewPeerConnection %v", err)
	}
	return peerConnection, nil
}

// ConnectToIceServices uses the configuration settings to establish a connection with ICE servers
func ConnectToIceServices(config *settings.Configuration) (*webrtc.PeerConnection, error) {
	return connectToIceServicesOption(config, false)
}

// ConnectToIceServicesDetached uses the configuration settings to establish a connection with ICE servers with detached channels
func ConnectToIceServicesDetached(config *settings.Configuration) (*webrtc.PeerConnection, error) {
	return connectToIceServicesOption(config, true)
}

// DisconnectIceService calls the close method for a WebRTC peer connection
func DisconnectIceService(pc *webrtc.PeerConnection) error {
	return pc.Close()
}
