// Copyright 2023 qbee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package client

// DeviceAttributes represents additional device attributes set for the device.
type DeviceAttributes struct {
	// DeviceName is the name of the device which is displayed in the UI (if set).
	DeviceName string `json:"device_name"`

	// Description is the extra device description.
	Description string `json:"description"`

	// Country is the country where device is located.
	Country string `json:"country"`

	// City is the city where device is located.
	City string `json:"city"`

	// Zip is the zip code where device is located.
	Zip string `json:"zip"`

	// Address is the address where device is located.
	Address string `json:"address"`

	// Latitude is the latitude where device is located.
	Latitude string `json:"latitude"`

	// Longitude is the longitude where device is located.
	Longitude string `json:"longitude"`
}
