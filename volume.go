package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Volume is a block of attachable storage for our IAAS products
type Volume struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	InstanceID    string    `json:"instance_id"`
	MountPoint    string    `json:"mountpoint"`
	SizeGigabytes int       `json:"size_gb"`
	Bootable      bool      `json:"bootable"`
	CreatedAt     time.Time `json:"created_at"`
}

// VolumeResult is the response from one of our simple API calls
type VolumeResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

// VolumeConfig are the settings required to create a new Volume
type VolumeConfig struct {
	Name          string `form:"name"`
	SizeGigabytes int    `form:"size_gb"`
	Bootable      bool   `form:"bootable"`
}

// ListVolumes returns all volumes owned by the calling API account
func (c *Client) ListVolumes() ([]Volume, error) {
	resp, err := c.SendGetRequest("/v2/volumes")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var volumes = make([]Volume, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&volumes); err != nil {
		return nil, err
	}

	return volumes, nil
}

// FindVolume finds a volume by either part of the ID or part of the name
func (c *Client) FindVolume(search string) (*Volume, error) {
	volumes, err := c.ListVolumes()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, volume := range volumes {
		if strings.Contains(volume.ID, search) || strings.Contains(volume.Name, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &volumes[found], nil
}

// NewVolume creates a new volume
func (c *Client) NewVolume(v *VolumeConfig) (*VolumeResult, error) {
	body, err := c.SendPostRequest("/v2/volumes/", v)
	if err != nil {
		return nil, err
	}

	var result = &VolumeResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ResizeVolume resizes a volume
func (c *Client) ResizeVolume(id string, size int) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/resize", id), map[string]int{
		"size_gb": size,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// AttachVolume attaches a volume to an intance
func (c *Client) AttachVolume(id string, instance string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/attach", id), map[string]string{
		"instance_id": instance,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DetachVolume attach volume from any instances
func (c *Client) DetachVolume(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/detach", id), "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DeleteVolume deletes a volumes
func (c *Client) DeleteVolume(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/volumes/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
